package hamt32

import "context"

// HamtTransient is the data structure which the Transient Hamt methods are
// called upon. In fact it is identical to the HamtFunctional data structure and
// all the table and leaf data structures it uses are the same ones used by the
// HamtTransient implementation. It is its own type so that the methods it calls
// are the transient version of the Hamt interface.
//
// The Transient version of the Hamt data structure, does all modifications
// in-place. So sharing this datastruture between threads is NOT safe unless
// you were to implement a locking stategy CORRECTLY.
type HamtTransient struct {
	hamtBase
}

// NewTransient constructs a new HamtTransient data structure based on the opt
// argument.
func NewTransient(opt int) *HamtTransient {
	var h = new(HamtTransient)

	h.hamtBase.init(opt)

	return h
}

// IsEmpty simply returns if the HamtTransient datastucture has no entries.
func (h *HamtTransient) IsEmpty() bool {
	return h.hamtBase.IsEmpty()
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtTransient data structure.
func (h *HamtTransient) Nentries() uint {
	return h.hamtBase.Nentries()
}

// ToFunctional just recasts the HamtFunctional pointer to a HamtFunctional
// underneath the Hamt interface.
//
// If you want a copy of the HamtTransient data structure over to a completely
// independent HamtFunctional data structure, you should first do a DeepCopy
// followed by a ToFunctional call.
func (h *HamtTransient) ToFunctional() Hamt {
	var nh = (*HamtFunctional)(h)
	return nh
}

// ToTransient does nothing to a HamtTransient pointer. This method
// only here for conformance with the Hamt interface.
func (h *HamtTransient) ToTransient() Hamt {
	return h
}

// DeepCopy() copies the HamtTransient data structure and every table it
// contains recursively.
func (h *HamtTransient) DeepCopy() Hamt {
	var nh = new(HamtTransient)
	nh.root = *h.root.deepCopy().(*fixedTable)
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

// Get retrieves the value related to the key in the HamtTransient
// data structure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtTransient data structure.
func (h *HamtTransient) Get(key []byte) (interface{}, bool) {
	return h.hamtBase.Get(key)
}

// Put stores a new (key,value) pair in the HamtTransient data structure. It
// returns a bool indicating if a new pair were added or if the value replaced
// the value in a previously stored (key,value) pair. Either way it returns and
// new HamtTransient data structure containing the modification.
func (h *HamtTransient) Put(key []byte, val interface{}) (Hamt, bool) {
	// Doing this in newFlatLeaf() and leafI.put().
	//key = copyKey(key)

	var hv = calcHashVal(key)
	var path, leaf, idx = h.find(hv)

	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool

	if leaf == nil {
		//check if upgrading allowed & if it is required
		if h.grade && curTable != &h.root &&
			(curTable.nentries()+1) == UpgradeThreshold {
			var newTable = upgradeToFixedTable(
				curTable.Hash(), depth, curTable.entries())

			var parentTable = path.peek()
			var parentIdx = hv.Index(depth - 1)
			parentTable.replace(parentIdx, newTable)

			curTable = newTable
		}
		curTable.insert(idx, newFlatLeaf(hv, key, val))
		added = true
	} else {
		// This is the condition that allows collision leafs to exist at a level
		// less than maxDepth. I don't know if I want to allow this...
		if leaf.Hash() == hv {
			var newLeaf leafI
			newLeaf, added = leaf.put(key, val)
			curTable.replace(idx, newLeaf)
		} else {
			var t = h.createTable(depth+1, leaf, newFlatLeaf(hv, key, val))
			curTable.replace(idx, t)
			added = true
		}
	}

	if added {
		h.nentries++
	}

	return h, added
}

// Del searches the HamtTransient for the key argument and returns three
// values: a Hamt data structure, a value, and a bool.
//
// If the key was found, then the bool returned is true and the value is the
// value related to that key.
//
// If key was not found, then the bool returned is false and the value is
// nil.
//
// In either case, the Hamt value is the original HamtTransient pointer as a
// Hamt interface.
func (h *HamtTransient) Del(key []byte) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	//key = copyKey(key)

	var hv = calcHashVal(key)
	var path, leaf, idx = h.find(hv)

	var curTable = path.pop()
	var depth = uint(path.len())

	if leaf == nil {
		return h, nil, false
	}

	var newLeaf, val, deleted = leaf.del(key)

	if !deleted {
		return h, nil, false
	}

	h.nentries--

	if newLeaf != nil { //leaf was a CollisionLeaf
		curTable.replace(idx, newLeaf)
	} else { //leaf was a FlatLeaf
		curTable.remove(idx)

		// Side-Effects of removing an KeyVal from the table
		if curTable != &h.root {
			switch {
			// if no entries left in table need to colapse down to parent
			case curTable.nentries() == 1:
				var lastNode = curTable.entries()[0].node
				if _, isLeaf := lastNode.(leafI); isLeaf {
					var parentTable = path.peek()
					var parentIdx = hv.Index(depth - 1)
					parentTable.replace(parentIdx, lastNode)
				}

				// else check if downgrade allowed and required
			case h.grade && curTable.nentries() == DowngradeThreshold:
				//when nentries is decr'd it will be <DowngradeThreshold
				var newTable = downgradeToSparseTable(
					curTable.Hash(), depth, curTable.entries())
				var parentTable = path.peek()
				var parentIdx = hv.Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
		}
	}

	return h, val, deleted
}

// String returns a simple string representation of the HamtTransient data
// structure.
func (h *HamtTransient) String() string {
	return "HamtTransient{" + h.hamtBase.String() + "}"
}

// LongString returns a complete recusive listing of the entire HamtTransient
// data structure.
func (h *HamtTransient) LongString(indent string) string {
	return "HamtTransient{\n" + indent + h.hamtBase.LongString(indent) + "\n}"
}

// Visit walks the Hamt executing the VisitFn then recursing into each of
// the subtrees in order. It returns the maximum table depth it reached in
// any branch.
func (h *HamtTransient) visit(fn visitFn) uint {
	return h.hamtBase.visit(fn)
}

// Stats walks the Hamt using Visit and populates a Stats data struture which
// it return.
func (h *HamtTransient) Stats() *Stats {
	return h.hamtBase.Stats()
}

// Iter returns an IterFunc to be called repeatedly to iterate over the Hamt.
// No modifications should happend during the lifetime of the iterator. For
// HamtFunctional this is not a problem, but for HamtTransient this constaint
// is up to the Library user.
//
//    var next = h.Iter()
//    for kv, ok:= next(); ok; kv, ok = next() {
//        doSomething(kv)
//    }
//
func (h *HamtTransient) Iter() IterFunc {
	return h.hamtBase.Iter()
}

// IterChan returns a readable channel. Calls to this method spawn an
// underlying goroutine that feeds the returned channel.
//
// The chanBufLen argument allows you to set the size of the channel's buffer
// for faster iteration.
//
// The context argument is allowed to be nil.
//
// The underlying goroutine is leaked if the iterator channel is not read till
// it is exhausted and the context is not canceled.
//
//    var ctx, cancel = context.WithCancel(context.Background())
//    defer cancel()
//    var iterChan = h.IterChan(20, ctx)
//    for kv := range iterChan {
//        if shouldStop(kv) {
//            break //would leak the goroutine except for the deferred cancel
//        }
//    }
//
func (h *HamtTransient) IterChan(
	chanBufLen int,
	ctx context.Context,
) <-chan KeyVal {
	return h.hamtBase.IterChan(chanBufLen, ctx)
}

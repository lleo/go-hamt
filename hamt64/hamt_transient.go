package hamt64

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

// NewTransient constructs a new HamtTransient data structure.
//
// The tblOpt argument is the table option defined by the constants
// HybridTables, SparseTables, xor FixedTables.
//
func NewTransient(tblOpt int) *HamtTransient {
	var h = new(HamtTransient)

	h.hamtBase.init(tblOpt)

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
	nh.nograde = h.nograde
	nh.startFixed = h.startFixed
	return nh
}

// Get retrieves the value related to the key in the HamtTransient
// data structure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtTransient data structure.
func (h *HamtTransient) Get(key KeyI) (interface{}, bool) {
	return h.hamtBase.Get(key)
}

// Put stores a new (key,value) pair in the HamtTransient data structure. It
// returns a bool indicating if a new pair were added or if the value replaced
// the value in a previously stored (key,value) pair. Either way it returns and
// new HamtTransient data structure containing the modification.
func (h *HamtTransient) Put(key KeyI, val interface{}) (Hamt, bool) {
	// Doing this in newFlatLeaf() and leafI.put().

	var hv = key.Hash()
	var path, leaf, idx = h.find(hv)

	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool

	if leaf == nil {
		//check if upgrading allowed & if it is required
		if !h.nograde && curTable != &h.root &&
			(curTable.nentries()+1) == UpgradeThreshold {
			var newTable = upgradeToFixedTable(
				curTable.Hash(), depth, curTable.entries())

			var parentTable = path.peek()
			var parentIdx = hv.Index(depth - 1)
			parentTable.replace(parentIdx, newTable)

			curTable = newTable
		}
		curTable.insert(idx, newFlatLeaf(key, val))
		added = true
	} else {
		// This is the condition that allows collision leafs to exist at a level
		// less than maxDepth. I don't know if I want to allow this...
		if leaf.Hash() == hv {
			var newLeaf leafI
			newLeaf, added = leaf.put(key, val)
			curTable.replace(idx, newLeaf)
		} else {
			var t = h.createTable(depth+1, leaf, newFlatLeaf(key, val))
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
func (h *HamtTransient) Del(key KeyI) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var hv = key.Hash()
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
			case !h.nograde && curTable.nentries() == DowngradeThreshold:
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

// walk traverses the Trie in pre-order traversal. For a Trie this is also a
// in-order traversal of all leaf nodes.
//
// walk returns false if the traversal stopped early.
func (h *HamtTransient) walk(fn visitFn) bool {
	return h.hamtBase.walk(fn)
}

// Range executes the given function for every KeyVal pair in the Hamt. KeyVal
// pairs are visited in a seeminly random order.
//
// Note: we say "seemingly random order", becuase there is a predictable order
// based on the hash value of the Keys and the insertion order of the KeyVal
// pairs, so you cannot reley on the "randomness" of the order of KeyVal pairs.
func (h *HamtTransient) Range(fn func(KeyI, interface{}) bool) {
	h.hamtBase.Range(fn)
}

// Stats walks the Hamt in a pre-order traversal and populates a Stats data
// struture which it returns.
func (h *HamtTransient) Stats() *Stats {
	return h.hamtBase.Stats()
}

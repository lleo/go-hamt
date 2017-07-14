package hamt64

// HamtTransient is the datastructure which the Transient Hamt methods are
// called upon. In fact it is identical to the HamtFunctional datastructure and
// all the table and leaf datastructures it uses are the same ones used by the
// HamtTransient implementation. It is its own type so that the methods it calls
// are the transient version of the hamt64.Hamt interface.
//
// The Transient version of the Hamt datastructure, does all modifications
// in-place. So sharing this datastruture between threads is NOT safe unless
// you were to implement a locking stategy CORRECTLY.
type HamtTransient struct {
	common
}

// NewTransient constructs a new HamtTransient datastructure based on the opt
// argument.
func NewTransient(opt int) *HamtTransient {
	var h = new(HamtTransient)

	h.common.init(opt)

	return h
}

// IsEmpty simply returns if the HamtTransient datastucture has no entries.
func (h *HamtTransient) IsEmpty() bool {
	return h.common.IsEmpty()
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtTransient datastructure.
func (h *HamtTransient) Nentries() uint {
	return h.common.Nentries()
}

// ToFunctional creates a HamtFunctional data structure and copies the values
// stored in the HamtTransient data structure over to the HamtFunctional
// data structure. In the case of root table it does a deep copy. Finnally, it
// returns a pointer to the HamtFunctional data structure as a hamt64.Hamt
// interface.
//
// If you are confident that modifications to the original HamtTransient would
// not impact the HamtFunctional data structure (eg. you no longer used the
// previous HamtTransient data structures), then you can simply recast a
// *HamtTransient to *HamtFunctional.
//
// The reason for ToFunctional() is to do a deep copy of all the data
// structures involved in the HamtFunctional. Of course, this can be very
// expensive.
func (h *HamtTransient) ToFunctional() Hamt {
	var nh = new(HamtFunctional)
	nh.root = h.root.deepCopy()
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
	//return &HamtFunctional{
	//	common{
	//		root:       h.root,
	//		nentries:   h.nentries,
	//		grade:      h.grade,
	//		startFixed: h.startFixed,
	//	},
	//}
}

// ToTransient does nothing to a HamtTransient datastructure. This method only
// returns the HamtTransient datastructure pointer as a hamt64.Hamt interface.
func (h *HamtTransient) ToTransient() Hamt {
	return h
}

// DeepCopy() copies the HamtTransient datastructure and every table it contains
// recursively. This is expensive, but usefull, if you want to use ToTransient()
// and ToFunctional().
func (h *HamtTransient) DeepCopy() Hamt {
	var nh = new(HamtTransient)
	nh.root = h.root.deepCopy()
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

// Get retrieves the value related to the key in the HamtTransient
// datastructure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtTransient datastructure.
func (h *HamtTransient) Get(bs []byte) (interface{}, bool) {
	return h.common.Get(bs)
}

// Put stores a new (key,value) pair in the HamtTransient datastructure. It
// returns a bool indicating if a new pair were added or if the value replaced
// the value in a previously stored (key,value) pair. Either way it returns and
// new HamtTransient datastructure containing the modification.
func (h *HamtTransient) Put(bs []byte, v interface{}) (Hamt, bool) {
	var k = newKey(bs)

	if h.IsEmpty() {
		h.root = h.createRootTable(newFlatLeaf(k, v))
		h.nentries++
		return h, true
	}

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool

	if leaf == nil {
		//check if upgrading allowed & if it is required
		if h.grade && (curTable.nentries()+1) == UpgradeThreshold {
			var newTable tableI
			newTable = upgradeToFixedTable(
				curTable.Hash(), depth, curTable.entries())
			if curTable == h.root {
				h.root = newTable
			} else {
				var parentTable = path.peek()
				var parentIdx = k.Hash().Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
			curTable = newTable
		}
		curTable.insert(idx, newFlatLeaf(k, v))
		added = true
	} else {
		// This is the condition that allows collision leafs to exist at a level
		// less than MaxDepth. I don't know if I want to allow this...
		if leaf.Hash() == k.Hash() {
			var newLeaf leafI
			// There are four possibilities here:
			// if leaf isa collision leaf
			//   k is identical to one of the kv pairs in collision leaf; hence
			//     we replace that ones value and added = false
			//   k is unique in the collision leaf and the kv pair is added;
			//     this is very rare; the underlying key basis is different but
			//     the Hash is identical.
			// if leaf isa flat leaf
			//   k is identical to the flat leaf's key; hence the value is
			//     replaced and added == false
			//   k is not identical to the flat leaf's key; and a collision leaf
			//     is created and added == true; again this is very rare; the
			//     underlying key basis is different but the Hash is identical
			newLeaf, added = leaf.put(k, v)
			curTable.replace(idx, newLeaf)
		} else {
			var tmpTable = h.createTable(depth+1, leaf, newFlatLeaf(k, v))
			curTable.replace(idx, tmpTable)
			added = true
		}
	}

	if added {
		h.nentries++
	}

	return h, added
}

// Del searches the HamtTransient for the key argument and returns three
// values: a Hamt datastuture, a value, and a bool. If the key was found then
// the bool returned is true and the value is the value related to that key and
// the returned Hamt is a new HamtTransient datastructure without. If the
// (key, value) pair. If key was not found, then the bool is false, the value is
// nil, and the Hamt value is the original HamtTransient datastructure.
func (h *HamtTransient) Del(bs []byte) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var k = newKey(bs)

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())

	if leaf == nil {
		return h, nil, false
	}

	var newLeaf, val, deleted = leaf.del(k)

	if !deleted {
		return h, nil, false
	}

	h.nentries--

	if newLeaf != nil { //leaf was a CollisionLeaf
		curTable.replace(idx, newLeaf)
	} else { //leaf was a FlatLeaf
		curTable.remove(idx)

		// Side-Effects of removing an iKeyVal from the table
		switch {
		// if no entries left in table need to colapse down to parent
		case curTable != h.root && curTable.nentries() == 1:
			var lastNode = curTable.entries()[0].node
			if _, isLeaf := lastNode.(leafI); isLeaf {
				var parentTable = path.peek()
				var parentIdx = k.Hash().Index(depth - 1)
				parentTable.replace(parentIdx, lastNode)
			}

			// else check if downgrade allowed and required
		case h.grade && curTable.nentries() == DowngradeThreshold:
			//when nentries is decr'd it will be <DowngradeThreshold
			var newTable = downgradeToSparseTable(
				curTable.Hash(), depth, curTable.entries())
			if curTable == h.root { //aka path.len() == 0 or path.peek() == nil
				h.root = newTable
			} else {
				var parentTable = path.peek()
				var parentIdx = k.Hash().Index(depth - 1)
				parentTable.replace(parentIdx, newTable)
			}
		}
	}

	return h, val, deleted
}

// String returns a string representation of the Hamt string.
func (h *HamtTransient) String() string {
	return h.common.String()
}

// LongString returns a complete listing of the entire Hamt data structure.
func (h *HamtTransient) LongString(indent string) string {
	return h.common.LongString(indent)
}

// Visit walks the Hamt executing the VisitFn then recursing into each of
// the subtrees in order. It returns the maximum table depth it reached in
// any branch.
func (h *HamtTransient) Visit(fn visitFn, arg interface{}) uint {
	return h.common.Visit(fn, arg)
}

// Count walks the Hamt using Visit and populates a Count data struture which
// it return.
func (h *HamtTransient) Count() (uint, *Counts) {
	return h.common.Count()
}

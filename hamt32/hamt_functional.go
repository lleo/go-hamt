package hamt32

// HamtFunctional is the data structure which the Funcitonal Hamt methods are
// called upon. In fact it is identical to the HamtTransient data structure and
// all the table and leaf data structures it uses are the same ones used by the
// HamtTransient implementation. It is its own type so that the methods it calls
// are the functional version of the hamt32.Hamt interface.
//
// Basically the functional versions implement a copy-on-write inmplementation
// of Put() and Del(). The original HamtFuncitonal isn't modified and Put()
// and Del() return a slightly modified copy of the HamtFunctional
// data structure. So sharing this data structure between threads is safe.
type HamtFunctional struct {
	hamtBase
}

// NewFunctional constructs a new HamtFunctional data structure based on the opt
// argument.
func NewFunctional(opt int) *HamtFunctional {
	var h = new(HamtFunctional)

	h.hamtBase.init(opt)

	return h
}

// IsEmpty simply returns if the HamtFunctional data structure has no entries.
func (h *HamtFunctional) IsEmpty() bool {
	return h.hamtBase.IsEmpty()
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtFunctional data structure.
func (h *HamtFunctional) Nentries() uint {
	return h.hamtBase.Nentries()
}

// ToFunctional does nothing to a HamtFunctional data structure. This method
// only returns the HamtFunctional data structure pointer as a hamt32.Hamt
// interface.
func (h *HamtFunctional) ToFunctional() Hamt {
	return h
}

// ToTransient creates a HamtTransient data structure and copies the values
// stored in the HamtFunctional data structure over to the HamtTransient data
// structure. In the case of the root table it does a deep copy. Finally, it
// returns a pointer to the HamtTransient data structure as a hamt32.Hamt
// interface.
//
// If you are confident that modifications to the new HamtTransient would not
// impact the original HamtFunctional data structures it was generated from (eg.
// you no longer used the previous HamtFunctional data structures), then you can
// simply recast a *HamtFunctional to *HamtTransient.
//
// The reason for ToTransient() is to do a deep copy of all the data structures
// involved in the HamtFunctional. Of course, this can be very expensive.
func (h *HamtFunctional) ToTransient() Hamt {
	var nh = new(HamtTransient)
	nh.root = *h.root.deepCopy().(*fixedTable)
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

// DeepCopy() copies the HamtFunctional data structure and every table it
// contains recursively. This method gets more expensive the deeper the Hamt
// becomes.
func (h *HamtFunctional) DeepCopy() Hamt {
	var nh = new(HamtFunctional)
	nh.root = *h.root.deepCopy().(*fixedTable)
	nh.nentries = h.nentries
	nh.grade = h.grade
	nh.startFixed = h.startFixed
	return nh
}

// persist() is ONLY called on a fresh copy of the current Hamt.
// Hence, modifying it is allowed.
func (h *HamtFunctional) persist(oldTable, newTable tableI, path tableStack) {
	// Removed the case where path.len() == 0 on the first call to nh.perist(),
	// because that case is handled in Put & Del now. It is handled in Put & Del
	// because otherwise we were allocating an extraneous fixedTable for the
	// old h.root.
	_ = assertOn && assert(path.len() == 0,
		"This case should be handled directly in Put & Del.")

	var depth = uint(path.len()) //guaranteed depth > 0
	var parentDepth = depth - 1

	var parentIdx = oldTable.Hash().Index(parentDepth)

	var oldParent = path.pop()

	var newParent tableI
	if path.len() == 0 {
		// This condition and the last if path.len() > 0; shaves off one call
		// to persist and one fixed table allocation (via oldParent.copy()).
		h.root = *oldParent.(*fixedTable)
		newParent = &h.root
	} else {
		newParent = oldParent.copy()
	}

	if newTable == nil {
		newParent.remove(parentIdx)
	} else {
		newParent.replace(parentIdx, newTable)
	}

	if path.len() > 0 {
		h.persist(oldParent, newParent, path)
	}

	return
}

// Get retrieves the value related to the key in the HamtFunctional
// data structure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtFunctional data structure.
func (h *HamtFunctional) Get(bs []byte) (interface{}, bool) {
	return h.hamtBase.Get(bs)
}

// Put stores a new (key,value) pair in the HamtFunctional data structure. It
// returns a bool indicating if a new pair was added (true) or if the value
// replaced (false). Either way it returns a new HamtFunctional data structure
// containing the modification.
func (h *HamtFunctional) Put(key []byte, v interface{}) (Hamt, bool) {
	var nh = new(HamtFunctional)
	*nh = *h

	var k = newKey(key)

	var path, leaf, idx = h.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())

	var added bool

	if curTable == &h.root {
		//copying all h.root into nh.root already done in *nh = *h
		if leaf == nil {
			nh.root.insert(idx, newFlatLeaf(k, v))
			added = true
		} else {
			var node nodeI
			if leaf.Hash() == k.Hash() {
				node, added = leaf.put(k, v)
			} else {
				node = nh.createTable(depth+1, leaf, newFlatLeaf(k, v))
				added = true
			}

			nh.root.replace(idx, node)
		}
	} else {
		var newTable tableI

		if leaf == nil {
			if nh.grade && (curTable.nentries()+1) == UpgradeThreshold {
				newTable = upgradeToFixedTable(
					curTable.Hash(), depth, curTable.entries())
			} else {
				newTable = curTable.copy()
			}

			newTable.insert(idx, newFlatLeaf(k, v))
			added = true
		} else {
			newTable = curTable.copy()

			var node nodeI
			if leaf.Hash() == k.Hash() {
				node, added = leaf.put(k, v)
			} else {
				node = nh.createTable(depth+1, leaf, newFlatLeaf(k, v))
				added = true
			}

			newTable.replace(idx, node)
		}

		nh.persist(curTable, newTable, path)
	}

	if added {
		nh.nentries++
	}

	return nh, added
}

// Del searches the HamtFunctional for the key argument and returns three
// values: a Hamt interface, a value, and a bool.
//
// If the key was found then the bool returned is true and the value is the
// value related to that key and the returned Hamt is the new HamtFunctional
// data structure pointer.
//
// If key was not found, then the bool is false, the value is nil, and the Hamt
// value is the original HamtFunctional data structure pointer.
func (h *HamtFunctional) Del(key []byte) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var k = newKey(key)
	var path, leaf, idx = h.find(k)

	if leaf == nil {
		return h, nil, false
	}

	var newLeaf, val, deleted = leaf.del(k)

	if !deleted {
		return h, nil, false
	}

	var curTable = path.pop()
	var depth = uint(path.len())

	var nh = new(HamtFunctional)
	*nh = *h

	nh.nentries--

	if curTable == &h.root {
		//copying all h.root into nh.root already done in *nh = *h
		if newLeaf == nil { //leaf was a FlatLeaf
			nh.root.remove(idx)
		} else { //leaf was a CollisionLeaf
			nh.root.replace(idx, newLeaf)
		}
	} else {
		var newTable = curTable.copy()

		if newLeaf == nil { //leaf was a FlatLeaf
			newTable.remove(idx)

			// Side-Effects of removing a iKeyVal from the table
			var nents = newTable.nentries()
			switch {
			case nents == 0:
				newTable = nil
			case h.grade && nents == DowngradeThreshold:
				newTable = downgradeToSparseTable(
					newTable.Hash(), depth, newTable.entries())
			}
		} else { //leaf was a CollisionLeaf
			newTable.replace(idx, newLeaf)
		}

		nh.persist(curTable, newTable, path)
	}

	return nh, val, deleted
}

// String returns a simple string representation of the HamtFunctional data
// structure.
func (h *HamtFunctional) String() string {
	return "HamtFunctional{" + h.hamtBase.String() + "}"
}

// LongString returns a complete recusive listing of the entire HamtFunctional
// data structure.
func (h *HamtFunctional) LongString(indent string) string {
	return "HamtFunctional{\n" + indent + h.hamtBase.LongString(indent) + "\n}"
}

// Visit walks the Hamt executing the VisitFn then recursing into each of
// the subtrees in order. It returns the maximum table depth it reached in
// any branch.
func (h *HamtFunctional) visit(fn visitFn) uint {
	return h.hamtBase.visit(fn)
}

// Stats walks the Hamt using Visit and populates a Stats data struture which
// it return.
func (h *HamtFunctional) Stats() *Stats {
	return h.hamtBase.Stats()
}

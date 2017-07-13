package hamt64

// HamtFunctional is the datastructure which the Funcitonal Hamt methods are
// called upon. In fact it is identical to the HamtTransient datastructure and
// all the table and leaf datastructures it uses are the same ones used by the
// HamtTransient implementation. It is its own type so that the methods it calls
// are the functional version of the hamt64.Hamt interface.
//
// Basically the functional versions implement a copy-on-write inmplementation
// of Put() and Del(), to the original HamtFuncitonal isn't modified and Put()
// and Del() return a slightly modified copy of the HamtFunctional
// datastructure. So sharing this datastructure between threads is safe.
type HamtFunctional struct {
	Common
}

// NewFunctional constructs a new HamtFunctional datastructure based on the opt
// argument.
func NewFunctional(opt int) *HamtFunctional {
	var h = new(HamtFunctional)

	switch opt {
	case HybridTables:
		h.grade = true
		h.compinit = true
	case SparseTables:
		h.grade = false
		h.compinit = true
	case FixedTables:
		fallthrough
	default:
		h.grade = false
		h.compinit = false
	}

	return h
}

// IsEmpty simply returns if the HamtFunctional datastucture has no entries.
func (h *HamtFunctional) IsEmpty() bool {
	return h.Common.IsEmpty()
}

// Nentries return the number of (key,value) pairs are stored in the
// HamtFunctional datastructure.
func (h *HamtFunctional) Nentries() uint {
	return h.Common.Nentries()
}

// ToFunctional does nothing to a HamtFunctional datastructure. This method only
// returns the HamtFunctional datastructure pointer as a hamt64.Hamt interface.
func (h *HamtFunctional) ToFunctional() Hamt {
	return h
}

// ToTransient creates a HamtTransient datastructure and simply copies the
// values stored in the HamtFunctional datastructure over to the HamtTransient
// datastructure, then it returns a pointer to the HamtTransient datastructure
// as a hamt64.Hamt interface.
//
// WARNING: given that ToTransient() just copies pointers to a new
// HamtFunctional datastructure, ANY modification of the new HamtTransient
// datastruture will change the previous HamtFunctional and any preceding
// HamtFunctional datastruture that share some of the same tables.
//
// If you use the previous HamtFunctional datastructures IN ANY WAY this
// convertion is mustly useless.
//
// The only way to avoid having the new HamtTransient from modifying the
// original HamtFunctional is to first perform a DeepCopy()
func (h *HamtFunctional) ToTransient() Hamt {
	return &HamtTransient{
		Common{
			root:     h.root,
			nentries: h.nentries,
			grade:    h.grade,
			compinit: h.compinit,
		},
	}
}

// DeepCopy() copies the HamtFunctional datastructure and every table it
// contains recursively. This is expensive, but usefull, if you want to use
// ToTransient() and ToFunctional().
func (h *HamtFunctional) DeepCopy() Hamt {
	return h.Common.DeepCopy()
}

// persist() is ONLY called on a fresh copy of the current Hamt.
// Hence, modifying it is allowed.
func (h *HamtFunctional) persist(oldTable, newTable tableI, path tableStack) {
	if h.IsEmpty() {
		h.root = newTable
		return
	}

	if oldTable == h.root {
		h.root = newTable
		return
	}

	var depth = uint(path.len())
	var parentDepth = depth - 1

	var parentIdx = oldTable.Hash().Index(parentDepth)

	var oldParent = path.pop()
	var newParent tableI = oldParent.copy()

	if newTable == nil {
		newParent.remove(parentIdx)
	} else {
		newParent.replace(parentIdx, newTable)
	}

	h.persist(oldParent, newParent, path) //recurses at most MaxDepth-1 times

	return
}

// Get retrieves the value related to the key in the HamtFunctional
// datastructure. It also return a bool to indicate the value was found. This
// allows you to store nil values in the HamtFunctional datastructure.
func (h *HamtFunctional) Get(bs []byte) (interface{}, bool) {
	return h.Common.Get(bs)
}

// Put stores a new (key,value) pair in the HamtFunctional datastructure. It
// returns a bool indicating if a new pair were added or if the value replaced
// the value in a previously stored (key,value) pair. Either way it returns and
// new HamtFunctional datastructure containing the modification.
func (h *HamtFunctional) Put(bs []byte, v interface{}) (Hamt, bool) {
	var nh = new(HamtFunctional)
	*nh = *h

	var k = newKey(bs)

	if nh.IsEmpty() {
		nh.root = nh.createRootTable(newFlatLeaf(k, v))
		nh.nentries++
		return nh, true
	}

	var path, leaf, idx = nh.find(k)

	var curTable = path.pop()
	var depth = uint(path.len())
	var added bool

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
		if leaf.Hash() == k.Hash() {
			var newLeaf leafI
			newLeaf, added = leaf.put(k, v)
			newTable.replace(idx, newLeaf)
		} else {
			var tmpTable = nh.createTable(depth+1, leaf, newFlatLeaf(k, v))
			newTable.replace(idx, tmpTable)
			added = true
		}
	}

	if added {
		nh.nentries++
	}

	nh.persist(curTable, newTable, path)

	return nh, added
}

// Del searches the HamtFunctional for the key argument and returns three
// values: a Hamt datastuture, a value, and a bool. If the key was found then
// the bool returned is true and the value is the value related to that key and
// the returned Hamt is a new HamtFunctional datastructure without. If the
// (key, value) pair. If key was not found, then the bool is false, the value is
// nil, and the Hamt value is the original HamtFunctional datastructure.
func (h *HamtFunctional) Del(bs []byte) (Hamt, interface{}, bool) {
	if h.IsEmpty() {
		return h, nil, false
	}

	var k = newKey(bs)
	var path, leaf, idx = h.find(k)

	var curTable = path.pop()

	if leaf == nil {
		return h, nil, false
	}

	var newLeaf, val, deleted = leaf.del(k)

	if !deleted {
		return h, nil, false
	}

	var depth = uint(path.len())
	var newTable tableI = curTable.copy()
	if newLeaf != nil { //leaf was a CollisionLeaf
		newTable.replace(idx, newLeaf)
	} else { //leaf was a FlatLeaf
		newTable.remove(idx)

		// Side-Effects of removing a KeyVal from the table
		switch {
		case newTable.nentries() == 0:
			newTable = nil
		case h.grade && newTable.nentries() == DowngradeThreshold:
			newTable = downgradeToSparseTable(
				newTable.Hash(), depth, newTable.entries())
		}
	}

	var nh = new(HamtFunctional)
	*nh = *h

	nh.nentries--

	nh.persist(curTable, newTable, path)

	return nh, val, deleted
}

// String returns a string representation of the HamtFunctional stastructure.
// Secifically it returns a representation of the datastructure with the
// nentries value of Nentries() and a representation of the root table.
func (h *HamtFunctional) String() string {
	return h.Common.String()
}

// LongString returns a complete listing of the entire Hamt data structure
// recursively indented..
func (h *HamtFunctional) LongString(indent string) string {
	return h.Common.LongString(indent)
}

func (h *HamtFunctional) Visit(fn visitFn, arg interface{}) uint {
	return h.Common.Visit(fn, arg)
}

func (h *HamtFunctional) Count() (uint, *Counts) {
	return h.Common.Count()
}

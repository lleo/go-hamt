/*

 */
package hamt32

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/key"
)

// NBITS constant is the number of bits(5) a 30bit hash value is split into
// to provied the index of a HAMT.
const NBITS uint = 5

// MAXDEPTH constant is the maximum depth(5) of NBITS values that constitute
// the path in a HAMT, from [0..MAXDEPTH] for a total of MAXDEPTH+1(6) levels.
// NBITS*(MAXDEPTH+1) == HASHBITS (ie 5*(5+1) == 30).
const MAXDEPTH uint = 5

// TABLE_CAPACITY constant is the number of table entries in a each node of
// a HAMT datastructure; its value is 2^5 == 32.
const TABLE_CAPACITY uint = uint(1 << NBITS)

// DOWNGRADE_THRESHOLD constant is the number of nodes a fullTable has shrunk to,
// before it is converted to a compressedTable.
const DOWNGRADE_THRESHOLD uint = TABLE_CAPACITY / 8

// UPGRADE_THRESHOLD constan is the number of nodes a compressedTable has grown to,
// before it is converted to a fullTable.
const UPGRADE_THRESHOLD uint = TABLE_CAPACITY / 2

func indexMask(depth uint) uint32 {
	return uint32(uint8(1<<NBITS)-1) << (depth * NBITS)
}

func index(h30 uint32, depth uint) uint {
	var idxMask = indexMask(depth)
	var idx = uint((h30 & idxMask) >> (depth * NBITS))
	return idx
}

func hashPathString(hashPath uint32, depth uint) string {
	if depth == 0 {
		return "/"
	}
	var strs = make([]string, depth)

	for d := depth; d > 0; d-- {
		var idx = index(hashPath, d-1)
		strs[d-1] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

func hash30String(h30 uint32) string {
	return hashPathString(h30, 6)
}

func hashPathMask(depth uint) uint32 {
	return uint32(1<<(depth*NBITS)) - 1
}

func buildHashPath(hashPath uint32, idx, depth uint) uint32 {
	return hashPath | uint32(idx<<(depth*NBITS))
}

type keyVal struct {
	key key.Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

const (
	HYBRID = iota
	COMPONLY
	FULLONLY
)

var OPTIONS = make(map[int]string, 3)

func init() {
	OPTIONS[0] = "HYBRID"
	OPTIONS[1] = "COMPONLY"
	OPTIONS[2] = "FULLONLY"
}

type Hamt struct {
	root            tableI
	nentries        int
	grade, fullinit bool
}

// Create a new hamt32.Hamt datastructure with the table options set to either
//   hamt32.HYBRID - initially start out with compressedTable, but when the table is
//                   half full upgrade to fullTable. If a fullTable shrinks to
//                   TABLE_CAPACITY/8(4) entries downgrade to compressed table.
//   hamt32.COMPONLY - Only use compressedTable no up/downgrading to/from fullTable.
//                     This uses the least amount of space.
//   hamt32.FULLONLY - Only use fullTable no up/downgrading from/to compressedTables.
//                     This is the fastest performance.
func New(opt int) *Hamt {
	var h = new(Hamt)
	if opt == COMPONLY {
		h.grade = false
		h.fullinit = false
	} else if opt == FULLONLY {
		h.grade = false
		h.fullinit = true
	} else /* opt == HYBRID */ {
		h.grade = true
		h.fullinit = false
	}
	return h
}

func (h *Hamt) IsEmpty() bool {
	return h.root == nil
}

func (h *Hamt) Get(k key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h30 = k.Hash30()

	var curTable = h.root //ISA tableI

	for depth := uint(0); depth <= MAXDEPTH; depth++ {
		var idx = index(h30, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			break
		}

		if leaf, isLeaf := curNode.(leafI); isLeaf {
			var val, found = leaf.get(k)
			return val, found
		}

		//else curNode MUST BE A tableI
		curTable = curNode.(tableI)
	}
	// curNode == nil || depth > MAXDEPTH

	return nil, false
}

func (h *Hamt) Put(k key.Key, v interface{}) bool {
	//var newLeaf = newFlatLeaf(k, v)
	var depth uint = 0
	var hashPath uint32 = 0
	var inserted = true

	if h.IsEmpty() {
		h.root = h.newRootTable(depth, hashPath, newFlatLeaf(k, v))
		h.nentries++
		return inserted
	}

	var path = newPathT()
	var curTable = h.root

	for depth = 0; depth <= MAXDEPTH; depth++ {
		var idx = index(k.Hash30(), depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			curTable.set(idx, newFlatLeaf(k, v))
			h.nentries++

			// upgrade?
			if h.grade {
				_, isCompressedTable := curTable.(*compressedTable)
				if isCompressedTable && curTable.nentries() >= UPGRADE_THRESHOLD {
					curTable = upgradeToFullTable(hashPath, curTable.entries())
					if depth == 0 {
						h.root = curTable
					} else {
						parentTable := path.peek()
						parentIdx := index(k.Hash30(), depth-1)
						parentTable.set(parentIdx, curTable)
					}
				}
			}

			return true //inserted
		}

		if curLeaf, isLeaf := curNode.(leafI); isLeaf {
			if curLeaf.Hash30() == k.Hash30() {
				// Accumulate collisionLeaf
				insLeaf, inserted := curLeaf.put(k, v)
				if inserted {
					curTable.set(idx, insLeaf)
					h.nentries++
				}
				return inserted
			}

			if depth == MAXDEPTH {
				// this test should be delete cuz it is logically impossible
				if curLeaf.Hash30() != k.Hash30() {
					// This should not happen cuz we had to walk up MAXDEPTH
					// levels to get here.
					panic("WTF!!!")
				}

				// Accumulate collisionLeaf
				insLeaf, inserted := curLeaf.put(k, v)
				if inserted {
					curTable.set(idx, insLeaf)
					h.nentries++
				}
				return inserted
			}

			hashPath = buildHashPath(hashPath, idx, depth)
			var collisionTable = h.newTable(depth+1, hashPath, curLeaf, newFlatLeaf(k, v))
			curTable.set(idx, collisionTable)
			h.nentries++

			return true
		}

		hashPath = buildHashPath(hashPath, idx, depth)
		path.push(curTable)
		curTable = curNode.(tableI)
	}

	//log.Println(path)
	//log.Printf("k=%s, v=%v", k, v)

	panic("WTF!!")
}

func (h *Hamt) Del(k key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h30 = k.Hash30()
	var depth uint
	var hashPath uint32

	var path = newPathT()
	var curTable = h.root

	for depth = 0; depth <= MAXDEPTH; depth++ {
		var idx = index(h30, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			return nil, false
		}

		if curLeaf, isLeaf := curNode.(leafI); isLeaf {
			v, delLeaf, deleted := curLeaf.del(k)
			if !deleted {
				return nil, false
			}
			// else a leaf key/value was deleted
			h.nentries--

			// If curLeaf was a collisionLeaf,
			// then delLeaf is either a slimmed down collisionLeaf or a flatLeaf.
			// If curLeaf was a flatLeaf then delLeaf is nil.
			curTable.set(idx, delLeaf)

			// downgrade?
			if h.grade {
				if delLeaf == nil {
					_, isFullTable := curTable.(*fullTable)
					if isFullTable && curTable.nentries() <= DOWNGRADE_THRESHOLD {
						curTable = downgradeToCompressedTable(hashPath, curTable.entries())
						if depth == 0 {
							h.root = curTable
						} else {
							parentTable := path.peek()
							parentIdx := index(h30, depth-1)
							parentTable.set(parentIdx, curTable)
						}
					}
				}
			}
			// If curTable has only one entry and that entry is a leafI,
			// then collapse that leafI down to the position curTable holds
			// in the parent Table; repeat test and collapse for parent table.

			// Identical for conditionals !!!
			//  curTable.nentries() == 1 && curTable != h.root
			//  curTable.nentries() == 1 && len(path) > 0
			for curTable.nentries() == 1 && depth > 0 {
				// _ = ASSERT && Assert(curTable != h.root, "curTable == h.root")
				// _ = ASSERT && Assert(depth == len(path), "depth != len(path)")

				var node = curTable.entries()[0].node
				var leaf, isLeaf = node.(leafI)
				if !isLeaf {
					break
				}

				// Collapse leaf down to where curTable is in parentTable

				var parentTable = path.pop()
				depth-- // OR depth = len(path)

				//var parentIdx = index(curTable.Hash30(), depth-1)
				parentIdx := index(curTable.Hash30(), depth)
				parentTable.set(parentIdx, leaf)

				curTable = parentTable
			}

			if curTable == h.root && curTable.nentries() == 0 {
				h.root = nil
			}

			return v, true
		} //if isLeaf

		// curNode is not nil
		// curNode is not a leafI
		// curNode MUST be a tableI
		hashPath = buildHashPath(hashPath, idx, depth)
		path.push(curTable)
		curTable = curNode.(tableI)
	} //for depth loop

	//log.Printf("WTF! this should never be called; k=%s", k)
	return nil, false
}

func (h *Hamt) String() string {
	return fmt.Sprintf("Hamt{ nentries: %d, root: %s }", h.nentries, h.root.LongString("", 0))
}

func (h *Hamt) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent + fmt.Sprintf("Hamt{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //Hamt"
	} else {
		str = indent + fmt.Sprintf("Hamt{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}

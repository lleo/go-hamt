/*
Package hamt32 implements a 32 node wide Hashed Array Mapped Trie. The hash key
is 30 bits wide and broken into 6 numbers of 5 bits each. Those 5bit numbers
allows us to index into a 32 node array. Each node is either a leaf or another
32 node table. So the 30bit hash allows us to index into a B+ Tree with a
branching factor of 32 and a Maximum depth of 6.

The basic insertion operation is to calculate a 30 bit hash value from your key
(a string in the case you use hamt.StringKey), then split it into six 5bit
 numbers. These six numbers represent a path thru the tree. For each level we
use the coresponding number as an index into the 32 cell array. If the cell is
empty we create a  leaf node there. If the cell is occupide by another table
we continue walking up the tree. If the cell is occupide by a leaf we promote
that cell to a new table and put the current leaf and new one into that table
in cells corresponding to that new level. If we are at the maximun depth of
the  tree and there is already a leaf there we insert our key,value pair into
that leaf.

The retrieval operation is a simmilar tree walk guided by the six 5bit numbers
till we find a leaf with the key,value pair in it.

The deletion operation is a walk to find the key, then delete the key from the
leaf. An empty leaf is removed from it's table. If the table has only one other
leaf in that level we will remove that leaf, replace the table in it's parent
table placing that last leaf down one level.

*/
package hamt32

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/key"
)

// nBits constant is the number of bits(5) a 30bit hash value is split into
// to provied the index of a HAMT.
const nBits uint = 5

// maxDepth constant is the maximum depth(5) of nBits values that constitute
// the path in a HAMT, from [0..maxDepth] for a total of maxDepth+1(6) levels.
// nBits*(maxDepth+1) == HASHBITS (ie 5*(5+1) == 30).
const maxDepth uint = 5

// tableCapacity constant is the number of table entries in a each node of
// a HAMT datastructure; its value is 2^5 == 32.
const tableCapacity uint = uint(1 << nBits)

// downgradeThreshold constant is the number of nodes a fullTable has shrunk to,
// before it is converted to a compressedTable.
const downgradeThreshold uint = tableCapacity / 8

// upgradeThreshold constan is the number of nodes a compressedTable has grown to,
// before it is converted to a fullTable.
const upgradeThreshold uint = tableCapacity / 2

func indexMask(depth uint) uint32 {
	return uint32(uint8(1<<nBits)-1) << (depth * nBits)
}

func index(h30 uint32, depth uint) uint {
	var idxMask = indexMask(depth)
	var idx = uint((h30 & idxMask) >> (depth * nBits))
	return idx
}

func hashPathString(hashPath uint32, depth uint) string {
	if depth == 0 {
		return "/"
	}
	var strs = make([]string, depth)

	for d := uint(0); d < depth; d++ {
		var idx = index(hashPath, d)
		strs[d] = fmt.Sprintf("%02d", idx)
	}

	return "/" + strings.Join(strs, "/")
}

func hash30String(h30 uint32) string {
	return hashPathString(h30, maxDepth)
}

func hashPathMask(depth uint) uint32 {
	return uint32(1<<(depth*nBits)) - 1
}

func buildHashPath(hashPath uint32, idx, depth uint) uint32 {
	return hashPath | uint32(idx<<(depth*nBits))
}

type keyVal struct {
	key key.Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

// These constants are the three configuration options to the hamt32.New()
const (
	Hybrid = iota
	CompressedOnly
	FullOnly
)

var optionsMap = make(map[int]string, 3)

func init() {
	optionsMap[Hybrid] = "Hybrid"
	optionsMap[CompressedOnly] = "CompressedOnly"
	optionsMap[FullOnly] = "FullOnly"
}

type Hamt struct {
	root            tableI
	nentries        int
	grade, fullinit bool
}

// Create a new hamt32.Hamt datastructure with the table optionsMap set to either
//   hamt32.Hybrid - initially start out with compressedTable, but when the table is
//                   half full upgrade to fullTable. If a fullTable shrinks to
//                   tableCapacity/8(4) entries downgrade to compressed table.
//   hamt32.CompressedOnly - Only use compressedTable no up/downgrading to/from fullTable.
//                           This uses the least amount of space.
//   hamt32.FullOnly - Only use fullTable no up/downgrading from/to compressedTables.
//                     This is the fastest performance.
func New(opt int) *Hamt {
	var h = new(Hamt)
	if opt == CompressedOnly {
		h.grade = false
		h.fullinit = false
	} else if opt == FullOnly {
		h.grade = false
		h.fullinit = true
	} else /* opt == Hybrid */ {
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

	for depth := uint(0); depth <= maxDepth; depth++ {
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
	// curNode == nil || depth > maxDepth

	return nil, false
}

func (h *Hamt) Put(k key.Key, v interface{}) bool {
	var depth uint
	var hashPath uint32

	if h.IsEmpty() {
		h.root = h.newRootTable(depth, hashPath, newFlatLeaf(k, v))
		h.nentries++
		return true
	}

	var path = newPathT()
	var curTable = h.root

	for depth = 0; depth <= maxDepth; depth++ {
		var idx = index(k.Hash30(), depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			curTable.set(idx, newFlatLeaf(k, v))
			h.nentries++

			// upgrade?
			if h.grade {
				_, isCompressedTable := curTable.(*compressedTable)
				if isCompressedTable && curTable.nentries() >= upgradeThreshold {
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
				// This is a minor optimization but since these two leaves
				// will collide all the way up the to maxDepth, we can
				// choose to create the collisionLeaf hear and now.

				// Accumulate collisionLeaf
				colLeaf, inserted := curLeaf.put(k, v)
				if inserted {
					curTable.set(idx, colLeaf)
					h.nentries++
				}
				return inserted
			}

			if depth == maxDepth {
				// this test should be delete cuz it is logically impossible
				if curLeaf.Hash30() != k.Hash30() {
					// This should not happen cuz we had to walk up maxDepth
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

	for depth = 0; depth <= maxDepth; depth++ {
		var idx = index(h30, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			return nil, false
		}

		if curLeaf, isLeaf := curNode.(leafI); isLeaf {
			val, delLeaf, deleted := curLeaf.del(k)
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
					if isFullTable && curTable.nentries() <= downgradeThreshold {
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
			//
			// These are identical for conditionals:
			//  curTable != h.root AND len(path) > 0 AND depth > 0
			//
			for curTable.nentries() == 1 && depth > 0 {
				// _ = ASSERT && Assert(curTable != h.root, "curTable == h.root")
				// _ = ASSERT && Assert(depth == len(path), "depth != len(path)")

				var node = curTable.entries()[0].node
				var leaf, isLeaf = node.(leafI)
				if !isLeaf {
					// We only collapse leafs
					break
				}

				// Collapse leaf down to where curTable is in parentTable

				var parentTable = path.pop()
				depth-- // OR depth = len(path)

				parentIdx := index(curTable.Hash30(), depth)
				parentTable.set(parentIdx, leaf)

				curTable = parentTable
			}

			// TODO: I should keep this table rather than throwing it away.
			// Instead using h.root == nil to detect emptyness, we should
			// trust our accounting and use h.nentries == 0.
			if curTable == h.root && curTable.nentries() == 0 {
				h.root = nil
			}

			return val, true
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

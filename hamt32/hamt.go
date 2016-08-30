package hamt32

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/key"
)

const NBITS uint = 5
const DEPTHLIMIT uint = 6
const TABLE_CAPACITY uint = uint(1 << NBITS)

func indexMask(depth uint) uint32 {
	return uint32(uint8(1<<NBITS)-1) << (depth * NBITS)
}

func index(h30 uint32, depth uint) uint {
	var idxMask = indexMask(depth)
	var idx = uint((h30 & idxMask) >> (depth * NBITS))
	return idx
}

func hashPathEqual(depth uint, a, b uint32) bool {
	//pathMask := uint32(1<<(depth*NBITS)) - 1
	var pathMask = hashPathMask(depth)

	return (a & pathMask) == (b & pathMask)
}

func hashPathMask(depth uint) uint32 {
	return uint32(1<<(depth*NBITS)) - 1
}

func buildHashPath(hashPath uint32, idx, depth uint) uint32 {
	return hashPath | uint32(idx<<(depth*NBITS))
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

func nodeMapString(nodeMap uint32) string {
	var strs = make([]string, 4)

	var top2 = nodeMap >> 30
	strs[0] = fmt.Sprintf("%02b", top2)

	const tenBitMask uint32 = 1<<10 - 1
	for i := uint(0); i < 3; i++ {
		var tenBitVal = (nodeMap & (tenBitMask << (i * 10))) >> (i * 10)
		strs[3-i] = fmt.Sprintf("%010b", tenBitVal)
	}

	return strings.Join(strs, " ")
}

type Hamt struct {
	root     tableI
	nentries int
}

func NewHamt() *Hamt {
	var h = new(Hamt)
	return h
}

type keyVal struct {
	key key.Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

func (h *Hamt) IsEmpty() bool {
	return h.root == nil
}

func (h *Hamt) Get(k key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h30 = key.Hash30(k)

	var curTable = h.root //ISA tableI

	for depth := uint(0); depth < DEPTHLIMIT; depth++ {
		var idx = index(h30, depth)
		var curNode = curTable.get(idx) //nodeI

		if curNode == nil {
			break
		}

		if leaf, isLeaf := curNode.(leafI); isLeaf {

			if hashPathEqual(depth, h30, leaf.hash30()) {
				var v, found = leaf.get(k)
				return v, found
			}

			return nil, false
		}

		//else curNode MUST BE A tableI
		curTable = curNode.(tableI)
	}
	// curNode == nil || depth >= DEPTHLIMIT

	return nil, false
}

func (h *Hamt) Put(k key.Key, v interface{}) bool {
	var h30 = key.Hash30(k)
	var newLeaf = newFlatLeaf(h30, k, v)
	var depth uint = 0

	if h.IsEmpty() {
		h.root = newCompressedTable(depth, h30, newLeaf)
		h.nentries++
		return true
	}

	var path = newPathT()
	var hashPath uint32 = 0
	var curTable = h.root
	var inserted = true

	for depth = 0; depth < DEPTHLIMIT; depth++ {
		var idx = index(h30, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			curTable.set(idx, newLeaf)
			h.nentries++
			break //from for-loop
		}

		if oldLeaf, isLeaf := curNode.(leafI); isLeaf {
			if oldLeaf.hash30() == h30 {
				var newLeaf leafI
				newLeaf, inserted = oldLeaf.put(k, v)
				if inserted {
					curTable.set(idx, newLeaf)
				}
			} else {
				hashPath = buildHashPath(hashPath, idx, depth)
				var newLeaf = newFlatLeaf(h30, k, v)
				var collisionTable = newCompressedTable_2(depth+1, hashPath, oldLeaf, newLeaf)
				//var collisionTable = newCompressedTable_2(depth, hashPath, oldLeaf, newLeaf)
				curTable.set(idx, collisionTable)
			}
			if inserted {
				h.nentries++
			}
			break //from for-loop
		}

		hashPath = buildHashPath(hashPath, idx, depth)
		path.push(curTable)
		curTable = curNode.(tableI)
	}

	var _, isCompressedTable = curTable.(*compressedTable)
	if isCompressedTable && curTable.nentries() > TABLE_CAPACITY/2 {
		if curTable == h.root {
			curTable = upgradeToFullTable(hashPath, curTable.entries())

			h.root = curTable
		} else {
			curTable = upgradeToFullTable(hashPath, curTable.entries())

			var parentTable = path.peek()

			var parentIdx = index(curTable.hash30(), depth-1)
			parentTable.set(parentIdx, curTable)
		}
	}

	return inserted
}

func (h *Hamt) Del(k key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h30 = key.Hash30(k)
	var depth uint = 0

	var path = newPathT()
	var hashPath uint32 = 0
	var curTable = h.root

	for depth = 0; depth < DEPTHLIMIT; depth++ {
		var idx = index(h30, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			return nil, false
		}

		if oldLeaf, isLeaf := curNode.(leafI); isLeaf {
			if oldLeaf.hash30() == h30 {
				if v, newLeaf, deleted := oldLeaf.del(k); deleted {
					//newLeaf MUST BE nil or a leaf slimmer by one
					if newLeaf != oldLeaf {
						//minor optimization, cuz curTable.set() can be non-trivial
						curTable.set(idx, newLeaf)
					}
					h.nentries--

					// demote curTable if it is a fullTable && shrank to small
					var _, isFullTable = curTable.(*fullTable)
					if isFullTable && curTable.nentries() < TABLE_CAPACITY/2 {
						if curTable == h.root {
							curTable = downgradeToCompressedTable(hashPath, curTable.entries())
							h.root = curTable
						} else {
							curTable = downgradeToCompressedTable(hashPath, curTable.entries())
							var parentTable = path.peek()
							var parentIdx = index(curTable.hash30(), depth-1)
							parentTable.set(parentIdx, curTable)
						}
					}

					if curTable != h.root && curTable.nentries() == 1 {
						var node = curTable.entries()[0].node
						if leaf, isLeaf := node.(leafI); isLeaf {
							// ONLY COLLAPSE LEAVES
							for depth > 0 {
								var parentTable = path.pop()
								var parentIdx = index(curTable.hash30(), depth-1)
								parentTable.set(parentIdx, leaf)

								curTable = parentTable
								depth--

								if parentTable.nentries() > 1 {
									break
								}
							}
						}
					}

					if curTable == h.root && curTable.nentries() == 0 {
						h.root = nil
					}

					return v, true
				} //if deleted
			} //if h30 == leaf.hash30

			return nil, false
		} //if isLeaf

		hashPath = buildHashPath(hashPath, idx, depth)
		path.push(curTable)
		curTable = curNode.(tableI)
	} //for depth loop

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

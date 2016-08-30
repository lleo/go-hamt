package hamt64

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/key"
)

const NBITS uint = 5
const DEPTHLIMIT uint = 6
const TABLE_CAPACITY uint = uint(1 << NBITS)

func indexMask(depth uint) uint64 {
	return uint64(uint8(1<<NBITS)-1) << (depth * NBITS)
}

func index(h60 uint64, depth uint) uint {
	var idxMask = indexMask(depth)
	var idx = uint((h60 & idxMask) >> (depth * NBITS))
	return idx
}

func hashPathEqual(depth uint, a, b uint64) bool {
	//pathMask := uint64(1<<(depth*NBITS)) - 1
	var pathMask = hashPathMask(depth)

	return (a & pathMask) == (b & pathMask)
}

func hashPathMask(depth uint) uint64 {
	return uint64(1<<(depth*NBITS)) - 1
}

func buildHashPath(hashPath uint64, idx, depth uint) uint64 {
	return hashPath | uint64(idx<<(depth*NBITS))
}

func hashPathString(hashPath uint64, depth uint) string {
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

func hash60String(h60 uint64) string {
	return hashPathString(h60, 6)
}

func nodeMapString(nodeMap uint64) string {
	var strs = make([]string, 4)

	var top2 = nodeMap >> 60
	strs[0] = fmt.Sprintf("%02b", top2)

	const tenBitMask uint64 = 1<<10 - 1
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

	var h60 = key.Hash60(k)

	var curTable = h.root //ISA tableI

	for depth := uint(0); depth < DEPTHLIMIT; depth++ {
		var idx = index(h60, depth)
		var curNode = curTable.get(idx) //nodeI

		if curNode == nil {
			break
		}

		if leaf, isLeaf := curNode.(leafI); isLeaf {

			if hashPathEqual(depth, h60, leaf.hash60()) {
				var val, found = leaf.get(k)
				return val, found
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
	var h60 = key.Hash60(k)
	var newLeaf = newFlatLeaf(h60, k, v)
	var depth uint = 0

	if h.IsEmpty() {
		h.root = newCompressedTable(depth, h60, newLeaf)
		h.nentries++
		return true
	}

	var path = newPathT()
	var hashPath uint64 = 0
	var curTable = h.root
	var inserted = true

	for depth = 0; depth < DEPTHLIMIT; depth++ {
		var idx = index(h60, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			curTable.set(idx, newLeaf)
			h.nentries++
			break //from for-loop
		}

		if oldLeaf, isLeaf := curNode.(leafI); isLeaf {
			if oldLeaf.hash60() == h60 {
				var newLeaf leafI
				newLeaf, inserted = oldLeaf.put(k, v)
				if inserted {
					curTable.set(idx, newLeaf)
				}
			} else {
				hashPath = buildHashPath(hashPath, idx, depth)
				var newLeaf = newFlatLeaf(h60, k, v)
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

			var parentIdx = index(curTable.hash60(), depth-1)
			parentTable.set(parentIdx, curTable)
		}
	}

	return inserted
}

func (h *Hamt) Del(k key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h60 = key.Hash60(k)
	var depth uint = 0

	var path = newPathT()
	var hashPath uint64 = 0
	var curTable = h.root

	for depth = 0; depth < DEPTHLIMIT; depth++ {
		var idx = index(h60, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			return nil, false
		}

		if oldLeaf, isLeaf := curNode.(leafI); isLeaf {
			if oldLeaf.hash60() == h60 {
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
							var parentIdx = index(curTable.hash60(), depth-1)
							parentTable.set(parentIdx, curTable)
						}
					}

					if curTable != h.root && curTable.nentries() == 1 {
						var node = curTable.entries()[0].node
						if leaf, isLeaf := node.(leafI); isLeaf {
							// ONLY COLLAPSE LEAVES
							for depth > 0 {
								var parentTable = path.pop()
								var parentIdx = index(curTable.hash60(), depth-1)
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
			} //if h60 == leaf.hash60

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

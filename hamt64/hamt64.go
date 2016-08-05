package hamt64

import (
	"fmt"
	"strings"

	"github.com/lleo/go-hamt/hamt_key"
)

const NBITS64 uint = 5
const DEPTHLIMIT64 uint = 6
const TABLE_CAPACITY64 uint = uint(1 << NBITS64)

func indexMask(depth uint) uint64 {
	return uint64(uint8(1<<NBITS64)-1) << (depth * NBITS64)
}

func index(h60 uint64, depth uint) uint {
	var idxMask = indexMask(depth)
	var idx = uint((h60 & idxMask) >> (depth * NBITS64))
	return idx
}

func hashPathEqual(depth uint, a, b uint64) bool {
	//pathMask := uint64(1<<(depth*NBITS64)) - 1
	var pathMask = hashPathMask(depth)

	return (a & pathMask) == (b & pathMask)
}

func hashPathMask(depth uint) uint64 {
	return uint64(1<<(depth*NBITS64)) - 1
}

func buildHashPath(hashPath uint64, idx, depth uint) uint64 {
	return hashPath | uint64(idx<<(depth*NBITS64))
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

//type Key interface {
//	Equals(Key) bool
//	Hash60() uint64
//	Hash60() uint64
//	String() string
//}

//type Hamt interface {
//	Get(hamt_key.Key) (interface{}, bool)
//	Put(hamt_key.Key, interface{}) bool
//	Del(hamt_key.Key) (interface{}, bool)
//	String() string
//	LongString(indent string) string
//}

type Hamt64 struct {
	root     table64I
	nentries int
}

func NewHamt64() *Hamt64 {
	var h = new(Hamt64)
	return h
}

type keyVal struct {
	key hamt_key.Key
	val interface{}
}

func (kv keyVal) String() string {
	return fmt.Sprintf("keyVal{%s, %v}", kv.key, kv.val)
}

func (h *Hamt64) IsEmpty() bool {
	return h.root == nil
}

func (h *Hamt64) Get(key hamt_key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h60 = key.Hash60()

	var curTable = h.root //ISA table64I

	for depth := uint(0); depth < DEPTHLIMIT64; depth++ {
		var idx = index(h60, depth)
		var curNode = curTable.get(idx) //node64I

		if curNode == nil {
			break
		}

		if leaf, isLeaf := curNode.(leaf64I); isLeaf {

			if hashPathEqual(depth, h60, leaf.hash60()) {
				var val, found = leaf.get(key)
				return val, found
			}

			return nil, false
		}

		//else curNode MUST BE A table64I
		curTable = curNode.(table64I)
	}
	// curNode == nil || depth >= DEPTHLIMIT64

	return nil, false
}

func (h *Hamt64) Put(key hamt_key.Key, val interface{}) bool {
	var h60 = key.Hash60()
	var newLeaf = NewFlatLeaf64(h60, key, val)
	var depth uint = 0

	if h.IsEmpty() {
		h.root = newCompressedTable64(depth, h60, newLeaf)
		h.nentries++
		return true
	}

	var path = newPath64T()
	var hashPath uint64 = 0
	var curTable = h.root
	var inserted = true

	for depth = 0; depth < DEPTHLIMIT64; depth++ {
		var idx = index(h60, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			curTable.set(idx, newLeaf)
			h.nentries++
			break //from for-loop
		}

		if oldLeaf, isLeaf := curNode.(leaf64I); isLeaf {
			if oldLeaf.hash60() == h60 {
				var newLeaf leaf64I
				newLeaf, inserted = oldLeaf.put(key, val)
				if inserted {
					curTable.set(idx, newLeaf)
				}
			} else {
				hashPath = buildHashPath(hashPath, idx, depth)
				var newLeaf = NewFlatLeaf64(h60, key, val)
				var collisionTable = newCompressedTable64_2(depth+1, hashPath, oldLeaf, newLeaf)
				//var collisionTable = newCompressedTable64_2(depth, hashPath, oldLeaf, newLeaf)
				curTable.set(idx, collisionTable)
			}
			if inserted {
				h.nentries++
			}
			break //from for-loop
		}

		hashPath = buildHashPath(hashPath, idx, depth)
		path.push(curTable)
		curTable = curNode.(table64I)
	}

	var _, isCompressedTable = curTable.(*compressedTable64)
	if isCompressedTable && curTable.nentries() > TABLE_CAPACITY64/2 {
		if curTable == h.root {
			curTable = UpgradeToFullTable64(hashPath, curTable.entries())

			h.root = curTable
		} else {
			curTable = UpgradeToFullTable64(hashPath, curTable.entries())

			var parentTable = path.peek()

			var parentIdx = index(curTable.hash60(), depth-1)
			parentTable.set(parentIdx, curTable)
		}
	}

	return inserted
}

func (h *Hamt64) Del(key hamt_key.Key) (interface{}, bool) {
	if h.IsEmpty() {
		return nil, false
	}

	var h60 = key.Hash60()
	var depth uint = 0

	var path = newPath64T()
	var hashPath uint64 = 0
	var curTable = h.root

	for depth = 0; depth < DEPTHLIMIT64; depth++ {
		var idx = index(h60, depth)
		var curNode = curTable.get(idx)

		if curNode == nil {
			return nil, false
		}

		if oldLeaf, isLeaf := curNode.(leaf64I); isLeaf {
			if oldLeaf.hash60() == h60 {
				if val, newLeaf, deleted := oldLeaf.del(key); deleted {
					//newLeaf MUST BE nil or a leaf slimmer by one
					if newLeaf != oldLeaf {
						//minor optimization, cuz curTable.set() can be non-trivial
						curTable.set(idx, newLeaf)
					}
					h.nentries--

					// demote curTable if it is a fullTable64 && shrank to small
					var _, isFullTable = curTable.(*fullTable64)
					if isFullTable && curTable.nentries() < TABLE_CAPACITY64/2 {
						if curTable == h.root {
							curTable = DowngradeToCompressedTable64(hashPath, curTable.entries())
							h.root = curTable
						} else {
							curTable = DowngradeToCompressedTable64(hashPath, curTable.entries())
							var parentTable = path.peek()
							var parentIdx = index(curTable.hash60(), depth-1)
							parentTable.set(parentIdx, curTable)
						}
					}

					if curTable != h.root && curTable.nentries() == 1 {
						var node = curTable.entries()[0].node
						if leaf, isLeaf := node.(leaf64I); isLeaf {
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

					return val, true
				} //if deleted
			} //if h60 == leaf.hash60

			return nil, false
		} //if isLeaf

		hashPath = buildHashPath(hashPath, idx, depth)
		path.push(curTable)
		curTable = curNode.(table64I)
	} //for depth loop

	return nil, false
}

func (h *Hamt64) String() string {
	return fmt.Sprintf("Hamt64{ nentries: %d, root: %s }", h.nentries, h.root.LongString("", 0))
}

func (h *Hamt64) LongString(indent string) string {
	var str string
	if h.root != nil {
		str = indent + fmt.Sprintf("Hamt64{ nentries: %d, root:\n", h.nentries)
		str += indent + h.root.LongString(indent, 0)
		str += indent + "} //Hamt64"
	} else {
		str = indent + fmt.Sprintf("Hamt64{ nentries: %d, root: nil }", h.nentries)
	}
	return str
}

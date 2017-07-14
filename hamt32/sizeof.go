package hamt32

import (
	"unsafe"
)

var SizeofFixedTable = unsafe.Sizeof(fixedTable{})
var SizeofSparseTable = unsafe.Sizeof(sparseTable{})
var SizeofBitmap = unsafe.Sizeof(bitmap{})
var SizeofNodeI = unsafe.Sizeof([1]nodeI{})

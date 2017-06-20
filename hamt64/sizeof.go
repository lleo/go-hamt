package hamt64

import (
	"unsafe"
)

var SizeofFixedTable = unsafe.Sizeof(fixedTable{})
var SizeofSparseTable = unsafe.Sizeof(sparseTable{})
var SizeofBitmap = unsafe.Sizeof(Bitmap{})
var SizeofNodeI = unsafe.Sizeof([1]nodeI{})

package hamt32

import "strings"

type tableStack interface {
	peek() tableI
	pop() tableI
	push(tableI) tableStack
	//shift() tableI
	//unshift(tableI) tableStack
	isEmpty() bool
	len() int
}

//
// []tableI implementation of "tableStack interface".
//

type tableSlice []tableI

// Constructs an tableSlice impl of the tableStack interface.
func newTableStack() tableStack {
	var ts = make(tableSlice, 0, MaxDepth)
	return &ts
}

// path.peek() returns the last entry without inserted with path.push(...)
func (path *tableSlice) peek() tableI {
	if len(*path) == 0 {
		return nil
	}
	return (*path)[len(*path)-1]
}

// Put a new tableI in the path object.
// You should never push nil, but we are not checking to prevent this.
func (path *tableSlice) push(tab tableI) tableStack {
	//_ = ASSERT && Assert(tab != nil, "tableSlice.push(nil) not allowed")
	*path = append(*path, tab)
	return path
}

// path.pop() returns & remmoves the last entry inserted with path.push(...).
func (path *tableSlice) pop() tableI {
	if len(*path) == 0 {
		//FIXME: should I do this or let the runtime panic on index out of range
		return nil
	}

	var parent = (*path)[len(*path)-1]
	*path = (*path)[:len(*path)-1]
	return parent
}

//func (path *tableSlice) shift() tableI {
//	var t tableI
//	t, *path = (*path)[0], (*path)[1:]
//	return t
//}
//
//func (path *tableSlice) unshift(t tableI) tableStack {
//	*path = append([]tableI{t}, (*path)...)
//	return path
//}

// path.isEmpty() returns true if there are no entries in the path object,
// otherwise it returns false.
func (path *tableSlice) isEmpty() bool {
	return len(*path) == 0
}

func (path *tableSlice) len() int {
	return len(*path)
}

// Convert path to a string representation. This is only good for debug messages.
// It is not a string format to convert back from.
func (path tableSlice) String() string {
	var paths = make([]string, len(path))

	for i, pv := range path {
		paths[i] = pv.String()
	}

	var vals = strings.Join(paths, ", ")

	return "[ " + vals + " ]"
}

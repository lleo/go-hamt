package hamt64

import "strings"

type pathT []tableI

// Constructs an empty pathT object.
func newPathT() pathT {
	return pathT(make([]tableI, 0, MAXDEPTH+1))
}

// path.peek() returns the last entry without inserted with path.push(...)
// modifying path.
func (path pathT) peek() tableI {
	if len(path) == 0 {
		return nil
	}
	return path[len(path)-1]
}

// path.pop() returns & remmoves the last entry inserted with path.push(...).
func (path *pathT) pop() tableI {
	if len(*path) == 0 {
		//should I do this or let the runtime panic on index out of range
		return nil
	}
	parent := (*path)[len(*path)-1]
	*path = (*path)[:len(*path)-1]
	return parent

}

// Put a new tableI in the path object.
// You should never push nil, but we are not checking to prevent this.
func (path *pathT) push(node tableI) {
	//_ = ASSERT && Assert(node != nil, "pathT.push(nil) not allowed")
	*path = append(*path, node)
}

// path.isEmpty() returns true if there are no entries in the path object,
// otherwise it returns false.
func (path *pathT) isEmpty() bool {
	return len(*path) == 0
}

// Convert path to a string representation. This is only good for debug messages.
// It is not a string format to convert back from.
func (path *pathT) String() string {
	s := "["
	pvs := []tableI(*path)
	strs := make([]string, 0, 2)
	for _, pv := range pvs {
		strs = append(strs, pv.String())
	}
	s += strings.Join(strs, " ")
	s += "]"

	return s
}

package hamt32

import "strings"

type path32T []table32I

// Constructs an empty path32T object.
func newPath32T() path32T {
	return path32T(make([]table32I, 0, DEPTHLIMIT32))
}

// path.peek() returns the last entry without inserted with path.push(...)
// modifying path.
func (path path32T) peek() table32I {
	if len(path) == 0 {
		return nil
	}
	return path[len(path)-1]
}

// path.pop() returns & remmoves the last entry inserted with path.push(...).
func (path *path32T) pop() table32I {
	if len(*path) == 0 {
		//should I do this or let the runtime panic on index out of range
		return nil
	}
	parent := (*path)[len(*path)-1]
	*path = (*path)[:len(*path)-1]
	return parent

}

// Put a new table32I in the path object.
// You should never push nil, but we are not checking to prevent this.
func (path *path32T) push(node table32I) {
	//_ = ASSERT && Assert(node != nil, "path32T.push(nil) not allowed")
	*path = append(*path, node)
}

// path.isEmpty() returns true if there are no entries in the path object,
// otherwise it returns false.
func (path *path32T) isEmpty() bool {
	return len(*path) == 0
}

// Convert path to a string representation. This is only good for debug messages.
// It is not a string format to convert back from.
func (path *path32T) String() string {
	s := "["
	pvs := []table32I(*path)
	strs := make([]string, 0, 2)
	for _, pv := range pvs {
		strs = append(strs, pv.String())
	}
	s += strings.Join(strs, " ")
	s += "]"

	return s
}

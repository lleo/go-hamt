package hamt64

import "strings"

type path64T []table64I

// Constructs an empty path64T object.
func newPath64T() path64T {
	return path64T(make([]table64I, 0, DEPTHLIMIT64))
}

// path.peek() returns the last entry without inserted with path.push(...)
// modifying path.
func (path path64T) peek() table64I {
	if len(path) == 0 {
		return nil
	}
	return path[len(path)-1]
}

// path.pop() returns & remmoves the last entry inserted with path.push(...).
func (path *path64T) pop() table64I {
	if len(*path) == 0 {
		//should I do this or let the runtime panic on index out of range
		return nil
	}
	parent := (*path)[len(*path)-1]
	*path = (*path)[:len(*path)-1]
	return parent

}

// Put a new table64I in the path object.
// You should never push nil, but we are not checking to prevent this.
func (path *path64T) push(node table64I) {
	//_ = ASSERT && Assert(node != nil, "path64T.push(nil) not allowed")
	*path = append(*path, node)
}

// path.isEmpty() returns true if there are no entries in the path object,
// otherwise it returns false.
func (path *path64T) isEmpty() bool {
	return len(*path) == 0
}

// Convert path to a string representation. This is only good for debug messages.
// It is not a string format to convert back from.
func (path *path64T) String() string {
	s := "["
	pvs := []table64I(*path)
	strs := make([]string, 0, 2)
	for _, pv := range pvs {
		strs = append(strs, pv.String())
	}
	s += strings.Join(strs, " ")
	s += "]"

	return s
}

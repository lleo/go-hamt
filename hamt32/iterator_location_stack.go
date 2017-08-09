package hamt32

import (
	"fmt"
	"strings"
)

type iterLocation struct {
	tab tableI
	idx uint
}

type iterLocStack []iterLocation

func newIterLocStack() iterLocStack {
	var stack iterLocStack = make([]iterLocation, 0, DepthLimit)
	return stack
}

func (is *iterLocStack) len() int {
	return len(*is)
}

func (is *iterLocStack) push(tab tableI, idx uint) *iterLocStack {
	_ = assertOn && assert(tab != nil, "tab in tableStat is nil")
	*is = append(*is, iterLocation{tab, idx})
	return is
}

func (is *iterLocStack) pop() (tableI, uint) {
	var last = is.len() - 1
	var tab = (*is)[last].tab
	var idx = (*is)[last].idx

	*is = (*is)[:last]

	return tab, idx
}

func (is *iterLocStack) String() string {
	var strs = make([]string, is.len())

	for i := 0; i < is.len(); i++ {
		var ise = (*is)[i]
		strs[i] = fmt.Sprintf("%d: tableI=%s; idx=%d;", i, ise.tab, ise.idx)
	}

	return strings.Join(strs, "\n") + "\n"
}

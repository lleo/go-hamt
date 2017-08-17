package hamt32

type tableIterStack []tableIterFunc

func newTableIterStack() tableIterStack {
	var ts tableIterStack = make([]tableIterFunc, 0, DepthLimit)
	return ts
}

func (ts *tableIterStack) push(f tableIterFunc) {
	(*ts) = append(*ts, f)
}

func (ts *tableIterStack) pop() tableIterFunc {
	if len(*ts) == 0 {
		return nil
	}

	var last = len(*ts) - 1
	var f = (*ts)[last]
	*ts = (*ts)[:last]

	return f
}

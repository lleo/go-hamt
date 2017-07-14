package hamt64

import (
	"fmt"
)

// iKeyVal is a simple struct used to transfer lists ([]iKeyVal) from one
// function to another.
type iKeyVal struct {
	Key *iKey
	Val interface{}
}

func (kv iKeyVal) String() string {
	return fmt.Sprintf("{%s, %v}", kv.Key, kv.Val)
}

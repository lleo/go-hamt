package key

import (
	"fmt"
)

// KeyVal is a simple struct used to transfer lists ([]KeyVal) from one
// function to another.
type KeyVal struct {
	Key Key
	Val interface{}
}

func (kv KeyVal) String() string {
	return fmt.Sprintf("{%s, %v}", kv.Key, kv.Val)
}

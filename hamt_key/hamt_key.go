/*
The hamt_key package contains a single Key interface. The hamt_key package was
created to prevent cicular depedencies betwee "github.com/lleo/go-hamt" and
either "github.com/lleo/go-hamt/hamt32" or "github.com/lleo/go-hamt/hamt64".

However the hamt_key pacakge is also used by the functional Hamt variation in
"github.com/lleo/go-hamt-functional".
*/
package hamt_key

type Key interface {
	Equals(Key) bool
	Hash30() uint32
	Hash60() uint64
	String() string
}

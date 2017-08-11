package hamt64

import "fmt"

// The assertOn constant determines whether or not assert() calls are called.
// When this constantant is false, statements of the form:
//     _ = assertOn && assert(...)
// become noops when compiled.
// NOTE: This constant SHOULD BE false for production code.
const assertOn bool = false

// assert() tests if test is false; if it is, it will panic with msg.
// assert() is the fastest as it is simple enough to be inlined.
func assert(test bool, msg string) bool {
	if !test {
		panic(msg)
	}
	return true
}

// assertf() tests if test is false; if it is, it will panic with a message
// formatted by msgFmt and msgArgs via fmt.Sprintf().
// assertf() is much slower as it is not inlined. I am guessing  this is due to
// the vararg nature of the arguments. However, nooping via the assertOn trick
// still applies.
func assertf(test bool, msgFmt string, msgArgs ...interface{}) bool {
	if !test {
		var msg = fmt.Sprintf(msgFmt, msgArgs...)
		panic(msg)
	}
	return true
}

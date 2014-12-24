//po:MsgId "argument of Offsetof implies indirection of an embedded field"
//po:MsgStr ""

package p

import "unsafe"

type T1 struct {
	X int
}

type T struct {
	*T1
}

var x = unsafe.Offsetof(T{}.X)

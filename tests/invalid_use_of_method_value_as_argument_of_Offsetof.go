//po:MsgId "invalid use of method value as argument of Offsetof"
//po:MsgStr "使用模板名%qE时不带实参表无效"

package p

import "unsafe"

type T struct{}

func (T) F() {}

var x = unsafe.Offsetof(T{}.F)

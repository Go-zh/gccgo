//po:MsgId "invalid use of method value as argument of Offsetof"
//po:MsgStr "%<unsafe.Offsetof%>不接受方法值作为参数"

package p

import "unsafe"

type T struct{}

func (T) F() {}

var x = unsafe.Offsetof(T{}.F)

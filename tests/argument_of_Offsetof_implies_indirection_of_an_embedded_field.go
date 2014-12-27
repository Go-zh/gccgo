//po:MsgId "argument of Offsetof implies indirection of an embedded field"
//po:MsgStr "%<unsafe.Offsetof%>的参数字段并不直接包含于结构体内：至少要经过一级以指针形式内嵌的结构体"

package p

import "unsafe"

type T1 struct {
	X int
}

type T struct {
	*T1
}

var x = unsafe.Offsetof(T{}.X)

//po:MsgId "argument must be a field reference"
//po:MsgStr "实参必须是对结构体字段的引用"

package p

import "unsafe"

type T struct{}

var x = unsafe.Offsetof(T{})

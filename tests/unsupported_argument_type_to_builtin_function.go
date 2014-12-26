//po:MsgId "unsupported argument type to builtin function"
//po:MsgStr "内建函数不支持这种实参类型"

package p

func F() {
	println(F())
	print(F())
}

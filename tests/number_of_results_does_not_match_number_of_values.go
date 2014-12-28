//po:MsgId "number of results does not match number of values"
//po:MsgStr "返回值的数量和赋值的数量不一样多"

package p

func F() (int, int)

func G() {
	_, _, _ = F()
}

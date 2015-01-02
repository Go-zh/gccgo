//po:MsgId "invalid left hand side of assignment"
//po:MsgStr "被赋值的值无效"

package p

func F() {
	[10]int{1, 2}[1]++
}

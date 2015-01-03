//po:MsgId "cannot fallthrough final case in switch"
//po:MsgStr "switch 语句的最后一个分支不能再 fallthrough"

package p

func F(x int) {
	switch x {
	default:
		fallthrough
	}
}

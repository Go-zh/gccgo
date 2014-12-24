//po:MsgId "invalid use of %<...%> with builtin function"
//po:MsgStr "内建函数实参无效"

package p

func F() {
	chans := []chan bool{}
	close(chans...)
}

//po:MsgId "invalid use of %<...%> with builtin function"
//po:MsgStr "该内建函数不支持%<...%>"

package p

func F() {
	chans := []chan bool{}
	close(chans...)
}

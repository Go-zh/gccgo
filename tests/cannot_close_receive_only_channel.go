//po:MsgId "cannot close receive-only channel"
//po:MsgStr "无法关闭只读信道"

package p

var x <-chan struct{}

func F() {
	close(x)
}

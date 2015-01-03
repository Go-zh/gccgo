//po:MsgId "invalid receive on send-only channel"
//po:MsgStr "无法从单向只发送信道中接收"

package p

var x chan<- int

var y = <-x

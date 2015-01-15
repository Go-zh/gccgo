//po:MsgId "method %s%s%s requires a pointer receiver"
//po:MsgStr "方法 %s%s%s 需要一个指针接收者"

package p

type T struct{}

func (*T) F() {}

type I interface {
	F()
}

var _ I = T{}

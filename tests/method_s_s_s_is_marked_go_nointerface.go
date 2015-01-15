//po:MsgId "method %s%s%s is marked go:nointerface"
//po:MsgStr "方法 %s%s%s 被标记为 go:nointerface"

package p

type I interface {
	F()
}

type x struct{}

//go:nointerface
func (x) F() {}

var _ I = x{}

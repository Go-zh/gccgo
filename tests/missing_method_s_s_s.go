//po:MsgId "missing method %s%s%s"
//po:MsgStr "缺少方法%s%s%s"

package p

type I interface {
	F()
	G()
}

type T struct{}

func (T) G() {}

var _ I = T{}

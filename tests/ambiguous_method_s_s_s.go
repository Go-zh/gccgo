//po:MsgId "ambiguous method %s%s%s"
//po:MsgStr "方法 %s%s%s 有歧义"

package p

type T1 struct{}

func (T1) F() {}

type T2 struct{}

func (T2) F() {}

type T struct {
	T1
	T2
}

type I interface {
	F()
}

var _ I = T{}

//po:MsgId "impossible type assertion: type does not implement interface"
//po:MsgStr "不可能的类型断言：所断言的类型并不实现该接口"

package p

type I interface {
	F()
}

type T1 struct{}

func (*T1) F() {}

var x I

var y = x.(T1)

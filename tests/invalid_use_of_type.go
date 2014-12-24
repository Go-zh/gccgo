//po:MsgId "invalid use of type"
//po:MsgStr "类型使用无效"

package p

type T int

func F() func () T {
	return func() T
}

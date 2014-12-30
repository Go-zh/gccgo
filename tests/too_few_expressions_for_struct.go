//po:MsgId "too few expressions for struct"
//po:MsgStr "结构字面的初始化字段值过少"

package p

type T struct {
	A, B, C int
}

var x = T{1}

//po:MsgId "too many expressions for struct"
//po:MsgStr "结构字面的初始化字段值过多"

package p

type T struct {
	A, B, C int
}

var x = T{1, 2, 3, 4}

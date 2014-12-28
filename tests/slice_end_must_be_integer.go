//po:MsgId "slice end must be integer"
//po:MsgStr "切片的端点索引必须是整数"

package p

var x []byte

var (
	//_ = x[2.718:] // index must be integer
	_ = x[:3.14] // slice end must be integer
)

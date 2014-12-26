//po:MsgId "too many arguments to make"
//po:MsgStr "make 的实参过多"

package p

var (
	z = make(map[int]int, 1, 2)
	x = make(chan struct{}, 1, 2)
	y = make([]byte, 1, 2, 3)
)

//po:MsgId "left argument must be a slice"
//po:MsgStr "%<copy%> 的第一个实参必须是一个切片"

package p

var x = copy(nil, []byte{1})

//po:MsgId "invalid use of %<...%> with non-slice"
//po:MsgStr "在文件作用域使用%<this%>无效"

package p

func F(int, []int) int

var x [1]int
var y = F(1, x...)

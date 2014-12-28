//po:MsgId "invalid use of %<...%> calling non-variadic function"
//po:MsgStr "使用<%...%>调用非变参函数"

package p

func F([]int) int

var x []int

var y = F(x...)

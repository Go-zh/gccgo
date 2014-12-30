//po:MsgId "type assertion only valid for interface types"
//po:MsgStr "类型断言仅可用于接口类型"

package p

type T int

var x T

var y = x.(int)

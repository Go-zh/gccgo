//po:MsgId "incompatible type for map index"
//po:MsgStr "字典索引的类型不兼容"

package p

var x map[int]int

type T int

var y = x[T(1)]

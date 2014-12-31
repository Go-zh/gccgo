//po:MsgId "slice of unaddressable value"
//po:MsgStr "无法对不可取址的值进行切片"

package p

var x = [10]byte{}[:]

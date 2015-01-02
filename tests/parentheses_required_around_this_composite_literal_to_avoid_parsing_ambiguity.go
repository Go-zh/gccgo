//po:MsgId "parentheses required around this composite literal to avoid parsing ambiguity"
//po:MsgStr "为避免语法解析时产生歧义，此复合字面必须用圆括号括起"

package p

type T struct{}

func F() {
	if &T{} == nil {
		return
	}
}

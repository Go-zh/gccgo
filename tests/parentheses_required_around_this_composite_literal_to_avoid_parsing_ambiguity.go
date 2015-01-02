//po:MsgId "parentheses required around this composite literal to avoid parsing ambiguity"
//po:MsgStr "为避免有歧义的语法分析，必须在用圆括号将这个复合字面括起来"

package p

type T struct{}

func F() {
	if &T{} == nil {
		return
	}
}

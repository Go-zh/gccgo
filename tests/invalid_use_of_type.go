package p

type T int

func F() func () T {
	return func() T
}

package math

type math struct {
	a int
	b int
	c int
}

func NewMath(a, b, c int) math {
	return math{a: a, b: b, c: c}
}

func (m math) Add() int {
	return m.a + m.b + m.c
}

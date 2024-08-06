package ruby

import "golang.org/x/exp/constraints"

func R[T constraints.Integer](start, end, step T) Enumerable[T] {
	return &enumerableImpl[T]{EnumeratorGenerator: &rangeEnumeratorGenerator[T]{start, end, step}}
}

type rangeEnumerator[T constraints.Integer] struct {
	start, end, step T
	pos              T
}

func (e *rangeEnumerator[T]) hasNext() bool {
	return e.pos < e.end
}

func (e *rangeEnumerator[T]) next() T {
	res := e.pos
	e.pos += e.step
	return res
}

type rangeEnumeratorGenerator[T constraints.Integer] struct {
	start, end, step T
}

func (g *rangeEnumeratorGenerator[T]) create() Enumerator[T] {
	return &rangeEnumerator[T]{g.start, g.end, g.step, g.start}
}

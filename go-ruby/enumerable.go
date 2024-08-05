package go_ruby

import (
	"golang.org/x/exp/constraints"
)

type Enumerator[T any] interface {
	hasNext() bool
	next() T
}

type EnumeratorGenerator[T any] interface {
	Create() Enumerator[T]
}

type Enumerable[T any] interface {
	Each(func(T))
	EachWithIndex(func(int, T))
	Includes(T, func(T, T) bool) bool
}

type SliceEnumeratorGenerator[T any] struct {
	data []T
}

type SliceEnumerator[T any] struct {
	data *[]T
	pos  int
}

func (g *SliceEnumeratorGenerator[T]) Create() Enumerator[T] {
	return &SliceEnumerator[T]{&g.data, 0}
}

func (g *SliceEnumerator[T]) hasNext() bool {
	return g.pos < len(*g.data)
}

func (g *SliceEnumerator[T]) next() T {
	return (*g.data)[g.pos]
}

func NewSliceEnumerable[T any](slice []T) Enumerable[T] {
	return &EnumerableImpl[T]{EnumeratorGenerator: &SliceEnumeratorGenerator[T]{slice}}
}

type EnumerableImpl[T any] struct {
	EnumeratorGenerator[T]
}

func (e *EnumerableImpl[T]) Each(f func(T)) {
	enumerator := e.EnumeratorGenerator.Create()
	for enumerator.hasNext() {
		f(enumerator.next())
	}
}

func (e *EnumerableImpl[T]) EachWithIndex(f func(int, T)) {
	enumerator := e.EnumeratorGenerator.Create()
	i := 0
	for enumerator.hasNext() {
		f(i, enumerator.next())
		i++
	}
}

func (e *EnumerableImpl[T]) Includes(t T, lessOrEqual func(T, T) bool) bool {
	enumerator := e.EnumeratorGenerator.Create()
	for enumerator.hasNext() {
		el := enumerator.next()
		if lessOrEqual(t, el) && lessOrEqual(el, t) {
			return true
		}
	}
	return false
}

type RangeEnumerator[T constraints.Integer] struct {
	start, end, step T
	pos              T
}

func (e *RangeEnumerator[T]) hasNext() bool {
	return e.pos < e.end
}

func (e *RangeEnumerator[T]) next() T {
	res := e.pos
	e.pos += e.step
	return res
}

type RangeEnumeratorGenerator[T constraints.Integer] struct {
	start, end, step T
}

func (g *RangeEnumeratorGenerator[T]) Create() Enumerator[T] {
	return &RangeEnumerator[T]{g.start, g.end, g.step, g.start}
}

func NewRange[T constraints.Integer](start, end, step T) Enumerable[T] {
	return &EnumerableImpl[T]{EnumeratorGenerator: &RangeEnumeratorGenerator[T]{start, end, step}}
}

package ruby

func E[T any](slice []T) Enumerable[T] {
	return &enumerableImpl[T]{EnumeratorGenerator: &sliceEnumeratorGenerator[T]{slice}}
}

type sliceEnumeratorGenerator[T any] struct {
	data []T
}

type sliceEnumerator[T any] struct {
	data *[]T
	pos  int
}

func (g *sliceEnumeratorGenerator[T]) create() Enumerator[T] {
	return &sliceEnumerator[T]{&g.data, 0}
}

func (g *sliceEnumerator[T]) hasNext() bool {
	return g.pos < len(*g.data)
}

func (g *sliceEnumerator[T]) next() T {
	return (*g.data)[g.pos]
}

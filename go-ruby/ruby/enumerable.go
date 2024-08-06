package ruby

type enumerableImpl[T any] struct {
	EnumeratorGenerator[T]
}

func (e *enumerableImpl[T]) Each(f func(T)) {
	enumerator := e.EnumeratorGenerator.create()
	for enumerator.hasNext() {
		f(enumerator.next())
	}
}

func (e *enumerableImpl[T]) EachWithIndex(f func(int, T)) {
	enumerator := e.EnumeratorGenerator.create()
	i := 0
	for enumerator.hasNext() {
		f(i, enumerator.next())
		i++
	}
}

func (e *enumerableImpl[T]) Includes(t T, lessOrEqual func(T, T) bool) bool {
	enumerator := e.EnumeratorGenerator.create()
	for enumerator.hasNext() {
		el := enumerator.next()
		if lessOrEqual(t, el) && lessOrEqual(el, t) {
			return true
		}
	}
	return false
}

func (e *enumerableImpl[T]) Entries() []T {
	a := make([]T, 0)
	e.Each(func(el T) {
		a = append(a, el)
	})
	return a
}

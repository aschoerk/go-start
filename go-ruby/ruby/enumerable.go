package ruby

import "reflect"

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

func isNil(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return v.IsNil()
	}
	return false
}

func (e *enumerableImpl[T]) All(f ...Predicate[T]) bool {
	if len(f) > 1 {
		panic("Invalid usage of All")
	}
	if len(f) == 0 {
		return e.All(func(x T) bool {
			return !isNil(reflect.ValueOf(x))
		})
	} else {
		enumerator := e.EnumeratorGenerator.create()
		for enumerator.hasNext() {
			if !f[0](enumerator.next()) {
				return false
			}
		}
		return true
	}
}

func (e *enumerableImpl[T]) Any(f func(T) bool) bool {
	enumerator := e.EnumeratorGenerator.create()
	for enumerator.hasNext() {
		if !f(enumerator.next()) {
			return true
		}
	}
	return false
}

func (e *enumerableImpl[T]) None(f func(T) bool) bool {
	return !e.Any(f)
}

func (e *enumerableImpl[T]) One(f func(T) bool) bool {
	enumerator := e.EnumeratorGenerator.create()
	found := false
	for enumerator.hasNext() {
		if f(enumerator.next()) {
			if found {
				return false
			} else {
				found = true
			}
		}
	}
	return found
}

func (e *enumerableImpl[T]) Count(f ...Predicate[T]) int {
	if len(f) > 1 {
		panic("Invalid usage of Count")
	}
	if len(f) == 0 {
		return e.Count(func(x T) bool {
			return true
		})
	} else {
		res := 0
		enumerator := e.EnumeratorGenerator.create()
		for enumerator.hasNext() {
			if f[0](enumerator.next()) {
				res++
			}
		}
		return res
	}
}

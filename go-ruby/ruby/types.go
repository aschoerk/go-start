package ruby

type Enumerator[T any] interface {
	hasNext() bool
	next() T
}

type EnumeratorGenerator[T any] interface {
	create() Enumerator[T]
}

type Enumerable[T any] interface {
	Each(func(T))
	EachWithIndex(func(int, T))
	Includes(T, func(T, T) bool) bool
	Entries() []T
}

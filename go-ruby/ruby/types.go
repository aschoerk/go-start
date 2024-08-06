package ruby

type Enumerator[T any] interface {
	hasNext() bool
	next() T
}

type EnumeratorGenerator[T any] interface {
	create() Enumerator[T]
}

type Predicate[T any] func(T) bool

type Enumerable[T any] interface {
	// Querying
	Includes(T, func(T, T) bool) bool
	// All() bool
	All(...Predicate[T]) bool
	Any(func(T) bool) bool
	None(func(T) bool) bool
	One(func(T) bool) bool
	Count(...Predicate[T]) int
	// Tally()

	// Iterating
	Each(func(T))
	EachWithIndex(func(int, T))
	Entries() []T
}

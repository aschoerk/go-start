package go_ruby

import (
	"testing"
)

func TestEach(t *testing.T) {

	var e Enumerable[int] = NewRange[int](1, 11, 1)
	var count int = 0
	e.Each(func(el int) {
		count++
		if el != count {
			t.Errorf("Expected %d, but got %d", count+1, el)
		}
		if !e.Includes(el, func(a int, b int) bool {
			return a <= b
		}) {
			t.Errorf("Expected Element %d, to be found", el)
		}
		if e.Includes(el+10, func(a int, b int) bool {
			return a <= b
		}) {
			t.Errorf("Expected Element %d, not to be found", el+10)
		}
	})

}

func TestEachWithIndex(t *testing.T) {

	e := NewRange[int](1, 11, 1)
	var count int = 0
	e.EachWithIndex(func(index int, el int) {
		count++
		if el != count {
			t.Errorf("Expected %d, but got %d", count, el)
		}
		if index != count-1 {
			t.Errorf("Expected index %d, but got %d", count-1, index)
		}
	})

}

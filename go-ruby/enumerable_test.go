package main

import (
	"testing"

	"aschoerk.de/go-ruby/ruby"
)

func TestEach(t *testing.T) {

	e := ruby.RStepped[int](1, 11, 1)
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

	e := ruby.RStepped[int](1, 11, 1)
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

func TestSlice(t *testing.T) {
	e := ruby.E(ruby.R[int](1, 3).Entries())
	if e.Count() != 2 {
		t.Errorf("Expected only 2 entries, but was %d", e.Count())
	}
	if e.Count(func(a int) bool { return a > 0 }) != 2 {
		t.Errorf("Expected only 2 entries, but was %d", e.Count())
	}
}

package main

import (
	"log"
	"sync"
)

type BooleanBuffers interface {
	current() BooleanBuffer
	next() bool
	prev() bool
	nextGeneration() (BooleanBuffer, bool)
	mu() *sync.Mutex
	changeSizes(uint, uint)
}

type BooleanBuffersImpl struct {
	size    uint
	ptr     uint
	buffers []BooleanBuffer
	mutex   sync.Mutex
}

func initBuffers(size uint, maxX uint, maxY uint) BooleanBuffers {
	var res BooleanBuffersImpl
	res.size = size
	res.ptr = 0
	res.buffers = make([]BooleanBuffer, size)
	res.buffers[0] = initBuffer(maxX, maxY)
	return &res
}

func (b *BooleanBuffersImpl) current() BooleanBuffer {
	var res BooleanBuffer = b.buffers[b.ptr]
	if res == nil {
		log.Fatal("current must always be not nil")
	}
	return res
}

func (b *BooleanBuffersImpl) relative(diff int) bool {
	example := b.current()
	var prevPtr = (int(b.ptr) + diff) % int(b.size)
	if prevPtr < 0 {
		prevPtr += int(b.size)
	}
	res := b.buffers[prevPtr]
	if res != nil {
		b.ptr = uint(prevPtr)
		res.changeSizeNotDestructing(example.maxX(), example.maxY())
		return true
	} else {
		return false
	}
}

func (b *BooleanBuffersImpl) progress() BooleanBuffer {
	example := b.current()
	var newPtr = (b.ptr + 1) % b.size

	res := b.buffers[newPtr]
	b.ptr = uint(newPtr)
	if res != nil {
		res.changeSizeNotDestructing(example.maxX(), example.maxY())
		return res
	} else {
		b.buffers[b.ptr] = initBuffer(example.maxX(), example.maxY())
		return b.buffers[b.ptr]
	}
}

func (b *BooleanBuffersImpl) next() bool {
	return b.relative(1)
}

func (b *BooleanBuffersImpl) prev() bool {
	return b.relative(-1)
}

func (b *BooleanBuffersImpl) nextGeneration() (BooleanBuffer, bool) {
	current := b.current()
	var prev BooleanBuffer = nil
	if b.prev() {
		prev = b.current()
		b.next()
	}
	var res BooleanBuffer = nil
	var changed bool
	newCurrent := b.progress()
	res, changed = current.nextGeneration(newCurrent)
	if prev != nil && res.equals(prev) {
		b.prev()
		b.prev()
	}

	if !changed {
		if !b.prev() {
			log.Fatal("prev after progress must always be not nil")
		}
	}
	return res, changed
}

func (b *BooleanBuffersImpl) mu() *sync.Mutex {
	return &b.mutex
}

func (b *BooleanBuffersImpl) changeSizes(maxX uint, maxY uint) {
	newB := b.current().changeSizeNotDestructing(maxX, maxY)
	b.buffers[b.ptr] = newB
}

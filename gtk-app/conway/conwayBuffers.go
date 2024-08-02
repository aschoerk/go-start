package conway

import (
	"log"
	"sync"
)

type ConwayBuffers interface {
	Current() ConwayBuffer
	Next() bool
	Prev() bool
	NextGeneration() (ConwayBuffer, bool)
	Mu() *sync.Mutex
	ChangeSizeNotDestructing(maxX uint32, maxY uint32)
	CurrentBounds() (uint32, uint32)
	changeSizes(uint32, uint32)
}

func InitBuffers(size uint32, maxX uint32, maxY uint32) ConwayBuffers {
	var res conwayBuffersImpl
	res.size = size
	res.ptr = 0
	res.buffers = make([]ConwayBuffer, size)
	res.buffers[0] = newBuffer(maxX, maxY)
	return &res
}

type conwayBuffersImpl struct {
	size    uint32
	ptr     uint32
	buffers []ConwayBuffer
	mutex   sync.Mutex
}

func (b *conwayBuffersImpl) Current() ConwayBuffer {
	var res ConwayBuffer = b.buffers[b.ptr]
	if res == nil {
		log.Fatal("current must always be not nil")
	}
	return res
}

func (b *conwayBuffersImpl) CurrentBounds() (uint32, uint32) {
	return currentMaxX, currentMaxY
}

func (b *conwayBuffersImpl) relative(diff int) bool {
	example := b.Current()
	var prevPtr = (int(b.ptr) + diff) % int(b.size)
	if prevPtr < 0 {
		prevPtr += int(b.size)
	}
	res := b.buffers[prevPtr]
	if res != nil {
		b.ptr = uint32(prevPtr)
		res.changeSizeNotDestructing(example.MaxX(), example.MaxY())
		return true
	} else {
		return false
	}
}

func (b *conwayBuffersImpl) progress() ConwayBuffer {
	example := b.Current()
	var newPtr = (b.ptr + 1) % b.size

	res := b.buffers[newPtr]
	b.ptr = uint32(newPtr)
	if res != nil {
		res.changeSizeNotDestructing(example.MaxX(), example.MaxY())
		return res
	} else {
		b.buffers[b.ptr] = newBuffer(example.MaxX(), example.MaxY())
		return b.buffers[b.ptr]
	}
}

func (b *conwayBuffersImpl) Next() bool {
	return b.relative(1)
}

func (b *conwayBuffersImpl) Prev() bool {
	return b.relative(-1)
}

func (b *conwayBuffersImpl) NextGeneration() (ConwayBuffer, bool) {
	current := b.Current()
	var prev ConwayBuffer = nil
	if b.Prev() {
		prev = b.Current()
		b.Next()
	}
	var res ConwayBuffer = nil
	var changed bool
	newCurrent := b.progress()
	res, changed = current.nextGeneration(newCurrent)
	if prev != nil && res.equals(prev) {
		b.Prev()
		b.Prev()
	}

	if !changed {
		if !b.Prev() {
			log.Fatal("prev after progress must always be not nil")
		}
	}
	return res, changed
}

func (b *conwayBuffersImpl) Mu() *sync.Mutex {
	return &b.mutex
}

func (b *conwayBuffersImpl) changeSizes(maxX uint32, maxY uint32) {
	newB := b.Current().changeSizeNotDestructing(maxX, maxY)
	b.buffers[b.ptr] = newB
}

func (b *conwayBuffersImpl) ChangeSizeNotDestructing(maxX uint32, maxY uint32) {
	b.Mu().Lock()
	b.buffers[b.ptr] = b.Current().changeSizeNotDestructing(maxX, maxY)
	b.Mu().Unlock()
}

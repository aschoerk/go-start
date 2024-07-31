package main

import (
	"log"
	"sync"
)

type BooleanBuffer interface {
	get(x uint, y uint) bool
	set(x uint, y uint, val bool)
	newEmptyBuffer() BooleanBuffer
	nextGeneration(BooleanBuffer) (BooleanBuffer, bool)
	changeSizeDestructing(maxX uint, maxY uint)
	changeSizeNotDestructing(maxX uint, maxY uint) BooleanBuffer
	maxX() uint
	maxY() uint
	equals(BooleanBuffer) bool
}

type BooleanBufferImpl struct {
	maxXVal uint
	maxYVal uint
	vals    []bool
}

func (b *BooleanBufferImpl) maxX() uint {
	return b.maxXVal
}

func (b *BooleanBufferImpl) maxY() uint {
	return b.maxYVal
}

func (b *BooleanBufferImpl) equals(other BooleanBuffer) bool {
	sizesOk := b.maxXVal == other.maxX() && b.maxYVal == other.maxY()
	if sizesOk {
		for i := uint(0); i < b.maxXVal; i++ {
			for j := uint(0); j < b.maxYVal; j++ {
				if b.get(i, j) != other.get(i, j) {
					return false
				}
			}
		}
		return true
	} else {
		return false
	}
}

func initBuffer(maxX uint, maxY uint) BooleanBuffer {
	var res = BooleanBufferImpl{
		maxXVal: maxX,
		maxYVal: maxY,
		vals:    make([]bool, maxX*maxY)}

	return &res
}

func (bb *BooleanBufferImpl) changeSizeDestructing(maxX uint, maxY uint) {
	bb.vals = make([]bool, maxX*maxY)
}

func (bb *BooleanBufferImpl) get(x uint, y uint) bool {
	if x < bb.maxXVal && y < bb.maxYVal {
		return bb.vals[x+y*bb.maxXVal]
	}
	log.Fatal("Invalid range")
	return false
}

func (bb *BooleanBufferImpl) set(x uint, y uint, val bool) {
	if x < bb.maxXVal && y < bb.maxYVal {
		bb.vals[x+y*bb.maxXVal] = val
	} else {
		log.Fatal("Invalid range")
	}
}

func (b *BooleanBufferImpl) newEmptyBuffer() BooleanBuffer {
	var res BooleanBufferImpl
	res.maxXVal = b.maxXVal
	res.maxYVal = b.maxYVal
	res.vals = make([]bool, b.maxXVal*b.maxYVal)
	return &res
}

func (b *BooleanBufferImpl) changeSizeNotDestructing(maxX uint, maxY uint) BooleanBuffer {
	if maxX != b.maxX() || maxY != b.maxY() {
		newB := initBuffer(maxX, maxY)
		for x := uint(0); x < min(maxX, b.maxX()); x++ {
			for y := uint(0); y < min(maxY, b.maxY()); y++ {
				newB.set(x, y, b.get(x, y))
			}
		}
		return newB
	} else {
		return b
	}
}

func (b *BooleanBufferImpl) countNeighbors(x, y uint) uint {
	var count uint
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := int(x)+dx, int(y)+dy
			if nx >= 0 && nx < int(b.maxXVal) && ny >= 0 && ny < int(b.maxYVal) && b.get(uint(nx), uint(ny)) {
				count++
			}
		}
	}
	return count
}

func (b *BooleanBufferImpl) nextGeneration(buffer BooleanBuffer) (BooleanBuffer, bool) {

	if buffer != nil {
		if buffer.maxX() != b.maxXVal || buffer.maxY() != b.maxYVal {
			buffer.changeSizeDestructing(b.maxXVal, b.maxYVal)
		}
	} else {
		buffer = initBuffer(b.maxXVal, b.maxYVal)
	}

	isChanged := false

	for y := uint(0); y < b.maxYVal; y++ {
		for x := uint(0); x < b.maxXVal; x++ {
			neighbors := b.countNeighbors(x, y)
			var tmp bool
			if b.get(x, y) {
				tmp = neighbors == 2 || neighbors == 3
			} else {
				tmp = neighbors == 3
			}
			if tmp != b.get(x, y) {
				isChanged = true
			}
			buffer.set(x, y, tmp)
		}
	}
	return buffer, isChanged
}

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

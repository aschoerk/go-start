package conway

import (
	"log"
	"sync/atomic"
)

var (
	currentMaxX uint32 = 0
	currentMaxY uint32 = 0
)

type ConwayBuffer interface {
	Get(x uint32, y uint32) bool
	Set(x uint32, y uint32, val bool)
	MaxX() uint32
	MaxY() uint32
	newEmptyBuffer() ConwayBuffer
	nextGeneration(ConwayBuffer) (ConwayBuffer, bool)
	changeSizeDestructing(maxX uint32, maxY uint32)
	changeSizeNotDestructing(maxX uint32, maxY uint32) ConwayBuffer
	equals(ConwayBuffer) bool
	handleAtomics()
}

type conwayBufferImpl struct {
	maxXVal uint32
	maxYVal uint32
	vals    []bool
}

func (b *conwayBufferImpl) MaxX() uint32 {
	return b.maxXVal
}

func (b *conwayBufferImpl) MaxY() uint32 {
	return b.maxYVal
}

func (b *conwayBufferImpl) equals(other ConwayBuffer) bool {
	sizesOk := b.maxXVal == other.MaxX() && b.maxYVal == other.MaxY()
	if sizesOk {
		for i := uint32(0); i < b.maxXVal; i++ {
			for j := uint32(0); j < b.maxYVal; j++ {
				if b.Get(i, j) != other.Get(i, j) {
					return false
				}
			}
		}
		return true
	} else {
		return false
	}
}

func newBuffer(maxX uint32, maxY uint32) ConwayBuffer {
	var res = conwayBufferImpl{
		maxXVal: maxX,
		maxYVal: maxY,
		vals:    make([]bool, maxX*maxY)}

	res.handleAtomics()
	return &res
}

func (b *conwayBufferImpl) handleAtomics() {
	atomic.StoreUint32(&currentMaxX, max(currentMaxX, b.MaxX()))
	atomic.StoreUint32(&currentMaxY, max(currentMaxY, b.MaxY()))
}

func (b *conwayBufferImpl) changeSizeDestructing(maxX uint32, maxY uint32) {
	b.maxXVal = maxX
	b.maxYVal = maxY
	b.vals = make([]bool, maxX*maxY)
	b.handleAtomics()
}

func (bb *conwayBufferImpl) Get(x uint32, y uint32) bool {
	if x < bb.maxXVal && y < bb.maxYVal {
		return bb.vals[x+y*bb.maxXVal]
	}
	log.Fatal("Invalid range")
	return false
}

func (bb *conwayBufferImpl) Set(x uint32, y uint32, val bool) {
	if x < bb.maxXVal && y < bb.maxYVal {
		bb.vals[x+y*bb.maxXVal] = val
	} else {
		log.Fatal("Invalid range")
	}
}

func (b *conwayBufferImpl) newEmptyBuffer() ConwayBuffer {
	return newBuffer(b.maxXVal, b.maxYVal)
}

func (b *conwayBufferImpl) changeSizeNotDestructing(maxX uint32, maxY uint32) ConwayBuffer {
	if maxX > b.MaxX() || maxY > b.MaxY() {
		newB := newBuffer(maxX, maxY)
		for x := uint32(0); x < min(maxX, b.MaxX()); x++ {
			for y := uint32(0); y < min(maxY, b.MaxY()); y++ {
				newB.Set(x, y, b.Get(x, y))
			}
		}
		return newB
	} else {
		return b
	}
}

func (b *conwayBufferImpl) countNeighbors(x, y uint32) uint32 {
	var count uint32
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := int(x)+dx, int(y)+dy
			if nx >= 0 && nx < int(b.maxXVal) && ny >= 0 && ny < int(b.maxYVal) && b.Get(uint32(nx), uint32(ny)) {
				count++
			}
		}
	}
	return count
}

func (b *conwayBufferImpl) nextGeneration(buffer ConwayBuffer) (ConwayBuffer, bool) {

	if buffer != nil {
		if buffer.MaxX() != b.maxXVal || buffer.MaxY() != b.maxYVal {
			buffer.changeSizeDestructing(b.maxXVal, b.maxYVal)
		}
	} else {
		buffer = newBuffer(b.maxXVal, b.maxYVal)
	}

	isChanged := false

	for y := uint32(0); y < b.maxYVal; y++ {
		for x := uint32(0); x < b.maxXVal; x++ {
			neighbors := b.countNeighbors(x, y)
			var tmp bool
			if b.Get(x, y) {
				tmp = neighbors == 2 || neighbors == 3
			} else {
				tmp = neighbors == 3
			}
			if tmp != b.Get(x, y) {
				isChanged = true
			}
			buffer.Set(x, y, tmp)
		}
	}
	return buffer, isChanged
}

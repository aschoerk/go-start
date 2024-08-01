package main

import "log"

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

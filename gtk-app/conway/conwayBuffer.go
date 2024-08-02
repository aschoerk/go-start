package conway

import "log"

type ConwayBuffer interface {
	Get(x uint, y uint) bool
	Set(x uint, y uint, val bool)
	ChangeSizeNotDestructing(maxX uint, maxY uint) ConwayBuffer
	MaxX() uint
	MaxY() uint
	newEmptyBuffer() ConwayBuffer
	nextGeneration(ConwayBuffer) (ConwayBuffer, bool)
	changeSizeDestructing(maxX uint, maxY uint)
	equals(ConwayBuffer) bool
}

type conwayBufferImpl struct {
	maxXVal uint
	maxYVal uint
	vals    []bool
}

func (b *conwayBufferImpl) MaxX() uint {
	return b.maxXVal
}

func (b *conwayBufferImpl) MaxY() uint {
	return b.maxYVal
}

func (b *conwayBufferImpl) equals(other ConwayBuffer) bool {
	sizesOk := b.maxXVal == other.MaxX() && b.maxYVal == other.MaxY()
	if sizesOk {
		for i := uint(0); i < b.maxXVal; i++ {
			for j := uint(0); j < b.maxYVal; j++ {
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

func initBuffer(maxX uint, maxY uint) ConwayBuffer {
	var res = conwayBufferImpl{
		maxXVal: maxX,
		maxYVal: maxY,
		vals:    make([]bool, maxX*maxY)}

	return &res
}

func (bb *conwayBufferImpl) changeSizeDestructing(maxX uint, maxY uint) {
	bb.vals = make([]bool, maxX*maxY)
}

func (bb *conwayBufferImpl) Get(x uint, y uint) bool {
	if x < bb.maxXVal && y < bb.maxYVal {
		return bb.vals[x+y*bb.maxXVal]
	}
	log.Fatal("Invalid range")
	return false
}

func (bb *conwayBufferImpl) Set(x uint, y uint, val bool) {
	if x < bb.maxXVal && y < bb.maxYVal {
		bb.vals[x+y*bb.maxXVal] = val
	} else {
		log.Fatal("Invalid range")
	}
}

func (b *conwayBufferImpl) newEmptyBuffer() ConwayBuffer {
	var res conwayBufferImpl
	res.maxXVal = b.maxXVal
	res.maxYVal = b.maxYVal
	res.vals = make([]bool, b.maxXVal*b.maxYVal)
	return &res
}

func (b *conwayBufferImpl) ChangeSizeNotDestructing(maxX uint, maxY uint) ConwayBuffer {
	if maxX != b.MaxX() || maxY != b.MaxY() {
		newB := initBuffer(maxX, maxY)
		for x := uint(0); x < min(maxX, b.MaxX()); x++ {
			for y := uint(0); y < min(maxY, b.MaxY()); y++ {
				newB.Set(x, y, b.Get(x, y))
			}
		}
		return newB
	} else {
		return b
	}
}

func (b *conwayBufferImpl) countNeighbors(x, y uint) uint {
	var count uint
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := int(x)+dx, int(y)+dy
			if nx >= 0 && nx < int(b.maxXVal) && ny >= 0 && ny < int(b.maxYVal) && b.Get(uint(nx), uint(ny)) {
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
		buffer = initBuffer(b.maxXVal, b.maxYVal)
	}

	isChanged := false

	for y := uint(0); y < b.maxYVal; y++ {
		for x := uint(0); x < b.maxXVal; x++ {
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

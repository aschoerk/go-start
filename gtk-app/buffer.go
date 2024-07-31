package main

import (
	"sync"
)

type Buffer struct {
	data    *[][]bool
	blocked int
	mu      sync.Mutex
}

func (b *Buffer) Update(newData *[][]bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = newData
}

func (b *Buffer) Get() *[][]bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.data
}

func (b *Buffer) ToggleCell(x, y int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if x >= 0 && x < len((*b.data)[0]) && y >= 0 && y < len(*b.data) {
		(*b.data)[y][x] = !(*b.data)[y][x]
	}
}

func (b *Buffer) NextGeneration() {
	b.mu.Lock()
	defer b.mu.Unlock()

	newData := make([][]bool, len((*b.data)))
	for i := range newData {
		newData[i] = make([]bool, len((*b.data)[i]))
	}

	for y := range *b.data {
		for x := range (*b.data)[y] {
			neighbors := b.countNeighbors(x, y)
			if (*b.data)[y][x] {
				newData[y][x] = neighbors == 2 || neighbors == 3
			} else {
				newData[y][x] = neighbors == 3
			}
		}
	}

	bufferHistory[actBufferHistoryIndex] = b.data

	actBufferHistoryIndex = (actBufferHistoryIndex + 1) % MAX_BUFFER_HISTORY

	b.data = &newData
}

func (b *Buffer) countNeighbors(x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < len((*b.data)[0]) && ny >= 0 && ny < len(*b.data) && (*b.data)[ny][nx] {
				count++
			}
		}
	}
	return count
}

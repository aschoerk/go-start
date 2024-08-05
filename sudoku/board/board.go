package sudoku

import (
	"fmt"
	"math"
	"reflect"
)

type SudokuBoard interface {
	Size() uint8
	Get(x uint8, y uint8) uint8
	Set(x uint8, y uint8, val uint8)
	PrintBoard()
	findEmptyCell() (uint8, uint8, bool)
	isValid(num, row, col uint8) bool
	SolveSudoku() bool
	SolveByHeuristic() (bool, *[]SudokuBoard)
	Equals(b SudokuBoard) bool
}

func CreateEmptyBoard(size uint8) SudokuBoard {
	return &sudokuBoardImpl{size, make([]uint8, size*size)}
}

func CreateBoard(arr *[9][9]uint8) SudokuBoard {
	tmp := CreateEmptyBoard(uint8(9))
	for i := uint8(0); i < 9; i++ {
		for j := uint8(0); j < 9; j++ {
			tmp.Set(i, j, arr[i][j])
		}
	}

	return tmp
}

type sudokuBoardImpl struct {
	size uint8
	vals []uint8
}

//go:inline
func (b *sudokuBoardImpl) Size() uint8 {
	return b.size
}

//go:inline
func (b *sudokuBoardImpl) Get(x uint8, y uint8) uint8 {
	return b.vals[y*b.size+x]
}

//go:inline
func (b *sudokuBoardImpl) Set(x uint8, y uint8, val uint8) {
	b.vals[y*b.size+x] = val
}

func (b *sudokuBoardImpl) xY(index uint8) (uint8, uint8) {
	return index % b.size, index / b.size
}

func (b *sudokuBoardImpl) Equals(c SudokuBoard) bool {
	if sudokuBoardImplPtr, ok := c.(*sudokuBoardImpl); ok {
		return reflect.DeepEqual(b.vals, sudokuBoardImplPtr.vals)
	} else {
		return false
	}
}

func (b *sudokuBoardImpl) isFilled() bool {
	for _, val := range b.vals {
		if val == 0 {
			return false
		}
	}
	return true
}

func (b *sudokuBoardImpl) copy() *sudokuBoardImpl {

	res := &sudokuBoardImpl{b.size, make([]uint8, b.size*b.size)}

	copy(res.vals, b.vals)
	return res
}

func (b *sudokuBoardImpl) PrintBoard() {
	for i := uint8(0); i < b.size; i++ {
		if i%3 == 0 && i != 0 {
			fmt.Println("---------------------")
		}
		for j := uint8(0); j < b.size; j++ {
			if j%3 == 0 && j != 0 {
				fmt.Print("| ")
			}
			fmt.Printf("%d ", b.Get(i, j))
		}
		fmt.Println()
	}
}

func (b *sudokuBoardImpl) findEmptyCell() (uint8, uint8, bool) {
	for i := uint8(0); i < b.size; i++ {
		for j := uint8(0); j < b.size; j++ {
			if b.Get(i, j) == 0 {
				return i, j, true
			}
		}
	}
	return math.MaxUint8, math.MaxUint8, false
}

// checks if num can already be found in the row, the col or the square containing the cell row,col
func (b *sudokuBoardImpl) isValid(num, x, y uint8) bool {
	// Check row
	for i := uint8(0); i < b.size; i++ {
		if b.Get(x, i) == num {
			return false
		}
	}

	// Check column
	for i := uint8(0); i < b.size; i++ {
		if b.Get(i, y) == num {
			return false
		}
	}

	// Check 3x3 box
	startX, startY := x-x%3, y-y%3
	for i := uint8(0); i < 3; i++ {
		for j := uint8(0); j < 3; j++ {
			if b.Get(i+startX, j+startY) == num {
				return false
			}
		}
	}

	return true
}

func (b *sudokuBoardImpl) SolveSudoku() bool {
	x, y, found := b.findEmptyCell()
	if !found {
		return true // Puzzle solved
	}

	for num := uint8(1); num <= 9; num++ {
		if b.isValid(num, x, y) {
			b.Set(x, y, num)

			if b.SolveSudoku() {
				return true
			}

			b.Set(x, y, 0) // Backtrack
		}
	}

	return false // No solution exists
}

type possibility struct {
	index   uint8
	entries []uint8
}

func (b *sudokuBoardImpl) createPossibilities() []possibility {
	res := make([]possibility, 0)
	for i := uint8(0); i < b.size*b.size; i++ {
		if b.vals[i] == 0 {
			x, y := b.xY(i)
			entries := make([]uint8, 0)
			for j := uint8(1); j <= b.size; j++ {
				if b.isValid(j, x, y) {
					entries = append(entries, j)
				}
			}
			res = append(res, possibility{i, entries})

		}
	}
	return res
}

func (b *sudokuBoardImpl) getNonUniques() []possibility {

	for {
		leftPossibilities := b.createPossibilities()

		foundUnique := false

		for _, poss := range leftPossibilities {
			if len(poss.entries) == 1 {
				b.vals[poss.index] = poss.entries[0]
				foundUnique = true
			}
		}
		if !foundUnique {
			return leftPossibilities
		}
	}

}

func (b *sudokuBoardImpl) SolveByHeuristic() (bool, *[]SudokuBoard) {

	possibilities := b.getNonUniques()
	if len(possibilities) == 0 {
		return b.isFilled(), nil
	}

	solutions := make([]SudokuBoard, 0)

	for _, possibility := range possibilities {
		for _, val := range possibility.entries {
			tmpB := b.copy()
			tmpB.vals[possibility.index] = val

			res, foundSolutions := tmpB.SolveByHeuristic()

			if res {
				solutions = append(solutions, *foundSolutions...)
			}
		}
	}

	if len(solutions) > 0 {
		return true, &solutions
	} else {
		return false, nil
	}

}

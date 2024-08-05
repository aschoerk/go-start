package main

import (
	"fmt"
	"math"
	"time"

	sudoku "aschoerk.de/sudoku/board"
)

const SIZE = 9

// findEmptyCell finds an empty cell in the Sudoku board
func findEmptyCell(board *[SIZE][SIZE]uint8) (uint8, uint8, bool) {
	for i := uint8(0); i < SIZE; i++ {
		for j := uint8(0); j < SIZE; j++ {
			if board[i][j] == 0 {
				return i, j, true
			}
		}
	}
	return math.MaxUint8, math.MaxUint8, false
}

// isValid checks if it's valid to place a number in a given position
func isValid(board *[SIZE][SIZE]uint8, num, row, col uint8) bool {
	// Check row
	for i := 0; i < SIZE; i++ {
		if board[row][i] == num {
			return false
		}
	}

	// Check column
	for i := 0; i < SIZE; i++ {
		if board[i][col] == num {
			return false
		}
	}

	// Check 3x3 box
	startRow, startCol := row-row%3, col-col%3
	for i := uint8(0); i < 3; i++ {
		for j := uint8(0); j < 3; j++ {
			if board[i+startRow][j+startCol] == num {
				return false
			}
		}
	}

	return true
}

// solveSudoku solves the Sudoku puzzle using backtracking
func solveSudoku(board *[SIZE][SIZE]uint8) bool {
	row, col, found := findEmptyCell(board)
	if !found {
		return true // Puzzle solved
	}

	for num := uint8(1); num <= 9; num++ {
		if isValid(board, num, row, col) {
			board[row][col] = num

			if solveSudoku(board) {
				return true
			}

			board[row][col] = 0 // Backtrack
		}
	}

	return false // No solution exists
}

func byHeuristicsBackTracking(b sudoku.SudokuBoard) bool {
	start := time.Now()

	defer func() {
		// Record the end time
		end := time.Now()

		// Calculate the duration
		duration := end.Sub(start)

		// Print the duration
		fmt.Printf("Processing time Heuristics: %s\n", duration)
	}()

	solved, _ := b.SolveByHeuristic()

	return solved

}

func byPureBacktracking(b sudoku.SudokuBoard) bool {
	start := time.Now()

	defer func() {
		// Record the end time
		end := time.Now()

		// Calculate the duration
		duration := end.Sub(start)

		// Print the duration
		fmt.Printf("Processing time Backtracking: %s\n", duration)
	}()
	return b.SolveSudoku()

}

func byArray(board *[SIZE][SIZE]uint8) bool {
	// Record the start time
	start := time.Now()

	defer func() {
		// Record the end time
		end := time.Now()

		// Calculate the duration
		duration := end.Sub(start)

		// Print the duration
		fmt.Printf("Processing time Backtracking Arrays: %s\n", duration)
	}()

	return solveSudoku(board)

}

func main() {
	board := [SIZE][SIZE]uint8{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 0, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 0, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	b := sudoku.CreateBoard(&board)

	fmt.Println("Sudoku Puzzle:")

	sudoku.CreateBoard(&board).PrintBoard()

	// if byHeuristicsBackTracking(b) {
	// 	fmt.Println("\nSolved Sudoku:")
	// 	b.PrintBoard()
	// }

	if byPureBacktracking(b) {
		fmt.Println("\nSolved Sudoku:")
		b.PrintBoard()
	}

	if byArray(&board) {
		fmt.Println("\nSolved Sudoku:")
		sudoku.CreateBoard(&board).PrintBoard()
	} else {
		fmt.Println("\nNo solution exists")
	}

}

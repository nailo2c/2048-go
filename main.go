package main

import (
	"os"

	"github.com/nsf/termbox-go"
)

const boardLen = 4

var gameFieldEndY = 0

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	board := initBoard(boardLen)
	drawGameField(board)
}

func initBoard(len int) [][]int {
	board := make([][]int, len)
	for i := range board {
		board[i] = make([]int, len)
	}
	return board
}

func drawGameField(board [][]int) {
	putNextNumber(board)
}

func putNextNumber(board [][]int) {
	emptyCells := findEmptyCells(board)
	if len(emptyCells) <= 0 {
		gameOver()
	}
}

func findEmptyCells(board [][]int) []int {
	emptyCells := []int{}
	for i, row := range board {
		for j, cell := range row {
			if cell == 0 {
				emptyCells = append(emptyCells, i*len(board)+j)
			}
		}
	}
	return emptyCells
}

func gameOver() {
	printTerminal(0, gameFieldEndY, []string{"Game Over!"})
	termbox.SetCursor(0, gameFieldEndY+1)
	termbox.Flush()
	os.Exit(0)
}

func printTerminal(startX, startY int, strs []string) int {
	for y, str := range strs {
		for x, ch := range str {
			termbox.SetCell(startX+x, startY+y, ch, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
	return startY + len(strs)
}

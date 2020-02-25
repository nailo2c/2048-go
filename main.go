package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

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
	putNextNumber(board)
	boardStartY := printTerminal(0, 0, []string{"Game 2048", ""})
	boardEndY := drawBoard(0, boardStartY, board)
	gameFieldEndY = printTerminal(0, boardEndY, []string{"Esc ←↑↓→", ""})
}

func putNextNumber(board [][]int) {
	emptyCells := findEmptyCells(board)
	if len(emptyCells) <= 0 {
		gameOver()
	}
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)
	emptyCell := emptyCells[rnd.Intn(len(emptyCells))]
	board[emptyCell/len(board)][emptyCell%len(board)] = 2
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

func drawBoard(startX, startY int, board [][]int) int {
	strs := []string{}
	for _, row := range board {
		var str string
		for _, cell := range row {
			str += fmt.Sprintf("%4d", cell)
		}
		str = strings.Replace(str, "0", ".", -1)
		strs = append(strs, str, "")
	}

	printTerminal(0, startY, strs)
	termbox.Flush()
	return startY + len(strs)
}

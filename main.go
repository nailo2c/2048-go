package main

import (
	"fmt"
	"math/rand"
	"os"
	"reflect"
	"time"

	"github.com/nsf/termbox-go"
)

const boardLen = 4

var boardStartY int
var gameFieldEndY int

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	board := initBoard(boardLen)
	drawGameField(board)
	startGame(board)
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
	boardStartY = printTerminal(0, 0, []string{"Game 2048", ""})
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
	termbox.Close()
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
			if cell == '0' {
				cell = '.'
			}
			str += fmt.Sprintf("%5d", cell)
		}
		strs = append(strs, str, "")
	}

	printTerminal(0, startY, strs)
	termbox.Flush()
	return startY + len(strs)
}

func startGame(board [][]int) {
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventError {
			if event.Err != nil {
				panic(event.Err)
			}
		}
		if event.Type == termbox.EventKey {
			prevBoard := copyBoard(board)
			switch event.Key {
			case termbox.KeyEsc:
				termbox.SetCursor(0, gameFieldEndY)
				termbox.Flush()
				os.Exit(0)
			case termbox.KeyArrowDown:
				board = rotateBoard(board, false)
				board = slideLeft(board)
				board = rotateBoard(board, true)
				notCange := reflect.DeepEqual(prevBoard, board)
				checkAndRefreshBoard(board, notCange)
			case termbox.KeyArrowLeft:
				board = slideLeft(board)
				notCange := reflect.DeepEqual(prevBoard, board)
				checkAndRefreshBoard(board, notCange)
			case termbox.KeyArrowRight:
				board = rotateBoard(board, true)
				board = rotateBoard(board, true)
				board = slideLeft(board)
				board = rotateBoard(board, false)
				board = rotateBoard(board, false)
				notCange := reflect.DeepEqual(prevBoard, board)
				checkAndRefreshBoard(board, notCange)
			case termbox.KeyArrowUp:
				board = rotateBoard(board, true)
				board = slideLeft(board)
				board = rotateBoard(board, false)
				notCange := reflect.DeepEqual(prevBoard, board)
				checkAndRefreshBoard(board, notCange)
			}
		}
	}
}

func rotateBoard(board [][]int, counterClockWise bool) [][]int {
	rotateBoard := make([][]int, len(board))
	for i, row := range board {
		rotateBoard[i] = make([]int, len(row))
		for j := range row {
			if counterClockWise {
				rotateBoard[i][j] = board[j][len(board)-i-1]
			} else {
				rotateBoard[i][j] = board[len(board)-j-1][i]
			}
		}
	}
	return rotateBoard
}

func slideLeft(board [][]int) [][]int {
	for _, row := range board {
		stopMerge := 0
		for j := 1; j < len(row); j++ {
			if row[j] != 0 {
				for k := j; k > stopMerge; k-- {
					if row[k-1] == 0 {
						row[k-1] = row[k]
						row[k] = 0
					} else if row[k-1] == row[k] {
						row[k-1] += row[k]
						row[k] = 0
						stopMerge = k
						break
					} else {
						break
					}
				}
			}
		}
	}
	return board
}

func checkAndRefreshBoard(board [][]int, boardNotChange bool) {
	checkWinner(board)
	if !boardNotChange {
		putNextNumber(board)
	}
	drawBoard(0, boardStartY, board)
}

func checkWinner(board [][]int) {
	for _, row := range board {
		for _, cell := range row {
			if cell == 2048 {
				drawBoard(0, boardStartY, board)
				gameWin()
			}
		}
	}
}

func gameWin() {
	printTerminal(0, gameFieldEndY, []string{"You Won!"})
	termbox.SetCursor(0, gameFieldEndY+1)
	termbox.Flush()
	os.Exit(0)
}

func copyBoard(board [][]int) [][]int {
	newBoard := make([][]int, len(board))
	for i, row := range board {
		newBoard[i] = make([]int, len(row))
		copy(newBoard[i], row)
	}

	return newBoard
}

package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	BoardWidth  = 10
	BoardHeight = 20
)

type Piece struct {
	x     int
	y     int
	shape [][]int
}

func newSquarePiece() Piece {
	shape := [][]int{
		{1, 1},
		{1, 1},
	}
	return Piece{
		x:     4,
		y:     0,
		shape: shape,
	}
}

func drawBoard(board [][]int, piece Piece) {
	for x := 0; x < BoardHeight; x++ {
		for y := 0; y < BoardWidth; y++ {
			filled := (board[x][y] == 1)
			if !filled && pieceAt(x, y, piece) {
				filled = true
			}
			if filled {
				fmt.Print("# ")
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Print("\r\n") // <-- changed from fmt.Println()
	}
}

func collapseTheRow(board [][]int, row int) {
	for y := row; y > 0; y-- {
		for x := BoardWidth - 1; x >= 0; x-- {
			board[y][x] = board[y-1][x]
		}
	}

	for i := 0; i < BoardWidth; i++ {
		board[0][i] = 0
	}
}

func removeIfLineIsPresent(board [][]int) {
	for row := BoardHeight - 1; row >= 0; row-- {
		cleanRow := true
		for col := BoardWidth - 1; col >= 0; col-- {
			if board[row][col] != 1 {
				cleanRow = false
			}
		}
		if cleanRow {
			collapseTheRow(board, row)
			row++
		}
	}
}

func pieceAt(row int, col int, piece Piece) bool {
	py := row - piece.y
	px := col - piece.x
	if py < 0 || py >= len(piece.shape) || px < 0 || px >= len(piece.shape[0]) {
		return false
	}
	return piece.shape[py][px] == 1
}

func placePieceOnBoard(board [][]int, piece Piece) {
	for py := 0; py < len(piece.shape); py++ {
		for px := 0; px < len(piece.shape[py]); px++ {
			if piece.shape[py][px] == 1 {
				board[piece.y+py][piece.x+px] = 1
			}
		}
	}
}

func clearBoard(board [][]int) {
	for x := 0; x < BoardHeight; x++ {
		for y := 0; y < BoardWidth; y++ {
			board[x][y] = 0
		}
	}
}

func canMoveSide(board [][]int, piece Piece, dx int) bool {
	newX := piece.x + dx
	if newX < 0 || newX+len(piece.shape[0]) > BoardWidth {
		return false
	}
	for py := 0; py < len(piece.shape); py++ {
		for px := 0; px < len(piece.shape[0]); px++ {
			if piece.shape[py][px] != 1 {
				continue
			}
			boardRow := piece.y + py
			boardCol := newX + px
			if board[boardRow][boardCol] == 1 {
				return false
			}
		}
	}
	return true
}

func canMove(board [][]int, piece Piece) bool {
	if piece.y+len(piece.shape) >= BoardHeight {
		return false
	}
	for py := 0; py < len(piece.shape); py++ {
		for px := 0; px < len(piece.shape[0]); px++ {
			if piece.shape[py][px] != 1 {
				continue
			}
			boardRowBelow := piece.y + py + 1
			boardCol := piece.x + px
			if board[boardRowBelow][boardCol] == 1 {
				return false
			}
		}
	}
	return true
}

func captureInput(ch chan string) {
    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        panic(err)
    }
    defer term.Restore(int(os.Stdin.Fd()), oldState)

    buf := make([]byte, 1)
    for {
        n, err := os.Stdin.Read(buf)
        if err != nil || n == 0 {
            continue
        }
        if buf[0] == 3 { // Ctrl+C -- was outside loop before, never executed
            term.Restore(int(os.Stdin.Fd()), oldState)
            os.Exit(0)
        }
        ch <- string(buf[:n])
    }
}

func main() {
	input := make(chan string)
	go captureInput(input)

	board := make([][]int, BoardHeight)
	for i := 0; i < BoardHeight; i++ {
		board[i] = make([]int, BoardWidth)
	}

	piece := newSquarePiece()
	key := ""
	for {
		// clearBoard(board)

		fmt.Print("\033[H\033[2J")
		removeIfLineIsPresent(board)
		drawBoard(board, piece)

		select {
		case key = <-input:
		default:
	      key = ""
		}

		if key == "a" && canMoveSide(board, piece, -1) {
			piece.x--
		}

		if key == "d" && canMoveSide(board, piece, 1) {
			piece.x++
		}

		if key == "s" && canMove(board, piece) {
			piece.y++
		}
		if canMove(board, piece) {
			piece.y++
		} else {
			placePieceOnBoard(board, piece)
			piece = newSquarePiece()
			
		}
		key = ""
		time.Sleep(300 * time.Millisecond)
	}
}

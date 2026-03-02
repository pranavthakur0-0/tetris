package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"golang.org/x/term"
)

const (
	BoardWidth  = 10
	BoardHeight = 20
)

var pieces = []Piece{
	{x: 4, y: 0, shape: [][]int{{1, 1}, {1, 1}}},         // O
	{x: 3, y: 0, shape: [][]int{{1, 1, 1, 1}}},           // I
	{x: 4, y: 0, shape: [][]int{{0, 1, 0}, {1, 1, 1}}},   // T
	{x: 4, y: 0, shape: [][]int{{1, 0}, {1, 0}, {1, 1}}}, // L
	{x: 4, y: 0, shape: [][]int{{0, 1}, {0, 1}, {1, 1}}}, // J
	{x: 4, y: 0, shape: [][]int{{0, 1, 1}, {1, 1, 0}}},   // S
	{x: 4, y: 0, shape: [][]int{{1, 1, 0}, {0, 1, 1}}},   // Z
}

type Piece struct {
	x     int
	y     int
	shape [][]int
}

func newPiece() Piece {
	p := pieces[rand.Intn(len(pieces))]
	return Piece{x:p.x, y:p.y, shape: p.shape}
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
		fmt.Print("\r\n")
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

    buf := make([]byte, 3)
    for {
        n, _ := os.Stdin.Read(buf)
        if n == 3 && buf[0] == '\x1b' && buf[1] == '[' {
            switch buf[2] {
            case 'A':
                ch <- "up"
            case 'B':
                ch <- "down"
            case 'C':
                ch <- "right"
            case 'D':
                ch <- "left"
            }
        } else if n == 1 {
            ch <- string(buf[:1])
        }
    }
}


func isGameOver(board [][]int) bool {
    for y := 0; y < BoardWidth; y++ {
        if board[0][y] == 1 {
            return true
        }
    }
    return false
}

func Rotate(piece Piece) Piece {
    rows := len(piece.shape)
    cols := len(piece.shape[0])
    rotate := make([][]int, cols)
    for i := range rotate {
        rotate[i] = make([]int, rows)
    }
    for row := 0; row < rows; row++ {
        for col := 0; col < cols; col++ {
            rotate[col][rows-1-row] = piece.shape[row][col]
        }
    }
    return Piece{x: piece.x, y: piece.y, shape: rotate}
}

func canRotate(board [][]int, rotated Piece) bool {
    if rotated.x < 0 || rotated.x+len(rotated.shape[0]) > BoardWidth {
        return false
    }
    if rotated.y+len(rotated.shape) > BoardHeight {
        return false
    }
    for row := 0; row < len(rotated.shape); row++ {
        for col := 0; col < len(rotated.shape[0]); col++ {
            if rotated.shape[row][col] != 1 {
                continue
            }
            if board[rotated.y+row][rotated.x+col] == 1 {
                return false
            }
        }
    }
    return true
}

func main() {
	input := make(chan string)
	go captureInput(input)

	board := make([][]int, BoardHeight)
	for i := 0; i < BoardHeight; i++ {
		board[i] = make([]int, BoardWidth)
	}

	piece := newPiece()
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	for {

		select {
		case key := <-input:
			switch key {
			case "left", "a":
				if canMoveSide(board, piece, -1) {
					piece.x--
				}
			case "right", "d":
				if canMoveSide(board, piece, 1) {
					piece.x++
				}
			case "up", "r":
				rotated := Rotate(piece)
				if canRotate(board, rotated) {
					piece = rotated
				}
			case "down", "s":
				if canMove(board, piece) {
					piece.y++
				}
			case "q":
				return
			}
		case <-ticker.C:
			if canMove(board, piece) {
				piece.y++
			} else {
				placePieceOnBoard(board, piece)
				if isGameOver(board) {
					fmt.Print("\033[2J\033[H")
					fmt.Print("Game over!\r\n")
					return  
				}
				piece = newPiece()

			}
		}
		fmt.Print("\033[2J\033[H")
		removeIfLineIsPresent(board)
		drawBoard(board, piece)
	}
}

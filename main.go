package main

import (
	"fmt"
	"time"
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
		fmt.Println()
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

func main() {
	board := make([][]int, BoardHeight)
	for i := 0; i < BoardHeight; i++ {
		board[i] = make([]int, BoardWidth)
	}

	piece := newSquarePiece()
	for {
		// clearBoard(board)

		fmt.Print("\033[H\033[2J")
		removeIfLineIsPresent(board)
		drawBoard(board, piece)

		var input string
		fmt.Scanln(&input)

		if input == "a" && canMoveSide(board, piece, -1) {
			piece.x--
		}

		if input == "d" && canMoveSide(board, piece, 1) {
			piece.x++
		}

		if input == "s" && canMove(board, piece) {
			piece.y++
		}
		if canMove(board, piece) {
			piece.y++
		} else {
			placePieceOnBoard(board, piece)
			piece = newSquarePiece()
		}
		time.Sleep(300 * time.Millisecond)
	}
}

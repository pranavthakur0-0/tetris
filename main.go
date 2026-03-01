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
	x int
	y int
	shape [][]int
}


func newSquarePiece() Piece {
	shape := [][]int{
		{0, 1, 0},
		{1, 1, 1},
	}
	return Piece{
		x : 4,
		y : 0,
		shape: shape,
	}
}




func drawBoard(board [][]int){
	for x := 0; x < BoardHeight; x++ {
		for y := 0; y < BoardWidth; y++{
			if board[x][y] == 0 {
                fmt.Print(". ")
            } else {
                fmt.Print("# ")
            }
		}
		fmt.Println()
	}
}

func placePieceOnBoard(board [][]int, piece Piece){
	for py := 0; py < len(piece.shape); py++ {
		for px := 0; px < len(piece.shape[py]); px++ {
			if piece.shape[py][px] == 1 {
                board[piece.y+py][piece.x+px] = 1
            }
		}
	}
}


func clearBoard(board [][]int){
	for x := 0; x < BoardHeight; x++ {
		for y := 0; y < BoardWidth; y++{
			 board[x][y] = 0 
		}
	}
}

func canMoveSide(piece Piece, dx int) bool {
	newX := piece.x + dx
	return newX >= 0 && newX+len(piece.shape[0]) <= BoardWidth
}

func canMove(piece Piece) bool {
	return piece.y+len(piece.shape) < BoardHeight
}



func main() {
	board := make([][]int, BoardHeight)
	for i := 0; i < BoardHeight; i++{
		board[i] = make([]int,BoardWidth)
	}
	piece := newSquarePiece();
	for {
		// clearBoard(board)
		
		fmt.Print("\033[H\033[2J")
		placePieceOnBoard(board, piece)
		drawBoard(board)

		var input string 
		fmt.Scanln(&input)

		if input == "a" && canMoveSide(piece, -1) {
			piece.x--
		}
	
		if input == "d" && canMoveSide(piece, 1) {
			piece.x++
		}
	
		if input == "s" && canMove(piece) {
			piece.y++
		}
		if canMove(piece) {
			piece.y++
		}else{
			piece = newSquarePiece()
		}

		// piece.y++clear
		time.Sleep(300 * time.Millisecond)
	}
	
}


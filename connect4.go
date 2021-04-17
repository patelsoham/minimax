package main

import (
	"fmt"
)

const (
	EMPTY       = iota
	P1          = iota
	P2          = iota
	EMPTY_COLOR = "\033[1;37m%s\033[0m" //White
	P1_COLOR    = "\033[1;31m%s\033[0m" //Red
	P2_COLOR    = "\033[1;36m%s\033[0m" //Teal
)

var colors = []string{EMPTY_COLOR, P1_COLOR, P2_COLOR}

type BitBoard struct {
	boards  []int64 //two 64 bit integers (longs) for each player's board (lower 6 bits for every byte represent a col)
	rows    int
	cols    int
	heights []int //each columns height (can eventually be converted to int64-> byte per col)
}

type Board struct {
	board   [][]int
	heights []int //each columns height (can eventually be converted to int64-> byte per col)
}

func getBitBoard(rows int, cols int) *BitBoard {
	return &BitBoard{make([]int64, 2), rows, cols, make([]int, cols)}
}

func (b *BitBoard) modBoard(col int, player int, delta int) {
	if delta > 0 {
		b.boards[player>>1] ^= (1 << uint(b.heights[col]))
		b.heights[col] += delta
	} else {
		b.heights[col] += delta
		b.boards[player>>1] ^= (1 << uint(b.heights[col]))
	}
	if b.heights[col] > b.rows || b.heights[col] < 0 {
		fmt.Printf("Invalid Height at col %d: %d\n", col, b.heights[col])
		panic(b.heights[col])
	}
}

//TODO: BitBoard printing
func (b *BitBoard) printBoard() {
	fmt.Printf("BitBoard\n")
	for i := 0; i < b.rows; i++ {
		for j := 0; j < b.cols; j++ {
			fmt.Printf(colors[1], "O ")
		}
		fmt.Printf("\n")
	}
}

func getBoard(row int, col int) *Board {
	arr := make([][]int, row)
	for i := range arr {
		arr[i] = make([]int, col)
	}
	return &Board{arr, make([]int, col)}
}

func (b *Board) modBoard(col int, player int, delta int) {
	if delta > 0 {
		b.board[len(b.board)-b.heights[col]-1][col] = player
		b.heights[col] += delta
	} else {
		b.heights[col] += delta
		b.board[len(b.board)-b.heights[col]-1][col] = player
	}
	if b.heights[col] > len(b.board) || b.heights[col] < 0 {
		fmt.Printf("Invalid Height at col %d: %d\n", col, b.heights[col])
		panic(b.heights[col])
	}
}

func (b *Board) printBoard() {
	fmt.Printf("2D Array Board\n")
	for i := range b.board {
		for j := range b.board[i] {
			fmt.Printf(colors[b.board[i][j]], "O ")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Heights for each column \n")
	fmt.Println(b.heights)
}

func hello() {
	//board := getBoard(6, 7)
	board := getBitBoard(6, 7)
	player := 1
	for j := 0; j < board.rows; j++ {
		for i := 0; i < board.cols; i++ {
			board.modBoard(i, player, 1)
			board.modBoard(i, 0, -1)
		}
		player ^= 3
	}
	board.printBoard()
}

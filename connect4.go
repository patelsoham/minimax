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
	boards              []int64 //two 64 bit integers (longs) for each player's board (lower 6 bits for every byte represent a col)
	rows, cols, heights int
}

type Board struct {
	board   [][]int //height 0 @ row n-1, height 1 @ row n-2, ..., height n @ row 0.
	heights int     //each columns height (can eventually be converted to int64-> byte per col)
}

func getBitBoard(rows int, cols int) *BitBoard {
	return &BitBoard{make([]int64, 2), rows, cols, 0}
}

func (b *BitBoard) modBoard(col int, player int, delta int) {
	cur_height := ((0xF << uint(col*4)) & b.heights) >> uint(col*4)
	//fmt.Printf("Before Placement: Col %d Player %d CurHeight %d\n", col, player, cur_height)
	if delta > 0 {
		b.boards[player>>1] ^= (1 << uint(cur_height+(col*8)))
		cur_height += delta
	} else {
		cur_height += delta
		b.boards[player>>1] ^= (1 << uint(cur_height+(col*8)))
	}
	//fmt.Printf("After Placement: Col %d Player %d CurHeight %d\n\n", col, player, cur_height)
	if cur_height > b.rows || cur_height < 0 {
		fmt.Printf("Invalid Height at col %d: %d\n", col, cur_height)
		panic(cur_height)
	}
	b.heights &= ^(0xF << uint(col*4))
	b.heights |= (cur_height << uint(col*4))
}

//Took this code from https://github.com/denkspuren/BitboardC4/blob/master/BitboardDesign.md
func (b *BitBoard) hasWon(player int) bool {
	var directions = []uint{1, 6, 7, 8}
	cur_board := b.boards[player>>1]
	var masked_board int64
	for i := range directions {
		masked_board = cur_board & (cur_board >> directions[i])
		if (masked_board & (masked_board >> (2 * directions[i]))) != 0 {
			return true
		}
	}
	return false
}

func (b *BitBoard) gameState() (int, int) {
	if b.hasWon(P1) {
		return 1, P1
	} else if b.hasWon(P2) {
		return 1, P2
	} else if len(movesAvailable(b.heights, b.rows, b.cols)) == 0 {
		return 0, 0
	}
	return -1, -1
}

//TODO: BitBoard printing
func (b *BitBoard) printBoard() {
	fmt.Printf("BitBoard\n")
	for i := b.rows - 1; i >= 0; i-- {
		for j := 0; j < b.cols; j++ {
			cur_cell_1 := ((((b.boards[0]) & (0xFF << uint(j*8))) >> uint(j*8)) >> uint(i)) & 1
			cur_cell_2 := ((((b.boards[1]) & (0xFF << uint(j*8))) >> uint(j*8)) >> uint(i)) & 1
			//fmt.Printf("Row %d, Col %d, p1 %d p2 %d ind %d\n", i, j, cur_cell_1, cur_cell_2, cur_cell_1+(cur_cell_2*2))
			fmt.Printf(colors[cur_cell_1+(cur_cell_2*2)], "O ")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("Heights for each column \n")
	for i := 0; i < b.cols; i++ {
		fmt.Printf("%d ", ((0xF<<uint(i*4))&b.heights)>>uint(i*4))
	}
	fmt.Printf("\n")
}

func getBoard(row int, col int) *Board {
	arr := make([][]int, row)
	for i := range arr {
		arr[i] = make([]int, col)
	}
	return &Board{arr, 0}
}

func (b *Board) modBoard(col int, player int, delta int) {
	cur_height := ((0xF << uint(col*4)) & b.heights) >> uint(col*4)
	//fmt.Printf("Before Placement: Col %d Player %d CurHeight %d\n", col, player, cur_height)
	if delta > 0 {
		b.board[len(b.board)-cur_height-1][col] = player
		cur_height += delta
	} else {
		cur_height += delta
		b.board[len(b.board)-cur_height-1][col] = player
	}
	//fmt.Printf("After Placement: Col %d Player %d CurHeight %d\n\n", col, player, cur_height)
	if cur_height > len(b.board) || cur_height < 0 {
		fmt.Printf("Invalid Height at col %d: %d\n", col, cur_height)
		panic(cur_height)
	}
	b.heights &= ^(0xF << uint(col*4))
	b.heights |= (cur_height << uint(col*4))
}

func (b *Board) hasWon(player int) bool {
	//Horizontal

	//Vertical

	//
	return false
}

func (b *Board) gameState() (int, int) {
	if b.hasWon(P1) {
		return 1, P1
	} else if b.hasWon(P2) {
		return 1, P2
	} else if len(movesAvailable(b.heights, len(b.board), len(b.board[0]))) == 0 {
		return 0, 0
	}
	return -1, -1
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
	for i := 0; i < len(b.board[0]); i++ {
		fmt.Printf("%d ", ((0xF<<uint(i*4))&b.heights)>>uint(i*4))
	}
	fmt.Printf("\n")
}

func movesAvailable(heights int, height_lim int, max_moves int) []int {
	//get slice of columns with spaces based on heights
	moves := make([]int, 0)
	for i := 0; i < max_moves; i++ {
		if ((0xF<<uint(i*4))&heights)>>uint(i*4) < height_lim {
			moves = append(moves, i)
		}
	}
	return moves
}

func hello() {
	player := 1
	//board := getBoard(6, 7)
	board := getBitBoard(6, 7)
	for j := 0; j < 7; j++ {
		for i := 0; i < 6; i++ {
			board.modBoard(j, player, 1)
		}
		player ^= 3

	}
	fmt.Println(board.hasWon(1))
	fmt.Println(board.gameState())
	board.printBoard()
}

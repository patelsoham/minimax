package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	EMPTY       = iota
	P1          = iota
	P2          = iota
	EMPTY_COLOR = "\033[1;37m%s\033[0m" //White
	P1_COLOR    = "\033[1;31m%s\033[0m" //Red
	P2_COLOR    = "\033[1;36m%s\033[0m" //Teal
	MAX         = 2147483647
	MIN         = -2147483648
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

func (b *BitBoard) copyBoard() *BitBoard {
	new_boards := make([]int64, 2)
	copy(new_boards, b.boards)
	return &BitBoard{new_boards, b.rows, b.cols, b.heights}
}

func (b *BitBoard) modBoard(col int, player int, delta int) {
	cur_height := ((0xF << uint(col*4)) & b.heights) >> uint(col*4)
	//fmt.Printf("Before Placement: Col %d Player %d CurHeight %d\n", col, player, cur_height)
	if delta > 0 {
		b.boards[player>>1] ^= (1 << uint(cur_height+(col*7)))
		cur_height += delta
	} else {
		cur_height += delta
		b.boards[player>>1] ^= (1 << uint(cur_height+(col*7)))
	}
	//fmt.Printf("After Placement: Col %d Player %d CurHeight %d\n\n", col, player, cur_height)
	if cur_height > b.rows || cur_height < 0 {
		fmt.Printf("Invalid Height at col %d: %d made by player %d\n", col, cur_height, player)
		b.printBoard()
		panic(cur_height)
	}
	b.heights &= ^(0xF << uint(col*4))
	b.heights |= (cur_height << uint(col*4))
}

//Took this optimization of hasWon from https://github.com/denkspuren/BitboardC4/blob/master/BitboardDesign.md
func (b *BitBoard) hasWon(player int) bool {
	var directions = []uint{1, 7, 6, 8}
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
			cur_cell_1 := ((((b.boards[0]) & (0xFF << uint(j*7))) >> uint(j*7)) >> uint(i)) & 1
			cur_cell_2 := ((((b.boards[1]) & (0xFF << uint(j*7))) >> uint(j*7)) >> uint(i)) & 1
			//fmt.Printf("Row %d, Col %d, p1 %d p2 %d ind %d\n", i, j, cur_cell_1, cur_cell_2, cur_cell_1+(cur_cell_2*2))
			fmt.Printf(colors[cur_cell_1+(cur_cell_2*2)], "O ")
		}
		fmt.Printf("\n")
	}
	fmt.Printf("%.64b\n %.64b\n", b.boards[0], b.boards[1])
	b.getHeights()
}

func (b *BitBoard) getHeights() {
	fmt.Printf("Heights for each column \n")
	for i := 0; i < b.cols; i++ {
		fmt.Printf("%d ", ((0xF<<uint(i*4))&b.heights)>>uint(i*4))
	}
	fmt.Printf("\n")
}

func countBits(n int64) int {
	count := 0
	for (n > 0) {
		n &= (n-1)
		count++
	}

	return count
}

func scoreWindow(window []int64, player int) int {
	score := 0
	opp_player := P1
	if player == P1 {
		opp_player = P2
	}

	player_count := countBits(window[player>>1]) 
	opp_player_count := countBits(window[opp_player>>1])
	empty_count := 4 - player_count - opp_player_count
	
	if player_count == 4 {
		score+= 100
	} else if player_count == 3 && empty_count == 1 {
		score+= 5
	} else if player_count == 2 && empty_count == 2 {
		score+= 2
	}

	if opp_player_count == 3 && empty_count == 1 {
		score-= 4
	} 	
	
	return score
}

// logic obtained from https://github.com/KeithGalli/Connect4-Python/
func (b *BitBoard) scoreBoard(player int) int {
	score := 0
	opp_player := P1
	if player == P1 {
		opp_player = P2
	}

	center_board := (b.boards[player >> 1] >> 21) & (0x7F)
	center_count := countBits(center_board)
	score+= center_count * 3
	
	// next we will score windows of 4 spots at a time in the rows, columns and diagonals to calculate our score

	// score rows
	for r := uint(0); r < 6; r++ {
		windows := make([]int64, 2)
		for c := uint(0); c < 4; c++ {
			curr_board_player := b.boards[player>>1]
			curr_board_opp := b.boards[opp_player>>1]

			bit_pos := (1 << (r+ 7*c)) | (1 << (r + 7*(c+1))) | (1 << (r + 7*(c+2))) | (1 << (r + 7*(c+3)))
			 
			curr_board_player &= int64(bit_pos)
			curr_board_opp &= int64(bit_pos)

			windows[player>>1] = curr_board_player
			windows[opp_player>>1] = curr_board_opp

			score+= scoreWindow(windows, player)
		}
	}

	// score cols
	for c := uint(0); c < 7; c++ {
		windows := make([]int64, 2)
		for r := uint(0); r < 4; r++ {
			curr_board_player := b.boards[player>>1]
			curr_board_opp := b.boards[opp_player>>1]

			bit_pos := (1 << (c*7 + r)) | (1 << (c*7 + r+1)) | (1 << (c*7 + r+2)) | (1 << (c*7 + r+3))
			 
			curr_board_player &= int64(bit_pos)
			curr_board_opp &= int64(bit_pos)

			windows[player>>1] = curr_board_player
			windows[opp_player>>1] = curr_board_opp

			score+= scoreWindow(windows, player)
		}
	}

	// score + diagonal
	for r := uint(0); r < 4; r++ {
		for c := uint(0); c < 4; c++ {

		}
	}

	// score - diagonal
	for r := uint(0); r < 4; r++ {
		for c := uint(0); c < 4; c++ {

		}
	}

	return score
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
	for c := 0; c < len(b.board[0]) - 3; c++ {
		for r := 0; r < len(b.board); r++ {
			if b.board[r][c] == player && b.board[r][c+1] == player && b.board[r][c+2] == player && b.board[r][c+3] == player {
				fmt.Println("horizontal win")
				return true
			}
		} 
	}
	//Vertical
	for r := 0; r < len(b.board) - 3; r++ {
		for c := 0; c < len(b.board[0]); c++ {
			if b.board[r][c] == player && b.board[r+1][c] == player && b.board[r+2][c] == player && b.board[r+3][c] == player {
				fmt.Println("vertical win")
				return true
			}
		} 
	}
	//Ascending Diagonals (L->R), Descending Diagonals (R->L)
	for r := 3; r < len(b.board); r++ {
		for c := 0; c < len(b.board[0]) - 3; c++ {
			if b.board[r][c] == player && b.board[r-1][c+1] == player && b.board[r-2][c+2] == player && b.board[r-3][c+3] == player {
				fmt.Println("ascending diagonals")
				return true
			}
		}
	}
	//Descending Diagonals (L->R), Ascending Diagonals (R->L)
	for r := 0; r < len(b.board) - 3; r++ {
		for c := 0; c < len(b.board[0]) - 3; c++ {
			if b.board[r][c] == player && b.board[r+1][c+1] == player && b.board[r+2][c+2] == player && b.board[r+3][c+3] == player {
				fmt.Println("descending diagonals")
				return true
			}
		}
	}
	
	return false
}

func (b *Board) gameState() (int, int) {
	if b.hasWon(P1) {
		return MAX, P1
	} else if b.hasWon(P2) {
		return MIN, P2
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

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func hello() {
	player := 1
	//board := getBoard(6, 7)
	rand.Seed(time.Now().UnixNano())
	board := getBitBoard(6, 7)
	for i := 0; i < 42; i++ {
		avail_moves := movesAvailable(board.heights, board.rows, board.cols)
		cur_col := rand.Intn(len(avail_moves))
		board.modBoard(avail_moves[cur_col], player, 1)
		player ^= 3
	}
	nb := board.copyBoard()
	nb.printBoard()
	fmt.Println(board.hasWon(1))
	fmt.Println(board.gameState())
	board.printBoard()
}

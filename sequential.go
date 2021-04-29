package main

import (
	"fmt"
	"math/rand"
	"time"
)

type move struct {
	val int
	col int
}

var prev_moves []move
var count int64 = 0
var per_10_mil time.Duration
var st time.Time

func seq_minimax(b *BitBoard, player int, depth int) (int, int) {
	game_res, player_res := b.gameState()
	if game_res != -1 && player_res != -1 {
		count += 1
		if count%10000000 == 0 {
			new_st := time.Now()
			per_10_mil = new_st.Sub(st)
			fmt.Printf("NAIVE: 10 million game states have been considered in %.5f, a total of %d game states considered\n", per_10_mil.Seconds(), count)
			st = new_st
		}
		return game_res, player_res
	}
	if depth == 0 {
		// we've reached the max tree depth, so use heuristic to give arbitrary score
		return b.scoreBoard(player), 0
	}
	avail_moves := movesAvailable(b.heights, b.rows, b.cols)
	if player == P1 {
		opt_val := MIN
		opt_move := rand.Intn(len(avail_moves))
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax(b, player^3, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val > opt_val {
				opt_val = val
				opt_move = i
			}
		}
		return opt_val, opt_move
	} else if player == P2 {
		opt_val := MAX
		opt_move := rand.Intn(len(avail_moves))
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax(b, (player ^ 3), depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val < opt_val {
				opt_val = val
				opt_move = i
			}
		}
		return opt_val, opt_move
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func seq_minimax_ab(b *BitBoard, player int, alpha int, beta int, depth int) (int, int) {
	game_res, player_res := b.gameState()
	if game_res != -1 && player_res != -1 {
		count += 1
		if count%10000000 == 0 {
			new_st := time.Now()
			per_10_mil = new_st.Sub(st)
			fmt.Printf("ALPHA-BETA: 10 million game states have been considered in %.5f, a total of %d game states considered\n", per_10_mil.Seconds(), count)
			st = new_st
		}
		return game_res, player_res
	}
	if player == P1 {
		opt_val := MIN
		avail_moves := movesAvailable(b.heights, b.rows, b.cols)
		opt_move := rand.Intn(len(avail_moves))
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax_ab(b, P2, alpha, beta, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val > opt_val {
				opt_val = val
				opt_move = i
				alpha = max(alpha, opt_val)
			}
			if beta <= alpha {
				break
			}
		}
		return opt_val, opt_move
	} else if player == P2 {
		opt_val := MAX
		avail_moves := movesAvailable(b.heights, b.rows, b.cols)
		opt_move := rand.Intn(len(avail_moves))
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax_ab(b, P1, alpha, beta, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val < opt_val {
				opt_val = val
				opt_move = i
				beta = min(beta, opt_val)
			}
			if beta <= alpha {
				break
			}
		}
		return opt_val, opt_move
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func hello_seq() {
	prev_moves = make([]move, 0)
	fmt.Printf("Hello From Sequential Minimax\n")
	board := getBitBoard(6, 7)
	player := 1
	//fmt.Printf("Player 1: %d, Player 2: %d\n Player: %d, Player: %d\n", P1, P2, player, player^3)
	st = time.Now()
	game_res, player_res := seq_minimax_ab(board, player, MIN, MAX, 10)
	//game_res, player_res := seq_minimax(board, player, 10)
	fmt.Printf("The game result %d for player %d\n", game_res, player_res)
	board.printBoard()

}

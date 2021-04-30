package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var count int64 = 0
var per_10_mil time.Duration
var st time.Time

func seq_minimax(b *BitBoard, player int, depth int) (int, int) {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 || player_res != -1 {
		count += 1
		return game_res, player_res
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
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 || player_res != -1 {
		count += 1
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
				metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
				return opt_val, opt_move
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
				metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
				return opt_val, opt_move
			}
		}
		return opt_val, opt_move
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func seq(impl int, depth int) {
	board := getBitBoard(6, 7)
	game_res, player_res, player := 0, 0, 1
	st = time.Now()
	if impl == SEQ {
		fmt.Printf("Sequential Ran\n")
		game_res, player_res = seq_minimax(board, player, depth)
	} else if impl == SEQ_AB {
		fmt.Printf("Sequential_AB Ran\n")
		game_res, player_res = seq_minimax_ab(board, player, MIN, MAX, depth)
	}
	fmt.Printf("The game result %d for player %d. %d boards explored. %d nodes pruned\n", game_res, player_res, count, metrics.nodesPruned)
	board.printBoard()
}

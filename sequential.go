package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var moves []int
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
		opt_move := avail_moves[rand.Intn(len(avail_moves))]
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax(b, player^3, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val > opt_val {
				opt_val = val
				opt_move = avail_moves[i]
			}
		}
		return opt_val, opt_move
	} else if player == P2 {
		opt_val := MAX
		opt_move := avail_moves[rand.Intn(len(avail_moves))]
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax(b, player^3, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val < opt_val {
				opt_val = val
				opt_move = avail_moves[i]
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
		opt_move := avail_moves[rand.Intn(len(avail_moves))]
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax_ab(b, player^3, alpha, beta, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val > opt_val {
				opt_val = val
				opt_move = avail_moves[i]
			}
			alpha = max(alpha, opt_val)
			if beta <= alpha {
				metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
				return opt_val, opt_move
			}
		}
		moves = append(moves, opt_move)
		return opt_val, opt_move
	} else if player == P2 {
		opt_val := MAX
		avail_moves := movesAvailable(b.heights, b.rows, b.cols)
		opt_move := avail_moves[rand.Intn(len(avail_moves))]
		for i := range avail_moves {
			b.modBoard(avail_moves[i], player, 1)
			val, _ := seq_minimax_ab(b, player^3, alpha, beta, depth-1)
			b.modBoard(avail_moves[i], player, -1)
			if val < opt_val {
				opt_val = val
				opt_move = avail_moves[i]
			}
			beta = min(beta, opt_val)
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
	st = time.Now()
	game_res, player_res, move, player := 0, 0, 0, 1
	moves := make([]int, 0)
	if impl == SEQ {
		fmt.Printf("Sequential Ran\n")
		g1, g2 := board.gameState(MAX, player)
		for g1 == -1 && g2 == -1 {
			game_res, move = seq_minimax(board, player, depth)
			//fmt.Printf("Move %d is placing in column %d by player %d with value %d\n", moves_count, move, player, game_res)
			moves = append(moves, move)
			board.modBoard(move, player, 1)
			player ^= 3
			moves_count += 1
			g1, g2 = board.gameState(MAX, player)
		}
		game_res, player_res = board.gameState(MAX, player)
	} else if impl == SEQ_AB {
		fmt.Printf("Sequential_AB Ran\n")
		g1, g2 := board.gameState(MAX, player)
		for g1 == -1 && g2 == -1 {
			game_res, move = seq_minimax_ab(board, player, MIN, MAX, depth)
			//fmt.Printf("Move %d is placing in column %d by player %d\n", moves_count, move, player)
			moves = append(moves, move)
			board.modBoard(move, player, 1)
			player ^= 3
			moves_count += 1
			g1, g2 = board.gameState(MAX, player)
		}
		game_res, player_res = board.gameState(MAX, player)
	}
	fmt.Printf("Total Moves Required: %d\n", len(moves))
	fmt.Println(moves)
	fmt.Printf("The game result %d for player %d. %d boards explored. %d nodes pruned\n", game_res, player_res, count, metrics.nodesPruned)
	board.printBoard()
}

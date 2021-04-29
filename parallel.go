package main

import (
	"fmt"
	"math/rand"
)

const (
	PERCENT_SEQ      = 0.0
	PERCENT_PARALLEL = 1 - PERCENT_SEQ
)

func parallel_minimax(b *BitBoard, player int, depth int, pdepth int, ret chan move) {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 && player_res != -1 {
		ret <- move{game_res, player_res}
	}
	avail_moves := movesAvailable(b.heights, b.rows, b.cols)
	if player == P1 {
		if pdepth == 0 {
			//I think this can be replaced with seq_minimax
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
			ret <- move{opt_val, opt_move}
		} else {
			//launch goroutines and compute in parallel
			result := make(chan move)
			opt_val := MIN
			opt_move := rand.Intn(len(avail_moves))
			for i := range avail_moves {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax(nb, player^3, depth-1, pdepth-1, result)
			}
			for i := 0; i < len(avail_moves); i++ {
				cur_res := <-result
				if cur_res.val > opt_val {
					opt_val = cur_res.val
					opt_move = cur_res.col
				}
			}
			ret <- move{opt_val, opt_move}
		}
	} else if player == P2 {
		if pdepth == 0 {
			//I think this can be replaced with seq_minimax
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
			ret <- move{opt_val, opt_move}
		} else {
			//launch goroutines and compute in parallel
			result := make(chan move)
			opt_val := MAX
			opt_move := rand.Intn(len(avail_moves))
			for i := range avail_moves {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax(nb, player^3, depth-1, pdepth-1, result)
			}
			//not sure if this is correct or not -> maybe for cur_res := range result {...}
			for i := 0; i < len(avail_moves); i++ {
				cur_res := <-result
				if cur_res.val < opt_val {
					opt_val = cur_res.val
					opt_move = cur_res.col
				}
			}
			ret <- move{opt_val, opt_move}
		}
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func parallel_minimax_ab(b *BitBoard, player int, depth int, pdepth int, alpha int, beta int, ret chan move) {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 && player_res != -1 {
		ret <- move{game_res, player_res}
	}
	avail_moves := movesAvailable(b.heights, b.rows, b.cols)
	if player == P1 {
		if pdepth == 0 {
			opt_val := MIN
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
					ret <- move{opt_val, opt_move}
					return
				}
			}
			ret <- move{opt_val, opt_move}
		} else {
			opt_val := MIN
			opt_move := rand.Intn(len(avail_moves))
			result := make(chan move)
			//compute percentage of tree sequentially to utilize alpha-beta pruning (principle variation search)
			for i := 0; i < int(PERCENT_SEQ*len(avail_moves)); i++ {
				b.modBoard(avail_moves[i], player, 1)
				parallel_minimax_ab(b, player^3, depth-1, pdepth-1, alpha, beta, result)
				b.modBoard(avail_moves[i], player, -1)
			}
			for i := int(PERCENT_SEQ * len(avail_moves)); i < len(avail_moves); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, result)
			}
			for i := 0; i < len(avail_moves); i++ {
				cur_result := <-result
				if cur_result.val > opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
					alpha = max(alpha, opt_val)
				}
				if beta <= alpha {
					ret <- move{opt_val, opt_move}
					return
				}
			}
		}
	} else if player == P2 {
		if pdepth == 0 {
			opt_val := MAX
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
					ret <- move{opt_val, opt_move}
					return
				}
			}
			ret <- move{opt_val, opt_move}
		} else {
			opt_val := MAX
			opt_move := rand.Intn(len(avail_moves))
			result := make(chan move)
			//compute percentage of tree sequentially to utilize alpha-beta pruning (principle variation search)
			//we don't need to make a copy of the board because the moves are done sequentially
			for i := 0; i < int(PERCENT_SEQ*len(avail_moves)); i++ {
				b.modBoard(avail_moves[i], player, 1)
				parallel_minimax_ab(b, player^3, depth-1, pdepth-1, alpha, beta, result)
				b.modBoard(avail_moves[i], player, -1)
			}
			for i := int(PERCENT_SEQ * len(avail_moves)); i < len(avail_moves); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, result)
			}
			//TODO: This needs some way of knowing whether a result is coming from a "break" due to alpha/beta or not
			for i := 0; i < len(avail_moves); i++ {
				cur_result := <-result
				if cur_result.val < opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
					beta = min(beta, opt_val)
				}
				if beta <= alpha {
					ret <- move{opt_val, opt_move}
					return
				}
			}
		}
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func test_parallel() {
	fmt.Printf("Hello From Parallel\n")
}

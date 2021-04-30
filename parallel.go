package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

var PERCENT_SEQ float64 = 0.0
var PERCENT_PARALLEL float64 = 1 - PERCENT_SEQ

type move struct {
	val int
	col int
}

func parallel_minimax(b *BitBoard, player int, depth int, pdepth int, ret chan move) {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 || player_res != -1 {
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

//TODO:
func parallel_minimax_ab(b *BitBoard, player int, depth int, pdepth int, alpha int, beta int, ret chan move, doReturn bool) move {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 || player_res != -1 {
		if !doReturn {
			ret <- move{game_res, player_res}
		}
		return move{game_res, player_res}
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
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, opt_move}
					}
					return move{opt_val, opt_move}
				}
			}
			if !doReturn {
				ret <- move{opt_val, opt_move}
			}
			return move{opt_val, opt_move}
		} else {
			opt_val := MIN
			opt_move := rand.Intn(len(avail_moves))
			result := make(chan move)
			//compute percentage of tree sequentially to utilize alpha-beta pruning (principle variation search)
			for i := 0; i < int(PERCENT_SEQ*float64(len(avail_moves))); i++ {
				b.modBoard(avail_moves[i], player, 1)
				cur_result := parallel_minimax_ab(b, player^3, depth-1, pdepth-1, alpha, beta, result, true)
				b.modBoard(avail_moves[i], player, -1)
				if cur_result.val > opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
					alpha = max(alpha, opt_val)
				}
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, opt_move}
					}
					return move{opt_val, opt_move}
				}
			}
			fmt.Printf("Parallel section entered\n")
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, result, false)
			}
			//Gather step for minimax
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				cur_result := <-result
				if cur_result.val > opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
					alpha = max(alpha, opt_val)
				}
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, opt_move}
					}
					return move{opt_val, opt_move}
				}
			}
			if !doReturn {
				ret <- move{opt_val, opt_move}
			}
			return move{opt_val, opt_move}
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
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, opt_move}
					}
					return move{opt_val, opt_move}
				}
			}
			if !doReturn {
				ret <- move{opt_val, opt_move}
			}
			return move{opt_val, opt_move}
		} else {
			opt_val := MAX
			opt_move := rand.Intn(len(avail_moves))
			result := make(chan move)
			//compute percentage of tree sequentially to utilize alpha-beta pruning (principle variation search)
			//we don't need to make a copy of the board because the moves are done sequentially
			for i := 0; i < int(PERCENT_SEQ*float64(len(avail_moves))); i++ {
				b.modBoard(avail_moves[i], player, 1)
				cur_result := parallel_minimax_ab(b, player^3, depth-1, pdepth-1, alpha, beta, result, true)
				b.modBoard(avail_moves[i], player, -1)
				if cur_result.val < opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
					beta = min(beta, opt_val)
				}
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int(math.Pow(7, float64(depth)-1)/6)
					if !doReturn {
						ret <- move{opt_val, opt_move}
					}
					return move{opt_val, opt_move}
				}
			}
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, result, false)
			}
			//TODO: This needs some way of knowing whether a result is coming from a "break" due to alpha/beta or not
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				cur_result := <-result
				if cur_result.val < opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
					beta = min(beta, opt_val)
				}
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, opt_move}
					}
					return move{opt_val, opt_move}
				}
			}
			if !doReturn {
				ret <- move{opt_val, opt_move}
			}
			return move{opt_val, opt_move}
		}
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func parallel(impl int, depth int, pdepth int, percent_ab float64) {
	board := getBitBoard(6, 7)
	game_res, player_res, player := 0, 0, 1
	st = time.Now()
	if impl == PARALLEL {
		fmt.Printf("Parallel Ran\n")
		ret := make(chan move)
		go parallel_minimax(board, player, depth, pdepth, ret)
		result := <-ret
		game_res = result.val
		player_res = result.col
	} else if impl == PARALLEL_AB {
		fmt.Printf("Parallel_AB Ran\n")
		ret := make(chan move)
		PERCENT_SEQ = percent_ab
		PERCENT_PARALLEL = 1 - PERCENT_SEQ
		go parallel_minimax_ab(board, player, depth, pdepth, MIN, MAX, ret, false)
		result := <-ret
		close(ret)
		game_res = result.val
		player_res = result.col
	}
	fmt.Printf("The game result %d for player %d. %d boards explored. %d nodes pruned\n", game_res, player_res, count, metrics.nodesPruned)
	board.printBoard()
}
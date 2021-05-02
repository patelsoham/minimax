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
	player int
	opt_move int
	col int
}

func parallel_minimax(b *BitBoard, player int, depth int, pdepth int, col int, ret chan move) {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 || player_res != -1 {
		count += 1
		ret <- move{game_res, player_res, -1, col}
		return
	}
	avail_moves := movesAvailable(b.heights, b.rows, b.cols)
	if player == P1 {
		if pdepth == 0 {
			opt_val, opt_move := seq_minimax(b, player, depth)
			ret <- move{opt_val, player, opt_move, col}
			return
		} else {
			//launch goroutines and compute in parallel
			result := make(chan move)
			opt_val := MIN
			opt_move := avail_moves[rand.Intn(len(avail_moves))]
			for i := range avail_moves {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax(nb, player^3, depth-1, pdepth-1, avail_moves[i], result)
			}
			for i := 0; i < len(avail_moves); i++ {
				cur_res := <-result
				if cur_res.val > opt_val {
					opt_val = cur_res.val
					opt_move = cur_res.col
				} else if cur_res.val == opt_val {
					opt_move = min(opt_move, cur_res.col)
				}
			}
			ret <- move{opt_val, player, opt_move, col}
			return
		}
	} else if player == P2 {
		if pdepth == 0 {
			opt_val, opt_move := seq_minimax(b, player, depth)
			ret <- move{opt_val, player, opt_move, col}
			return
		} else {
			//launch goroutines and compute in parallel
			result := make(chan move)
			opt_val := MAX
			opt_move := avail_moves[rand.Intn(len(avail_moves))]
			for i := range avail_moves {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax(nb, player^3, depth-1, pdepth-1, avail_moves[i], result)
			}
			//not sure if this is correct or not -> maybe for cur_res := range result {...}
			for i := 0; i < len(avail_moves); i++ {
				cur_res := <-result
				if cur_res.val < opt_val {
					opt_val = cur_res.val
					opt_move = cur_res.col
				} else if cur_res.val == opt_val {
					opt_move = min(opt_move, cur_res.col)
				}
			}
			ret <- move{opt_val, player, opt_move, col}
			return
		}
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

//TODO:
func parallel_minimax_ab(b *BitBoard, player int, depth int, pdepth int, alpha int, beta int, col int, ret chan move, doReturn bool) move {
	game_res, player_res := b.gameState(depth, player)
	if game_res != -1 || player_res != -1 {
		count += 1
		if !doReturn {
			ret <- move{game_res, player_res, -1, col}
		}
		return move{game_res, player_res, -1, col}
	}
	avail_moves := movesAvailable(b.heights, b.rows, b.cols)
	if player == P1 {
		if pdepth == 0 {
			opt_val := MIN
			opt_move := avail_moves[rand.Intn(len(avail_moves))]
			for i := range avail_moves {
				b.modBoard(avail_moves[i], player, 1)
				val, _ := seq_minimax_ab(b, P2, alpha, beta, depth-1)
				b.modBoard(avail_moves[i], player, -1)
				if val > opt_val {
					opt_val = val
					opt_move = avail_moves[i]
				} 
				alpha = max(alpha, opt_val)
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, player, opt_move, col}
					}
					return move{opt_val, player, opt_move, col}
				}
			}
			if !doReturn {
				ret <- move{opt_val, player, opt_move, col}
			}
			return move{opt_val, player, opt_move, col}
		} else {
			opt_val := MIN
			opt_move := avail_moves[rand.Intn(len(avail_moves))]
			result := make(chan move)
			//compute percentage of tree sequentially to utilize alpha-beta pruning (principle variation search)
			for i := 0; i < int(PERCENT_SEQ*float64(len(avail_moves))); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				cur_result := parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, avail_moves[i], result, true)
				if cur_result.val > opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
				} else if cur_result.val == opt_val {
					opt_move = min(opt_move, cur_result.col)
				}
				alpha = max(alpha, opt_val)
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, player, opt_move, col}
					}
					return move{opt_val, player, opt_move, col}
				}
			}
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, avail_moves[i], result, false)
			}
			//Gather step for minimax
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				cur_result := <-result
				if cur_result.val > opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
				} else if cur_result.val == opt_val {
					opt_move = min(opt_move, cur_result.col)
				}
				alpha = max(alpha, opt_val)
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, player, opt_move, col}
					}
					return move{opt_val, player, opt_move, col}
				}
			}
			if !doReturn {
				ret <- move{opt_val, player, opt_move, col}
			}
			return move{opt_val, player, opt_move, col}
		}
	} else if player == P2 {
		if pdepth == 0 {
			opt_val := MAX
			opt_move := avail_moves[rand.Intn(len(avail_moves))]
			for i := range avail_moves {
				b.modBoard(avail_moves[i], player, 1)
				val, _ := seq_minimax_ab(b, P1, alpha, beta, depth-1)
				b.modBoard(avail_moves[i], player, -1)
				if val < opt_val {
					opt_val = val
					opt_move = avail_moves[i]
				}
				beta = min(beta, opt_val)
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, player, opt_move, col}
					}
					return move{opt_val, player, opt_move, col}
				}
			}
			if !doReturn {
				ret <- move{opt_val, player, opt_move, col}
			}
			return move{opt_val, player, opt_move, col}
		} else {
			opt_val := MAX
			opt_move := avail_moves[rand.Intn(len(avail_moves))]
			result := make(chan move)
			//compute percentage of tree sequentially to utilize alpha-beta pruning (principle variation search)
			//we don't need to make a copy of the board because the moves are done sequentially
			for i := 0; i < int(PERCENT_SEQ*float64(len(avail_moves))); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				cur_result := parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, avail_moves[i], result, true)
				if cur_result.val < opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
				} else if cur_result.val == opt_val {
					opt_move = min(opt_move, cur_result.col)
				}
				beta = min(beta, opt_val)
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int(math.Pow(7, float64(depth)-1)/6)
					if !doReturn {
						ret <- move{opt_val, player, opt_move, col}
					}
					return move{opt_val, player, opt_move, col}
				}
			}
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				nb := b.copyBoard()
				nb.modBoard(avail_moves[i], player, 1)
				go parallel_minimax_ab(nb, player^3, depth-1, pdepth-1, alpha, beta, avail_moves[i], result, false)
			}
			//TODO: This needs some way of knowing whether a result is coming from a "break" due to alpha/beta or not
			for i := int(PERCENT_SEQ * float64(len(avail_moves))); i < len(avail_moves); i++ {
				cur_result := <-result
				if cur_result.val < opt_val {
					opt_val = cur_result.val
					opt_move = cur_result.col
				} else if cur_result.val == opt_val {
					opt_move = min(opt_move, cur_result.col)
				}
				beta = min(beta, opt_val)
				if beta <= alpha {
					metrics.nodesPruned += (len(avail_moves) - i - 1) * int((math.Pow(7, float64(depth))-1)/6)
					if !doReturn {
						ret <- move{opt_val, player, opt_move, col}
					}
					return move{opt_val, player, opt_move, col}
				}
			}
			if !doReturn {
				ret <- move{opt_val, player, opt_move, col}
			}
			return move{opt_val, player, opt_move, col}
		}
	} else {
		fmt.Printf("Invalid player number %d\n", player)
		panic(player)
	}
}

func parallel(impl int, depth int, pdepth int, percent_ab float64) {
	board := getBitBoard(6, 7)
	game_res, player_res, cur_move, player := 0, 0, 0, 1
	st = time.Now()
	moves := make([]int, 0)
	if impl == PARALLEL {
		fmt.Printf("Parallel Ran\n")
		ret := make(chan move)
		g1, g2 := board.gameState(MAX, player)
		for g1 == -1 && g2 == -1 {
			go parallel_minimax(board, player, depth, pdepth, rand.Intn(7), ret)
			result := <-ret
			game_res, cur_move = result.val, result.opt_move
			//fmt.Printf("Move %d is placing in column %d by player %d with value %d\n", moves_count, cur_move, player, game_res)
			moves = append(moves, cur_move)
			board.modBoard(cur_move, player, 1)
			player ^= 3
			moves_count += 1
			g1, g2 = board.gameState(MAX, player)
		}
		game_res, player_res = board.gameState(MAX, player)
	} else if impl == PARALLEL_AB {
		fmt.Printf("Parallel_AB Ran\n")
		ret := make(chan move)
		PERCENT_SEQ = percent_ab
		PERCENT_PARALLEL = 1 - PERCENT_SEQ
		g1, g2 := board.gameState(MAX, player)
		for g1 == -1 && g2 == -1 {
			go parallel_minimax_ab(board, player, depth, pdepth, MIN, MAX, rand.Intn(7), ret, false)
			result := <-ret
			game_res, cur_move = result.val, result.opt_move
			//fmt.Printf("Move %d is placing in column %d by player %d\n", moves_count, cur_move, player)
			moves = append(moves, cur_move)
			board.modBoard(cur_move, player, 1)
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

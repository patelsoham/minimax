package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type benchmarks struct {
	endtoend    time.Duration
	gameState   time.Duration
}

var lock sync.Mutex

func incrementGameState(toAdd time.Duration) {
	return // comment this line to calculate the amount time spent calculating the game state for metrics 
	lock.Lock()
	metrics.gameState += toAdd
	lock.Unlock()
}

func getGameState() float64 {
	return metrics.gameState.Seconds()
}

var metrics benchmarks
var moves_count int = 0

const (
	SEQ         = iota
	SEQ_AB      = iota
	PARALLEL    = iota
	PARALLEL_AB = iota
	NUM_IMPL = iota
)

func main() {
	impl, depth, pdepth, ab_percent_sequential := 0, 0, 0, 0.0
	metrics := benchmarks{0, 0}
	flag.IntVar(&impl, "i", 0, "which implementation to run: 0 = sequential, 1 = sequential_ab, 2 = parallel, 3 = parallel_ab")
	flag.IntVar(&depth, "d", 5, "max depth of algorithms")
	flag.IntVar(&pdepth, "pd", -1, "max depth computed in parallel")
	flag.Float64Var(&ab_percent_sequential, "ab", 0.5, "percentage of the parallel AB solution done in sequential")
	flag.Parse()
	if pdepth == -1 || pdepth > depth {
		pdepth = depth
	}
	fmt.Printf("impl: %d, depth: %d, pdepth: %d, PERCENT_SEQ: %.5f\n", impl, depth, pdepth, ab_percent_sequential)
	st := time.Now()
	switch {
	case impl < 2:
		seq(impl, depth)
	case impl >= 2 && impl < NUM_IMPL:
		parallel(impl, depth, pdepth, ab_percent_sequential)
	default:
		fmt.Println("Default case entered")
		test_helpers()
	}
	metrics.endtoend += time.Now().Sub(st)
	// print metrics
	fmt.Printf("End to End Time: %.5f\nGame State Computation Time: %.10f\n", metrics.endtoend.Seconds(), getGameState())
}

package main

import (
	"flag"
	"fmt"
	"gonum.org/v1/gonum/stat"
	"os"
	"runtime/pprof"
	"time"
)

var randomSeed = flag.Int64("seed", time.Now().UnixNano(), "Random seed for the randomness source.")
var samples = flag.Int("samples", 1000, "Number of samples in the simulation.")
var prof = flag.String("prof", "", "filepath to write CPU profile to.")

func main() {
	flag.Parse()
	if *prof != "" {
		f, err := os.Create(*prof)
		if err != nil {
			fmt.Printf("failed to create file: %s", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	fmt.Printf("Seed: %d\n", *randomSeed)
	g, err := NewGameFromState(&GameStateInput{
		Camels: map[BoardPosition][]Color{
			0: {Blue, Green, Red, Yellow, Purple},
			5: {White, Black},
		},
	})
	fmt.Printf("Game state:\n%s\n", g)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	times := make([]float64, *samples, *samples)
	for i := range *samples {
		start := time.Now()
		g.ComputeLegRankingDistribution()
		end := time.Now()
		times[i] = float64(end.Sub(start).Nanoseconds())
	}
	mean, variance := stat.MeanVariance(times, nil)
	fmt.Println("Search time stats:")
	fmt.Printf("Mean: %5.2f ms\n", mean/1000000)
	fmt.Printf("Variance: %f ms squared \n", variance/float64(1000000*1000000))
}

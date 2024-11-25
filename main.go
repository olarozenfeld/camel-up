package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var randomSeed = flag.Int64("seed", time.Now().UnixNano(), "Random seed for the randomness source.")
var samples = flag.Int("samples", 1000, "Number of samples in the simulation.")

func main() {
	flag.Parse()
	fmt.Printf("Seed: %d\n", *randomSeed)
	r := rand.New(rand.NewSource(*randomSeed))
	g, err := NewGameFromState(&GameStateInput{
		Camels: map[BoardPosition][]Color{
			0: {Blue, Green, Red, Yellow, Purple},
			5: {White, Black},
		},
		DiePyramid: NewDiePyramid(r),
	})
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	start := time.Now()
	dist := g.SimulateLegRankingDistribution(*samples)
	end := time.Now()

	fmt.Printf("time: %s\n", end.Sub(start))
	fmt.Printf("Distribution: \n%s\n", dist)
}

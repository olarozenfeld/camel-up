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
			0: {Blue, Black, White},
			1: {Green},
			2: {Red},
			3: {Yellow},
			4: {Purple},
		},
		DiePyramid: NewDiePyramidWithDice(r, []Color{Purple}),
	})
	fmt.Printf("Game state:\n%s\n", g)
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

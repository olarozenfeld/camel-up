package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var randomSeed = flag.Int64("seed", time.Now().UnixNano(), "Random seed for the randomness source.")

func main() {
	flag.Parse()
	fmt.Printf("Seed: %d\n", *randomSeed)
	r := rand.New(rand.NewSource(*randomSeed))
	g, err := NewGameFromState(&GameStateInput{
		Camels: map[BoardPosition][]Color{
			1:  {Yellow, Green},
			3:  {Red, Blue},
			8:  {Purple},
			14: {White, Black},
		},
		Cheers: map[BoardPosition]string{
			2: "",
		},
		DiePyramid: NewDiePyramid(r),
	})
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", g)
	d := &RankingDistribution{}
	//	d.RecordRanking(&g.ranking)
	rs := [NumMovesPerLeg]Color{Yellow, Green, Red, Blue, Purple}
	for i := 0; i < 10000; i++ {
		r.Shuffle(NumMovesPerLeg, func(i, j int) {
			rs[i], rs[j] = rs[j], rs[i]
		})
		d.RecordRanking(&rs)
	}
	fmt.Printf("%s\n", d)
}

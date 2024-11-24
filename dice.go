package main

import (
	"fmt"
	"math/rand"
)

type RollValue int // 1,2,3

type DieRoll struct {
	Color Color
	Value RollValue
}

func (r *DieRoll) String() string {
	return colorPrinters[r.Color](" %d ", r.Value)
}

type DiePyramid struct {
	r        *rand.Rand
	numRolls int
	dice     []Color
}

var ErrOutOfDice = fmt.Errorf("out of dice")

// Creates a pyramid with all 6 dice available (Black stands for the grey die).
func NewDiePyramid(r *rand.Rand) *DiePyramid {
	return NewDiePyramidWithDice(r, []Color{Green, Yellow, Red, Blue, Purple, Black})
}

// Creates a die pyramid with only the specified N dice in it, prepared for rolling
// N-1 of them. Used for simulations/computations where some of the dice have
// been already rolled out.
// TODO: validate input ([2-6], unique colors, no White); copy input slice.
func NewDiePyramidWithDice(r *rand.Rand, dice []Color) *DiePyramid {
	result := &DiePyramid{r: r, dice: dice}
	result.Reset()
	return result
}

// Resets the pyramid to the starting dice.
func (p *DiePyramid) Reset() {
	// Shuffle the colors (Black stands for the grey die).
	p.r.Shuffle(len(p.dice), func(i, j int) {
		p.dice[i], p.dice[j] = p.dice[j], p.dice[i]
	})
	p.numRolls = 0
}

func (p *DiePyramid) RemainingRolls() int {
	return len(p.dice) - 1 - p.numRolls
}

func (p *DiePyramid) IsEmpty() bool {
	return p.RemainingRolls() == 0
}

func (p *DiePyramid) Roll() (DieRoll, error) {
	result := DieRoll{-1, -1}
	if p.numRolls == len(p.dice)-1 {
		return result, ErrOutOfDice
	}
	result.Color = p.dice[p.numRolls]
	if result.Color == Black { // Grey, can be Black/White
		result.Color = Color(p.r.Intn(2)) + Black
	}
	result.Value = RollValue(p.r.Intn(3) + 1)
	p.numRolls++
	return result, nil
}

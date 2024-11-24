package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"gonum.org/v1/gonum/stat"
)

const numSamples = 1000
const stDevCutoff = 30

type TestResultDistribution struct {
	colors []Color
	// Count results by: step X color X value
	result [5][7][3]int
}

func (d *TestResultDistribution) RecordRoll(step int, roll *DieRoll) {
	d.result[step][roll.Color][roll.Value-1]++
}

func (d *TestResultDistribution) Validate() error {
	// Everything for unused steps should be 0.
	for s := len(d.colors) - 1; s < 5; s++ {
		for c := Green; c <= White; c++ {
			for r := 0; r < 3; r++ {
				if d.result[s][c][r] != 0 {
					return fmt.Errorf("expected result[%d][%s][%d] to be 0, got %d", s, c, r, d.result[s][c][r])
				}
			}
		}
	}
	usedColors := make(map[Color]bool)
	for _, c := range d.colors {
		usedColors[c] = true
		if c == Black {
			usedColors[White] = true
		}
	}
	// Everything for unused colors should be 0.
	for c := Green; c <= White; c++ {
		if usedColors[c] {
			continue
		}
		for s := 0; s < 5; s++ {
			for r := 0; r < 3; r++ {
				if d.result[s][c][r] != 0 {
					return fmt.Errorf("expected result[%d][%s][%d] to be 0, got %d", s, c, r, d.result[s][c][r])
				}
			}
		}
	}
	// White and Black share a distribution; we expect the used
	// colors to have approximately equal number of samples, with
	// Black and White having around half of that amount.
	var normalResults []float64
	var crazyResults []float64
	for c := range usedColors {
		for s := 0; s < len(d.colors)-1; s++ {
			for r := 0; r < 3; r++ {
				sample := float64(d.result[s][c][r])
				if c.IsCrazy() {
					crazyResults = append(crazyResults, sample)
				} else {
					normalResults = append(normalResults, sample)
				}
			}
		}
	}
	normMean, normStd := stat.MeanStdDev(normalResults, nil)
	crazyMean, crazyStd := stat.MeanStdDev(crazyResults, nil)
	if normStd > stDevCutoff {
		return fmt.Errorf("expected standard deviation of normal colors to be less than %d, got %f (samples: %v)", stDevCutoff, normStd, d.result)
	}
	if crazyResults != nil {
		if crazyStd > stDevCutoff {
			return fmt.Errorf("expected standard deviation of crazy colors to be less than %d, got %f (samples: %v)", stDevCutoff, crazyStd, d.result)
		}
		if math.Abs(crazyMean*2-normMean) > stDevCutoff {
			return fmt.Errorf("expected crazy colors mean to be within stddev of half of normal colors mean, got samples: %v", d.result)
		}
	}
	return nil
}

func TestDiePyramidDistributions(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// For k=2 to 6 choose random k colors
	allColors := []Color{Green, Yellow, Red, Blue, Purple, Black}
	for k := 2; k <= 6; k++ {
		// Choose random k colors:
		r.Shuffle(len(allColors), func(i, j int) {
			allColors[i], allColors[j] = allColors[j], allColors[i]
		})
		p := NewDiePyramidWithDice(r, allColors[:k])
		d := &TestResultDistribution{colors: allColors[:k]}
		for i := 0; i < numSamples; i++ {
			for s := 0; s < k-1; s++ {
				roll, err := p.Roll()
				if err != nil {
					t.Errorf("Roll() failed: %v", err)
				}
				d.RecordRoll(s, &roll)
			}
			if !p.IsEmpty() {
				t.Error("expected IsEmpty() to be true")
			}
			p.Reset()
		}
		if err := d.Validate(); err != nil {
			t.Errorf("results failed validation: %v\n%v\n", err, d)
		}
	}
}

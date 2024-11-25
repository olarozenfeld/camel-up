package main

import (
	"fmt"
	"math"
	"strings"
)

// A ranking distribution is a result of either a simulation or a computed
// enumeration of all possible leg results. It counts all the possible rankings
// for all the racing camels.

type RankingDistribution struct {
	TotalRankings int
	// Rank x Color
	Rankings [NumRacingCamels][NumRacingCamels]int
}

func (d *RankingDistribution) RecordRanking(ranking *[NumRacingCamels]Color) {
	d.TotalRankings++
	for r, c := range ranking {
		d.Rankings[r][c]++
	}
}

func (d *RankingDistribution) String() string {
	var s strings.Builder
	fmt.Fprintf(&s, "Total rankings: %d\n", d.TotalRankings)
	digits := int(math.Log10(float64(d.TotalRankings))) + 1
	headerPattern := strings.Repeat(fmt.Sprintf("\t%%%ds", digits+9), 5)
	s.WriteString(colorPrinters[White](headerPattern+"\n", "Last", "4th", "3rd", "2nd", "First"))
	for c := Green; c < Black; c++ {
		fmt.Fprintf(&s, "%s\t", c)
		for r := Last; r <= First; r++ {
			samples := d.Rankings[c][r]
			percentage := float64(samples) * 100 / float64(d.TotalRankings)
			pattern := fmt.Sprintf("%%%dd (%%5.2f%%%%)\t", digits)
			fmt.Fprintf(&s, pattern, samples, percentage)
		}
		s.WriteString("\n")
	}
	return s.String()
}

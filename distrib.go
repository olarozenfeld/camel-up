package main

import (
	"fmt"
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
	s.WriteString(colorPrinters[White]("\tLast\t\t4th\t\t3rd\t\t2nd\t\tFirst\n"))
	for c := Green; c < Black; c++ {
		fmt.Fprintf(&s, "%s\t", c)
		for r := Last; r <= First; r++ {
			samples := d.Rankings[c][r]
			percentage := float64(samples) * 100 / float64(d.TotalRankings)
			fmt.Fprintf(&s, "%d (%.2f%%)\t", samples, percentage)
		}
		s.WriteString("\n")
	}
	return s.String()
}

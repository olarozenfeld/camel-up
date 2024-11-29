package main

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestNewGameFromStateFromInputSuccess(t *testing.T) {
	testCases := []struct {
		name        string
		input       *GameStateInput
		wantRanking [NumRacingCamels]Color
	}{
		{
			name: "normal leads",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Green, Red, Blue, Purple},
		},
		{
			name: "normal on crazy",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Blue},
					8:  {Red, Purple},
					14: {Black, White},
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Green, Blue, Red, Purple},
		},
		{
			name: "crazy on normal",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Green},
					3:  {Blue},
					8:  {Black, Red},
					14: {Yellow, Purple, White},
				},
			},
			wantRanking: [NumRacingCamels]Color{Green, Blue, Red, Yellow, Purple},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := NewGameFromState(tc.input)
			if err != nil {
				t.Error(err)
			}
			if board.ranking != tc.wantRanking {
				t.Errorf("want ranking %s, got %s", tc.wantRanking, board.ranking)
			}
		})
	}
}

func TestNewGameFromStateFromInputFailure(t *testing.T) {
	testCases := []struct {
		name      string
		input     *GameStateInput
		wantError string
	}{
		{
			name: "position too high",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					16: {White, Black},
				},
			},
			wantError: "invalid board position: 16",
		},
		{
			name: "position too low",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					-1: {Purple},
					15: {White, Black},
				},
			},
			wantError: "invalid board position: -1",
		},
		{
			name: "doubled camel",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue, Yellow},
					8:  {Purple},
					15: {White, Black},
				},
			},
			wantError: "camel appears twice in input",
		},
		{
			name: "missing camel",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red},
					8:  {Purple},
					15: {White, Black},
				},
			},
			wantError: "camel is not placed on the board",
		},
		{
			name: "non-empty cheer",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					15: {White, Black},
				},
				Cheers: map[BoardPosition]string{3: ""},
			},
			wantError: "invalid cheer position 3, not empty",
		},
		{
			name: "non-empty boo",
			input: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					15: {White, Black},
				},
				Boos: map[BoardPosition]string{3: ""},
			},
			wantError: "invalid boo position 3, not empty",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			board, err := NewGameFromState(tc.input)
			if err == nil {
				t.Errorf("wanted error with %s, got:\n%s", tc.wantError, board)
			} else if !strings.Contains(err.Error(), tc.wantError) {
				t.Errorf("wanted error with %s, got:\n%s", tc.wantError, err)
			}
		})
	}
}

// Compares camel tokens of different boards.
func (c *camel) equals(o *camel) bool {
	if c == nil {
		return o == nil
	}
	if o == nil {
		return false
	}
	return c.Color == o.Color && c.Position == o.Position && c.Next.shallowEquals(o.Next) && c.Prev.shallowEquals(o.Prev) && c.OtherCrazy.shallowEquals(o.OtherCrazy)
}

func (c *camel) shallowString() string {
	if c == nil {
		return "nil"
	}
	return c.Color.String()
}

func (c *camel) String() string {
	if c == nil {
		return "nil"
	}
	return fmt.Sprintf("{%s, %d, %s, %s, %s}", c.Color, c.Position, c.Prev.shallowString(), c.Next.shallowString(), c.OtherCrazy.shallowString())
}

func (s *boardSpace) String() string {
	if s == nil {
		return "nil"
	}
	return fmt.Sprintf("{ Cheer: %d Boo: %d StackBottom: %s StackTop: %s}\n", s.Cheer, s.Boo, s.StackBottom.shallowString(), s.StackTop.shallowString())
}

func (g *Game) fullString() string {
	var s strings.Builder
	for i := range g.boardSpaces {
		s.WriteString(g.boardSpaces[i].String())
	}
	return s.String()
}

// This is needed to avoid infinite loops.
func (c *camel) shallowEquals(o *camel) bool {
	if c == nil {
		return o == nil
	}
	if o == nil {
		return false
	}
	return c.Color == o.Color
}

// Compares board spaces of different boards.
func (s *boardSpace) equals(o *boardSpace) bool {
	if s == nil {
		return o == nil
	}
	if o == nil {
		return false
	}
	return s.Cheer == o.Cheer && s.Boo == o.Boo && s.StackTop.equals(o.StackTop) && s.StackBottom.equals(o.StackBottom)
}

// Compares boards of different games.
func (g *Game) equals(o *Game) bool {
	//	fmt.Printf(">>>>>>>>>>>> g.fullString= %s\n", g.fullString())
	//	fmt.Printf(">>>>>>>>>>>> o.fullString= %s\n", o.fullString())
	for i := range g.camelTokens {
		if !g.camelTokens[i].equals(&o.camelTokens[i]) {
			//			fmt.Printf(">>>>>>>>>>>> camelTokens %d\n", i)
			//			fmt.Printf(">>>>>>>>>>>> g.camelTokens[0]= %v\n", &g.camelTokens[i])
			//			fmt.Printf(">>>>>>>>>>>> o.camelTokens[0]= %v\n", &o.camelTokens[i])
			//			fmt.Printf(">>>>>>>>>>>> g.camelTokens[0].Prev= %v\n", g.camelTokens[i].Prev)
			//			fmt.Printf(">>>>>>>>>>>> o.camelTokens[0].Prev= %v\n", o.camelTokens[i].Prev)
			return false
		}
	}
	for i := range g.boardSpaces {
		if !g.boardSpaces[i].equals(&o.boardSpaces[i]) {
			//			fmt.Printf(">>>>>>>>>>>> boardSpaces %d\n", i)
			return false
		}
	}
	return true
}

func TestApplyCamelMove(t *testing.T) {
	testCases := []struct {
		startState  *GameStateInput
		dieRoll     *DieRoll
		wantState   *GameStateInput
		wantRanking [NumRacingCamels]Color
	}{
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
			},
			dieRoll: &DieRoll{Green, 2},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow},
					3:  {Red, Blue, Green},
					8:  {Purple},
					14: {White, Black},
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Red, Blue, Green, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
			},
			dieRoll: &DieRoll{Yellow, 2},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					3:  {Red, Blue, Yellow, Green},
					8:  {Purple},
					14: {White, Black},
				},
			},
			wantRanking: [NumRacingCamels]Color{Red, Blue, Yellow, Green, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
				Cheers: map[BoardPosition]string{
					2: "",
				},
			},
			dieRoll: &DieRoll{Green, 1},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow},
					3:  {Red, Blue, Green},
					8:  {Purple},
					14: {White, Black},
				},
				Cheers: map[BoardPosition]string{
					2: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Red, Blue, Green, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
				Boos: map[BoardPosition]string{
					4: "",
				},
			},
			dieRoll: &DieRoll{Green, 3},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow},
					3:  {Green, Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
				Boos: map[BoardPosition]string{
					4: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Green, Red, Blue, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
				Boos: map[BoardPosition]string{
					4: "",
				},
			},
			dieRoll: &DieRoll{Yellow, 3},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					3:  {Yellow, Green, Red, Blue},
					8:  {Purple},
					14: {White, Black},
				},
				Boos: map[BoardPosition]string{
					4: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Green, Red, Blue, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					9:  {Red, Blue},
					10: {White, Purple},
					14: {Yellow, Green, Black},
				},
				Cheers: map[BoardPosition]string{
					15: "",
				},
			},
			dieRoll: &DieRoll{Yellow, 1},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					9:  {Red, Blue},
					10: {White, Purple},
					0:  {Yellow, Green, Black},
				},
				Cheers: map[BoardPosition]string{
					15: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Red, Blue, Purple, Yellow, Green},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0:  {Red, Blue},
					10: {White, Purple},
					14: {Yellow, Black, Green},
				},
			},
			dieRoll: &DieRoll{Yellow, 2},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0:  {Red, Blue, Yellow, Black, Green},
					10: {White, Purple},
				},
			},
			wantRanking: [NumRacingCamels]Color{Red, Blue, Purple, Yellow, Green},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {Red, Blue, White, Purple},
					2: {Green, Black, Yellow},
				},
				Boos: map[BoardPosition]string{
					1: "",
				},
			},
			dieRoll: &DieRoll{Blue, 1},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {Blue, White, Purple, Red},
					2: {Green, Black, Yellow},
				},
				Boos: map[BoardPosition]string{
					1: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Blue, Purple, Red, Green, Yellow},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					12: {Purple},
					14: {White, Black},
				},
			},
			dieRoll: &DieRoll{White, 2},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					3:  {Red, Blue},
					12: {Purple, Black},
					14: {White},
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Green, Red, Blue, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					12: {Red},
					13: {Purple},
					14: {White, Black, Blue},
				},
			},
			dieRoll: &DieRoll{White, 2},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green},
					12: {Red, Black, Blue},
					13: {Purple},
					14: {White},
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Green, Red, Blue, Purple},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {Red, Blue, White, Purple, Black, Yellow},
					2: {Green},
				},
				Boos: map[BoardPosition]string{
					15: "",
				},
			},
			dieRoll: &DieRoll{White, 1},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {White, Purple, Black, Yellow, Red, Blue},
					2: {Green},
				},
				Boos: map[BoardPosition]string{
					15: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Purple, Yellow, Red, Blue, Green},
		},
		{
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {Red, Blue, White, Purple, Black, Yellow},
					2: {Green},
				},
				Boos: map[BoardPosition]string{
					15: "",
				},
			},
			dieRoll: &DieRoll{Black, 1},
			wantState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {Black, Yellow, Red, Blue, White, Purple},
					2: {Green},
				},
				Boos: map[BoardPosition]string{
					15: "",
				},
			},
			wantRanking: [NumRacingCamels]Color{Yellow, Red, Blue, Purple, Green},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
			start, err := NewGameFromState(tc.startState)
			if err != nil {
				t.Fatal(err)
			}
			want, err := NewGameFromState(tc.wantState)
			if err != nil {
				t.Fatal(err)
			}
			start.applyCamelMove(tc.dieRoll)
			if !start.equals(want) {
				t.Errorf("want board state:\n%s\nGot board state:\n%s\n", want, start)
			}
			if start.ranking != tc.wantRanking {
				t.Errorf("want ranking %s, got %s", tc.wantRanking, start.ranking)
			}
			cp, err := NewGameFromState(tc.startState)
			if err != nil {
				t.Fatal(err)
			}
			start.undoLastCamelMove()
			if !start.equals(cp) {
				t.Errorf("want board state after undo:\n%s\nGot board state:\n%s\n", cp, start)
			}
		})
	}
}

// The test cases and their correct distributions are curtesy of https://github.com/nishchalchandna/camel_up/blob/main/lib/search_test.go
// Also verified by SimulateLegRankingDistribution.
func TestComputeLegRankingDistribution(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	testCases := []struct {
		desc             string
		startState       *GameStateInput
		wantDistribution *RankingDistribution
	}{
		{
			desc: "base case: all camels, no dice to roll",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0: {Blue, Black, White},
					1: {Green},
					2: {Red},
					3: {Yellow},
					4: {Purple},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Purple}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 1,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 1, 0, 0, 0}, // Green
					{0, 0, 0, 1, 0}, // Yellow
					{0, 0, 1, 0, 0}, // Red
					{1, 0, 0, 0, 0}, // Blue
					{0, 0, 0, 0, 1}, // Purple
				},
			},
		},
		{
			desc: "one die left: blue, green",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Red, Yellow, Purple},
					8:  {Green},
					12: {Blue},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Green, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 12,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 0, 0, 12, 0}, // Green
					{0, 12, 0, 0, 0}, // Yellow
					{12, 0, 0, 0, 0}, // Red
					{0, 0, 0, 0, 12}, // Blue
					{0, 0, 12, 0, 0}, // Purple
				},
			},
		},
		{
			desc: "one die left: blue, green stacked",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Red, Yellow, Purple},
					8:  {Green, Blue},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Green, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 12,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 0, 0, 12, 0}, // Green
					{0, 12, 0, 0, 0}, // Yellow
					{12, 0, 0, 0, 0}, // Red
					{0, 0, 0, 0, 12}, // Blue
					{0, 0, 12, 0, 0}, // Purple
				},
			},
		},
		{
			desc: "one die left: blue, green split",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Red, Yellow, Purple},
					8:  {Green},
					9:  {Blue},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Green, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 12,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 0, 0, 6, 6},  // Green
					{0, 12, 0, 0, 0}, // Yellow
					{12, 0, 0, 0, 0}, // Red
					{0, 0, 0, 6, 6},  // Blue
					{0, 0, 12, 0, 0}, // Purple
				},
			},
		},
		{
			desc: "one die left: blue, green, red",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Purple},
					8:  {Green},
					9:  {Blue},
					10: {Red},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Green, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 12,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 0, 6, 2, 4},  // Green
					{12, 0, 0, 0, 0}, // Yellow
					{0, 0, 0, 10, 2}, // Red
					{0, 0, 6, 0, 6},  // Blue
					{0, 12, 0, 0, 0}, // Purple
				},
			},
		},
		{
			desc: "one die left: blue, green, red stacked",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Purple},
					8:  {Green},
					9:  {Blue, Red},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Green, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 12,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 0, 6, 0, 6},  // Green
					{12, 0, 0, 0, 0}, // Yellow
					{0, 0, 0, 6, 6},  // Red
					{0, 0, 6, 6, 0},  // Blue
					{0, 12, 0, 0, 0}, // Purple
				},
			},
		},
		{
			desc: "all unstacked, 2nd and 4th dice remain",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					0:  {Yellow},
					1:  {Purple},
					2:  {Green},
					3:  {Blue},
					4:  {Red},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Purple, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 12,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 6, 6, 0, 0},  // Green
					{12, 0, 0, 0, 0}, // Yellow
					{0, 0, 0, 8, 4},  // Red
					{0, 0, 4, 2, 6},  // Blue
					{0, 6, 2, 2, 2},  // Purple
				},
			},
		},
		{
			desc: "all stacked, top three dice remain",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green, Red, Blue, Purple},
					13: {Black, White},
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Purple, Red, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 216,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 216, 0, 0, 0},   // Green
					{216, 0, 0, 0, 0},   // Yellow
					{0, 0, 168, 24, 24}, // Red
					{0, 0, 24, 144, 48}, // Blue
					{0, 0, 24, 48, 144}, // Purple
				},
			},
		},
		{
			desc: "all stacked, top three dice remain with Boo in front",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green, Red, Blue, Purple},
					13: {Black, White},
				},
				Boos: map[BoardPosition]string{
					2: "",
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Purple, Red, Blue}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 216,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 120, 40, 32, 24},  // Green
					{120, 40, 32, 24, 0},  // Yellow
					{40, 0, 96, 44, 36},   // Red
					{32, 32, 12, 92, 48},  // Blue
					{24, 24, 36, 24, 108}, // Purple
				},
			},
		},
		{
			desc: "two dice left - camels far apart",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1:  {Yellow, Green, Red},
					5:  {Blue},
					9:  {Purple},
					13: {Black, White},
				},
				Cheers: map[BoardPosition]string{
					8: "",
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Blue, Purple, Red}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 216,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 216, 0, 0, 0},  // Green
					{216, 0, 0, 0, 0},  // Yellow
					{0, 0, 216, 0, 0},  // Red
					{0, 0, 0, 180, 36}, // Blue
					{0, 0, 0, 36, 180}, // Purple
				},
			},
		},
		{
			desc: "two dice left - camels far apart - crazy camel - simple",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1: {Yellow, Green, White, Red},
					2: {Black},
					5: {Blue},
					9: {Purple},
				},
				Cheers: map[BoardPosition]string{
					8: "",
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Blue, Purple, Black}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 216,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 72, 144, 0, 0}, // Green
					{72, 144, 0, 0, 0}, // Yellow
					{144, 0, 72, 0, 0}, // Red
					{0, 0, 0, 188, 28}, // Blue
					{0, 0, 0, 28, 188}, // Purple
				},
			},
		},
		{
			desc: "two dice left - camels far apart - crazy camel - complicated",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					4:  {Yellow, Green, White, Red},
					8:  {Black, Blue},
					12: {Purple},
				},
				Cheers: map[BoardPosition]string{
					11: "",
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Blue, Purple, Black}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 216,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 126, 90, 0, 0}, // Green
					{126, 90, 0, 0, 0}, // Yellow
					{90, 0, 126, 0, 0}, // Red
					{0, 0, 0, 186, 30}, // Blue
					{0, 0, 0, 30, 186}, // Purple
				},
			},
		},
		{
			desc: "two dice left - camels far apart - crazy camel - complicated, game ends",
			startState: &GameStateInput{
				Camels: map[BoardPosition][]Color{
					1: {Yellow, Green, White, Red},
					5: {Black, Blue},
					9: {Purple},
				},
				Cheers: map[BoardPosition]string{
					8: "",
				},
				DiePyramid: NewDiePyramidWithDice(r, []Color{Blue, Purple, Black}),
			},
			wantDistribution: &RankingDistribution{
				TotalRankings: 216,
				Rankings: [NumRacingCamels][NumRacingCamels]int{
					{0, 126, 90, 0, 0}, // Green
					{126, 90, 0, 0, 0}, // Yellow
					{90, 0, 126, 0, 0}, // Red
					{0, 0, 0, 190, 26}, // Blue
					{0, 0, 0, 26, 190}, // Purple
				},
			},
		},
	}
	t.Parallel()
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			g, err := NewGameFromState(tc.startState)
			if err != nil {
				t.Fatal(err)
			}
			computation := g.ComputeLegRankingDistribution()
			if *tc.wantDistribution != *computation {
				t.Errorf("want distribution:\n%s, got:\n%s", tc.wantDistribution, computation)
			}
		})
	}
}

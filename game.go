package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type BoardPosition int

const (
	StartPosition  BoardPosition = 0
	FinishPosition BoardPosition = 15
	BoardSize                    = 16
)

func (p BoardPosition) Add(k int) BoardPosition {
	return BoardPosition((BoardSize + int(p) + k) % BoardSize)
}

// A camel token on the board.
type camel struct {
	Color      Color
	Position   BoardPosition
	Next       *camel // The camel on top of this one.
	Prev       *camel // The camel below this one.
	OtherCrazy *camel // White<->Black
}

func (c *camel) IsCrazy() bool {
	return c.Color.IsCrazy()
}

type Player int

const NoPlayer Player = -1

type boardSpace struct {
	Cheer       Player
	Boo         Player
	StackBottom *camel
	StackTop    *camel
}

func (s *boardSpace) HasCheer() bool {
	return s.Cheer != NoPlayer
}

func (s *boardSpace) HasBoo() bool {
	return s.Boo != NoPlayer
}

type Rank int

const (
	Last            Rank = 0
	First           Rank = 4
	NumCamels            = 7
	NumRacingCamels      = 5
	NumMovesPerLeg       = 5
)

type undoableMove struct {
	stackBottom *camel
	stackTop    *camel
	srcPos      BoardPosition
}

type Game struct {
	numPlayers int
	// TODO: add player names.
	camelTokens   [NumCamels]camel
	boardSpaces   [BoardSize]boardSpace
	ranking       [NumRacingCamels]Color
	gameOver      bool
	diePyramid    *DiePyramid
	legMovesIndex int
	legCamelMoves [NumMovesPerLeg]undoableMove
}

type GameStateInput struct {
	Players []string // Player names, in order.
	Camels  map[BoardPosition][]Color
	Cheers  map[BoardPosition]string // To player name, for now ignored.
	Boos    map[BoardPosition]string // To player name, for now ignored.
	// The pyramid is only needed for simulations and for computations with
	// some dice already rolled out.
	DiePyramid *DiePyramid
}

type MoveType int

const (
	RollDie MoveType = iota
	PlaceCheer
	PlaceBoo
	BuyTicket
	BetOnWinner
	BetOnLoser
	MakePact
)

type Move struct {
	Type   MoveType
	Player Player
	// Should be a union based on type. For now, only die rolls.
	DieRoll DieRoll
}

func NewGameFromState(i *GameStateInput) (*Game, error) {
	board := &Game{numPlayers: len(i.Players), diePyramid: i.DiePyramid}
	if board.diePyramid == nil {
		board.diePyramid = NewDiePyramid(rand.New(rand.NewSource(*randomSeed)))
	}
	board.legMovesIndex = NumMovesPerLeg - board.diePyramid.RemainingRolls()
	for c := Green; c <= White; c++ {
		board.camelTokens[c].Color = c
		board.camelTokens[c].Position = -1
	}
	board.camelTokens[White].OtherCrazy = &board.camelTokens[Black]
	board.camelTokens[Black].OtherCrazy = &board.camelTokens[White]
	for s := StartPosition; s <= FinishPosition; s++ {
		board.boardSpaces[s].Cheer = NoPlayer
		board.boardSpaces[s].Boo = NoPlayer
	}
	for p, colors := range i.Camels {
		if len(colors) == 0 {
			break
		}
		if p > FinishPosition || p < StartPosition {
			return nil, fmt.Errorf("invalid board position: %d", p)
		}
		space := &board.boardSpaces[p]
		var prev *camel
		for i := 0; i < len(colors); i++ {
			c := colors[i]
			next := &board.camelTokens[c]
			if prev == nil {
				space.StackBottom = next
			} else {
				prev.Next = next
			}
			if next.Position != -1 {
				return nil, fmt.Errorf("%s camel appears twice in input", c)
			}
			next.Position = p
			next.Prev = prev
			space.StackTop = next
			prev = next
		}
	}
	// Check that all camels are accounted for:
	for c := Green; c <= White; c++ {
		if board.camelTokens[c].Position == -1 {
			return nil, fmt.Errorf("%s camel is not placed on the board", c)
		}
	}
	for p := range i.Cheers { // Players ignored for now.
		if p > FinishPosition || p <= StartPosition {
			return nil, fmt.Errorf("invalid board position: %d", p)
		}
		s := &board.boardSpaces[p]
		if s.StackBottom != nil || s.Cheer != NoPlayer || s.Boo != NoPlayer {
			return nil, fmt.Errorf("invalid cheer position %d, not empty", p)
		}
		s.Cheer = 0 // just some player number for now.
	}
	for p := range i.Boos { // Players ignored for now.
		if p > FinishPosition || p <= StartPosition {
			return nil, fmt.Errorf("invalid board position: %d", p)
		}
		s := &board.boardSpaces[p]
		if s.StackBottom != nil || s.Cheer != NoPlayer || s.Boo != NoPlayer {
			return nil, fmt.Errorf("invalid boo position %d, not empty", p)
		}
		s.Boo = 0 // just some player number for now.
	}
	board.computeRanking()
	return board, nil
}

func (g *Game) computeRankingGameOver() {
	// Special cases: if the game is over because a crazy camel crossed
	// over in the opposite direction, all the camels it was carrying (if any)
	// are considered the least advanced.
	// Similarly, when a normal camel crosses over, everyone in its stack are
	// the most advanced (there can be no Boo token on space 0).
	var specialBottom, specialTop *camel
	move := &g.legCamelMoves[g.legMovesIndex-1]
	if move.stackBottom.IsCrazy() {
		if move.stackBottom.Next == nil {
			g.computeRankingRegularCase()
			return
		}
		curRank := Last
		specialBottom = move.stackBottom.Next
		for c := specialBottom; c != move.stackTop.Next; c = c.Next {
			if !c.IsCrazy() {
				specialTop = c
				g.ranking[curRank] = c.Color
				curRank++
			}
		}
		for i := StartPosition; i <= FinishPosition; i++ {
			for c := g.boardSpaces[i].StackBottom; c != nil; c = c.Next {
				if c == specialBottom {
					c = specialTop
				} else if !c.IsCrazy() {
					g.ranking[curRank] = c.Color
					curRank++
				}
			}
		}
	} else {
		curRank := First
		specialTop = move.stackTop
		for c := specialTop; c != move.stackBottom.Prev; c = c.Prev {
			if !c.IsCrazy() {
				specialBottom = c
				g.ranking[curRank] = c.Color
				curRank--
			}
		}
		for i := FinishPosition; i >= StartPosition; i-- {
			for c := g.boardSpaces[i].StackTop; c != nil; c = c.Prev {
				if c == specialTop {
					c = specialBottom
				} else if !c.IsCrazy() {
					g.ranking[curRank] = c.Color
					curRank--
				}
			}
		}
	}
}

func (g *Game) computeRankingRegularCase() {
	curRank := Last
	for i := StartPosition; i <= FinishPosition; i++ {
		for c := g.boardSpaces[i].StackBottom; c != nil; c = c.Next {
			if !c.IsCrazy() {
				g.ranking[curRank] = c.Color
				curRank++
			}
		}
	}
}

func (g *Game) computeRanking() {
	if g.gameOver && g.legMovesIndex > 0 {
		g.computeRankingGameOver()
	} else {
		g.computeRankingRegularCase()
	}
}

func (g *Game) HasCheer(pos BoardPosition) bool {
	return g.boardSpaces[pos].HasCheer()
}

func (g *Game) HasBoo(pos BoardPosition) bool {
	return g.boardSpaces[pos].HasBoo()
}

func (g *Game) moveStack(bottom *camel, top *camel, destPos BoardPosition, pushBelowStack bool) {
	sourceSp := &g.boardSpaces[bottom.Position]
	// Disconnect stack from source space, whether it is currently at the top or bottom:
	if sourceSp.StackTop == top {
		sourceSp.StackTop = bottom.Prev
		if bottom.Prev == nil {
			sourceSp.StackBottom = nil
		} else {
			bottom.Prev.Next = nil
		}
	} else {
		sourceSp.StackBottom = top.Next
		if top.Next == nil {
			sourceSp.StackTop = nil
		} else {
			top.Next.Prev = nil
			top.Next = nil
		}
	}
	destSp := &g.boardSpaces[destPos]
	prevBottom := destSp.StackBottom
	bottom.Prev = nil
	if pushBelowStack {
		destSp.StackBottom = bottom
		top.Next = prevBottom
		if prevBottom == nil {
			destSp.StackTop = top
		} else {
			prevBottom.Prev = top
		}
	} else {
		if prevBottom == nil {
			destSp.StackBottom = bottom
		} else {
			bottom.Prev = destSp.StackTop
			destSp.StackTop.Next = bottom
		}
		destSp.StackTop = top
	}
	for ; bottom != top.Next; bottom = bottom.Next {
		bottom.Position = destPos
	}
}

// Applies move within the current leg of the race. This may
// cause the leg to be over, and/or the game to be over.
// The move is undoable via the undoCamelMove function.
func (g *Game) applyCamelMove(r *DieRoll) {
	move := &g.legCamelMoves[g.legMovesIndex]
	g.legMovesIndex++
	c := &g.camelTokens[r.Color]
	moveDirection := 1
	if c.IsCrazy() {
		// Crazy camels have special rules.
		moveDirection = -1
		other := c.OtherCrazy
		if c.Next == other || (c.Next == nil && other.Next != nil && other.Next != c) {
			c = other
		}
	}
	move.srcPos = c.Position
	move.stackBottom = c
	move.stackTop = g.boardSpaces[c.Position].StackTop
	destPos := c.Position.Add(int(r.Value) * moveDirection)
	if g.HasCheer(destPos) {
		destPos = destPos.Add(moveDirection)
	}
	g.gameOver = int(c.Position-destPos)*moveDirection > 0
	pushBelowStack := false
	if g.HasBoo(destPos) {
		// The game is still over even if we are now below the finish line again.
		destPos = destPos.Add(-moveDirection)
		pushBelowStack = true
	}
	g.moveStack(c, move.stackTop, destPos, pushBelowStack)
	g.computeRanking()
}

func (g *Game) undoLastCamelMove() {
	g.gameOver = false
	g.legMovesIndex--
	move := &g.legCamelMoves[g.legMovesIndex]
	// To undo, we always push to the top of the src stack.
	g.moveStack(move.stackBottom, move.stackTop, move.srcPos, false)
	g.computeRanking()
}

// Computes all the possible outcomes for the current leg.
func (g *Game) ComputeLegRankingDistribution() *RankingDistribution {
	// The remaining weight with N dice in the bag is 6^(N-1) * N!
	remainingWeights := [6]int{1, 2 * 6, 6 * 36, 24 * 216, 120 * 1296}
	powersOf2 := [6]int{1, 2, 4, 8, 16, 32}
	d := &RankingDistribution{}
	if g.diePyramid.IsEmpty() {
		// All dice were rolled: only the current board remains.
		d.RecordRanking(&g.ranking)
		return d
	}
	colors := g.diePyramid.RemainingDice()
	movesInLeg := g.diePyramid.RemainingRolls()
	var used [6]bool
	colorIndices := [NumMovesPerLeg]int{-1, -1, -1, -1, -1}
	var values [NumMovesPerLeg]RollValue
	var roll DieRoll
	for curDie := 0; curDie >= 0; {
		if colorIndices[curDie] >= 0 {
			g.undoLastCamelMove()

			c := colors[colorIndices[curDie]]
			values[curDie]++
			if c == Black && values[curDie] > 6 || c != Black && values[curDie] > 3 {
				// Try to find next unused color:
				used[c] = false
				for i := colorIndices[curDie] + 1; i <= movesInLeg; i++ {
					next := colors[i]
					if !used[next] {
						used[next] = true
						colorIndices[curDie] = i
						values[curDie] = 1
						break
					}
				}
				if values[curDie] > 1 { // not found
					colorIndices[curDie] = -1
					curDie--
					continue
				}
			}
		} else {
			// Choose first unused color.
			for i, c := range colors {
				if !used[c] {
					used[c] = true
					colorIndices[curDie] = i
					break
				}
			}
			values[curDie] = 1
		}
		roll.Color = colors[colorIndices[curDie]]
		roll.Value = values[curDie]
		if roll.Value > 3 {
			roll.Value -= 3
			roll.Color = White
		}
		g.applyCamelMove(&roll)
		if g.gameOver || curDie == movesInLeg-1 {
			weightIndex := curDie
			if !used[Black] {
				weightIndex++
			}
			d.RecordWeightedRanking(&g.ranking, powersOf2[weightIndex]*remainingWeights[movesInLeg-curDie-1])
		} else {
			curDie++
		}
	}
	return d
}

// Simulates the current leg numSamples times. It is implemented in order to
// test/validate the results of ComputeLegRankingDistribution.
func (g *Game) SimulateLegRankingDistribution(numSamples int) *RankingDistribution {
	d := &RankingDistribution{}
	gameCopy := *g
	for s := 0; s < numSamples; s++ {
		for !g.LegOver() {
			r, _ := g.diePyramid.Roll()
			g.applyCamelMove(&r)
		}
		d.RecordRanking(&g.ranking)
		g.diePyramid.Reset()
		*g = gameCopy
	}
	return d
}

func (g *Game) LegOver() bool {
	return g.gameOver || g.legMovesIndex == NumMovesPerLeg
}

// Pretty print the board to the user.
func (g *Game) String() string {
	var s strings.Builder
	for p := StartPosition; p <= FinishPosition; p++ {
		fmt.Fprintf(&s, "%2d: ", p+1)
		sp := &g.boardSpaces[p]
		for c := sp.StackBottom; c != nil; c = c.Next {
			fmt.Fprintf(&s, "%s ", c.Color)
		}
		if sp.HasCheer() {
			s.WriteString(colorPrinters[White](">>"))
		}
		if sp.HasBoo() {
			s.WriteString(colorPrinters[White]("<<"))
		}
		s.WriteString("\n")
	}
	if g.gameOver {
		s.WriteString("Game is over! Final ranking (>>): ")
	} else {
		s.WriteString("Game is in progress. Current ranking (>>): ")
	}
	for _, c := range g.ranking {
		fmt.Fprintf(&s, "%s ", c)
	}
	return s.String()
}

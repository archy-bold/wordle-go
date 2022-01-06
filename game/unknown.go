package game

import (
	"fmt"

	"github.com/pkg/errors"
)

type UnknownGame struct {
	complete bool
	attempts int
	length   int
	grid     Grid
}

func (g *UnknownGame) Play(word string) (bool, error) {
	// Do nothing, we need to know the result for an unknown game
	return false, nil
}

func (g *UnknownGame) AddResult(row []GridCell) (bool, error) {
	// Handle errors where the length isn't right
	if len(row) != g.length {
		return false, errors.Wrap(ErrWrongWordLength, fmt.Sprint(len(row)))
	}

	g.grid[g.attempts] = row
	g.attempts++

	// Check if it's a winner
	numCorrect := 0
	for _, cell := range row {
		if cell.Status == STATUS_CORRECT {
			numCorrect++
		}
	}

	if numCorrect == g.length {
		g.complete = true
	}

	return g.complete, nil
}

func (g *UnknownGame) HasEnded() bool {
	return g.complete || g.attempts == len(g.grid)
}

func (g *UnknownGame) GetScore() (int, int) {
	return g.attempts, len(g.grid)
}

func (g *UnknownGame) GetLastPlay() []GridCell {
	if g.attempts == 0 {
		return nil
	}
	return g.grid[g.attempts-1]
}

func (g *UnknownGame) OutputForConsole() string {
	return outputGridForConsole(g.grid, g.length)
}

func (g *UnknownGame) OutputToShare() string {
	score := fmt.Sprint(g.attempts)
	if !g.complete && g.HasEnded() {
		score = "X"
	}
	return outputGridToShare(g.grid, score, len(g.grid))
}

// CreateGame creates a game for the given answer and number of allowed tries
func CreateUnknownGame(length int, tries int) Game {
	// TODO include valid entries
	grid := make(Grid, tries)

	return &UnknownGame{false, 0, length, grid}
}

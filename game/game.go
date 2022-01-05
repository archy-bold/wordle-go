package game

import (
	"strings"
)

// CellStatus represents whether a cell is in a state of empty, correct, incorrect or wrong
type CellStatus int

const (
	STATUS_EMPTY     CellStatus = 0
	STATUS_CORRECT   CellStatus = 1
	STATUS_INCORRECT CellStatus = 2
	STATUS_WRONG     CellStatus = 3
)

// GridCell represents a cell within a game grid
type GridCell struct {
	Letter string
	Status CellStatus
}

// Grid represents a game grid
type Grid [][]GridCell

// Game represents a game that can be played
type Game interface {
	// Play plays a word in the game
	Play(word string) (bool, error)
	// HasEnded returns whether the game has ended, whether with success or failure
	HasEnded() bool
	// GetScore gets the running or final score for the game
	GetScore() (int, int)
	// GetLastPlay returns the result of the last play
	GetLastPlay() []GridCell
	// OutputForConsole returns a string representation of the game for the command line
	OutputForConsole() string
}

type game struct {
	complete bool
	attempts int
	answer   string
	grid     Grid
}

func (g *game) Play(word string) (bool, error) {
	// TODO check the word length here and error if too long/short

	// Create the row for the grid
	parts := strings.Split(word, "")
	answerParts := strings.Split(g.answer, "")
	row := make([]GridCell, len(parts))
	numCorrect := 0
	for i, chr := range parts {
		var status CellStatus
		if chr == answerParts[i] {
			status = STATUS_CORRECT
			numCorrect++
		} else if stringInSlice(chr, answerParts) {
			status = STATUS_INCORRECT
		} else {
			status = STATUS_WRONG
		}

		row[i] = GridCell{chr, status}
	}

	// Update the game
	g.grid[g.attempts] = row
	g.attempts++

	return word == g.answer, nil
}

func (g *game) HasEnded() bool {
	return g.complete || g.attempts == len(g.grid)
}

func (g *game) GetScore() (int, int) {
	return g.attempts, len(g.grid)
}

func (g *game) GetLastPlay() []GridCell {
	return g.grid[g.attempts-1]
}

func (g *game) OutputForConsole() string {
	str := "\n" + strings.Repeat("-", len(g.answer)+2) + "\n"
	for _, row := range g.grid {
		if len(row) == 0 {
			break
		}

		str += "|"
		for _, cell := range row {
			switch cell.Status {
			case STATUS_CORRECT:
				str += COLOUR_GREEN
			case STATUS_INCORRECT:
				str += COLOUR_YELLOW
			}
			str += cell.Letter + COLOUR_RESET
		}
		str += "|\n"
	}
	str += strings.Repeat("-", len(g.answer)+2) + "\n"

	return str
}

// CreateGame creates a game for the given answer and number of allowed tries
func CreateGame(answer string, tries int) Game {
	// TODO include valid entries
	grid := make([][]GridCell, tries)

	return &game{false, 0, answer, grid}
}
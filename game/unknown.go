package game

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
	// TODO handle errors where the length isn't right
	g.grid = append(g.grid, row)
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
	return g.grid[g.attempts-1]
}

func (g *UnknownGame) OutputForConsole() string {
	return outputGridForConsole(g.grid, g.length)
}

func (g *UnknownGame) OutputToShare() string {
	return outputGridToShare(g.grid, g.attempts, len(g.grid))
}

// CreateGame creates a game for the given answer and number of allowed tries
func CreateUnknownGame(length int, tries int) Game {
	// TODO include valid entries
	grid := make([][]GridCell, tries)

	return &UnknownGame{false, 0, length, grid}
}

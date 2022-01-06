package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	gameTapirStart    = &game{false, 0, "tapir", make(Grid, 6)}
	gameTapirFinished = &game{true, 4, "tapir", Grid{
		{GridCell{"g", STATUS_WRONG}, GridCell{"r", STATUS_INCORRECT}, GridCell{"o", STATUS_WRONG}, GridCell{"u", STATUS_WRONG}, GridCell{"p", STATUS_INCORRECT}},
		{GridCell{"p", STATUS_INCORRECT}, GridCell{"r", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}, GridCell{"n", STATUS_WRONG}, GridCell{"k", STATUS_WRONG}},
		{GridCell{"s", STATUS_WRONG}, GridCell{"p", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}, GridCell{"r", STATUS_INCORRECT}, GridCell{"e", STATUS_WRONG}},
		{GridCell{"t", STATUS_CORRECT}, GridCell{"a", STATUS_CORRECT}, GridCell{"p", STATUS_CORRECT}, GridCell{"i", STATUS_CORRECT}, GridCell{"r", STATUS_CORRECT}},
		nil,
		nil,
	}}
	gameAtStart    = &game{false, 0, "at", make(Grid, 1)}
	gameAtFinished = &game{false, 1, "at", Grid{
		{GridCell{"t", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}},
	}}
)

var createGameTests = map[string]struct {
	answer   string
	tries    int
	expected *game
}{
	"5 letter, 6 tries": {"tapir", 6, &game{false, 0, "tapir", Grid{nil, nil, nil, nil, nil, nil}}},
	"3 letter, 3 tries": {"bat", 3, &game{false, 0, "bat", Grid{nil, nil, nil}}},
	// "5 letter, uppercase": {"TAPIR", 6, &game{false, 0, "tapir", Grid{nil, nil, nil, nil, nil, nil}}},
}

func Test_CreateGame(t *testing.T) {
	for tn, tt := range createGameTests {
		g := CreateGame(tt.answer, tt.tries)

		assert.Equalf(t, tt.expected, g, "Expected game to match for test '%s'", tn)
	}
}

var gamePlayTests = map[string]struct {
	g            *game
	tries        []string
	expected     []bool
	expectedErr  string
	expectedGrid Grid
}{
	"5-letter, won": {
		g:            gameTapirStart,
		tries:        []string{"group", "prank", "spare", "tapir"},
		expected:     []bool{false, false, false, true},
		expectedGrid: gameTapirFinished.grid,
	},
	"2-letter, lost": {
		g:            gameAtStart,
		tries:        []string{"ta"},
		expected:     []bool{false},
		expectedGrid: gameAtFinished.grid,
	},
}

func Test_game_Play(t *testing.T) {
	for tn, tt := range gamePlayTests {
		// Copy first
		g := &game{tt.g.complete, tt.g.attempts, tt.g.answer, make(Grid, len(tt.g.grid))}
		for i, row := range tt.g.grid {
			copy(g.grid[i], row)
		}
		for i, word := range tt.tries {
			res, _ := g.Play(word)

			// Make the assertions
			assert.Equalf(t, tt.expected[i], res, "Expected play outcome to match for test '%s'", tn)
			assert.Equalf(t, tt.expected[i], g.complete, "Expected complete to match for test '%s'", tn)
			assert.Equalf(t, tt.expectedGrid[i], g.grid[i], "Expected grid row to match for test '%s'", tn)
			assert.Equalf(t, i+1, g.attempts, "Expected attempts to match for test '%s'", tn)
		}
	}
}

var gameHasEndedTests = map[string]struct {
	g        *game
	expected bool
}{
	"5-letter start":    {gameTapirStart, false},
	"5-letter finished": {gameTapirFinished, true},
	"2-letter start":    {gameAtStart, false},
	"2-letter finished": {gameAtFinished, true},
}

func Test_game_HasEnded(t *testing.T) {
	for tn, tt := range gameHasEndedTests {
		assert.Equalf(t, tt.expected, tt.g.HasEnded(), "Expected result to match for test '%s'", tn)
	}
}

var gameGetScoreTests = map[string]struct {
	g             *game
	expectedScore int
	expectedOf    int
}{
	"5-letter start":    {gameTapirStart, 0, 6},
	"5-letter finished": {gameTapirFinished, 4, 6},
	"2-letter start":    {gameAtStart, 0, 1},
	"2-letter finished": {gameAtFinished, 1, 1},
}

func Test_game_GetScore(t *testing.T) {
	for tn, tt := range gameGetScoreTests {
		score, of := tt.g.GetScore()

		assert.Equalf(t, tt.expectedScore, score, "Expected score to match for test '%s'", tn)
		assert.Equalf(t, tt.expectedOf, of, "Expected of to match for test '%s'", tn)
	}
}

var gameGetLastPlayTests = map[string]struct {
	g        *game
	expected []GridCell
}{
	"5-letter start":    {gameTapirStart, nil},
	"5-letter finished": {gameTapirFinished, gameTapirFinished.grid[3]},
	"2-letter start":    {gameAtStart, nil},
	"2-letter finished": {gameAtFinished, gameAtFinished.grid[0]},
}

func Test_game_GetLastPlay(t *testing.T) {
	for tn, tt := range gameGetLastPlayTests {
		assert.Equalf(t, tt.expected, tt.g.GetLastPlay(), "Expected result to match for test '%s'", tn)
	}
}

var gameOutputForConsoleTests = map[string]struct {
	g        *game
	expected string
}{
	"5-letter start": {
		g:        gameTapirStart,
		expected: "\n-------\n-------\n",
	},
	"5-letter finished": {
		g: gameTapirFinished,
		expected: "\n-------\n" +
			"|g" + COLOUR_RESET + COLOUR_YELLOW + "r" + COLOUR_RESET + "o" + COLOUR_RESET + "u" + COLOUR_RESET + COLOUR_YELLOW + "p" + COLOUR_RESET + "|\n" +
			"|" + COLOUR_YELLOW + "p" + COLOUR_RESET + COLOUR_YELLOW + "r" + COLOUR_RESET + COLOUR_YELLOW + "a" + COLOUR_RESET + "n" + COLOUR_RESET + "k" + COLOUR_RESET + "|\n" +
			"|s" + COLOUR_RESET + COLOUR_YELLOW + "p" + COLOUR_RESET + COLOUR_YELLOW + "a" + COLOUR_RESET + COLOUR_YELLOW + "r" + COLOUR_RESET + "e" + COLOUR_RESET + "|\n" +
			"|" + COLOUR_GREEN + "t" + COLOUR_RESET + COLOUR_GREEN + "a" + COLOUR_RESET + COLOUR_GREEN + "p" + COLOUR_RESET + COLOUR_GREEN + "i" + COLOUR_RESET + COLOUR_GREEN + "r" + COLOUR_RESET + "|\n" +
			"-------\n",
	},
	"2-letter start": {
		g:        gameAtStart,
		expected: "\n----\n----\n",
	},
	"2-letter finished": {
		g:        gameAtFinished,
		expected: "\n----\n|" + COLOUR_YELLOW + "t" + COLOUR_RESET + COLOUR_YELLOW + "a" + COLOUR_RESET + "|\n----\n",
	},
}

func Test_game_OutputForConsole(t *testing.T) {
	for tn, tt := range gameOutputForConsoleTests {
		assert.Equalf(t, tt.expected, tt.g.OutputForConsole(), "Expected result to match for test '%s'", tn)
	}
}

var gameOutputToShareTests = map[string]struct {
	g        *game
	expected string
}{
	"5-letter start": {
		g:        gameTapirStart,
		expected: "Wordle 0/6\n\n\n",
	},
	"5-letter finished": {
		g: gameTapirFinished,
		expected: "Wordle 4/6\n\n" +
			"⬜🟨⬜⬜🟨\n" +
			"🟨🟨🟨⬜⬜\n" +
			"⬜🟨🟨🟨⬜\n" +
			"🟩🟩🟩🟩🟩\n\n",
	},
	"2-letter start": {
		g:        gameAtStart,
		expected: "Wordle 0/1\n\n\n",
	},
	"2-letter finished": {
		g:        gameAtFinished,
		expected: "Wordle X/1\n\n🟨🟨\n\n",
	},
}

func Test_game_OutputToShare(t *testing.T) {
	for tn, tt := range gameOutputToShareTests {
		assert.Equalf(t, tt.expected, tt.g.OutputToShare(), "Expected result to match for test '%s'", tn)
	}
}

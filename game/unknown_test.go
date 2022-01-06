package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	unknownGameTapirStart    = &UnknownGame{false, 0, 5, make(Grid, 6)}
	unknownGameTapirFinished = &UnknownGame{true, 4, 5, Grid{
		{GridCell{"g", STATUS_WRONG}, GridCell{"r", STATUS_INCORRECT}, GridCell{"o", STATUS_WRONG}, GridCell{"u", STATUS_WRONG}, GridCell{"p", STATUS_INCORRECT}},
		{GridCell{"p", STATUS_INCORRECT}, GridCell{"r", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}, GridCell{"n", STATUS_WRONG}, GridCell{"k", STATUS_WRONG}},
		{GridCell{"s", STATUS_WRONG}, GridCell{"p", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}, GridCell{"r", STATUS_INCORRECT}, GridCell{"e", STATUS_WRONG}},
		{GridCell{"t", STATUS_CORRECT}, GridCell{"a", STATUS_CORRECT}, GridCell{"p", STATUS_CORRECT}, GridCell{"i", STATUS_CORRECT}, GridCell{"r", STATUS_CORRECT}},
		nil,
		nil,
	}}
	unknownGameAtStart    = &UnknownGame{false, 0, 2, make(Grid, 1)}
	unknownGameAtFinished = &UnknownGame{false, 1, 2, Grid{
		{GridCell{"t", STATUS_INCORRECT}, GridCell{"a", STATUS_INCORRECT}},
	}}
)

var createUnknownGameTests = map[string]struct {
	length   int
	tries    int
	expected *UnknownGame
}{
	"5 letter, 6 tries": {5, 6, &UnknownGame{false, 0, 5, Grid{nil, nil, nil, nil, nil, nil}}},
	"3 letter, 3 tries": {3, 3, &UnknownGame{false, 0, 3, Grid{nil, nil, nil}}},
}

func Test_CreateUnknownGame(t *testing.T) {
	for tn, tt := range createUnknownGameTests {
		g := CreateUnknownGame(tt.length, tt.tries)

		assert.Equalf(t, tt.expected, g, "Expected game to match for test '%s'", tn)
	}
}

var unknownGamePlayTests = map[string]struct {
	g        *UnknownGame
	tries    []string
	expected []bool
}{
	"5-letter": {
		g:        unknownGameTapirStart,
		tries:    []string{"group", "prank", "spare", "tapir"},
		expected: []bool{false, false, false, false},
	},
	"5-letter, mixed case, mixed length": {
		g:        unknownGameTapirStart,
		tries:    []string{"grOUp", "PRAnk", "strong", "TAPIR", "tape"},
		expected: []bool{false, false, false, false, false},
	},
	"2-letter": {
		g:        unknownGameAtStart,
		tries:    []string{"ta"},
		expected: []bool{false},
	},
}

func Test_UnknownGame_Play(t *testing.T) {
	for tn, tt := range unknownGamePlayTests {
		// Copy first
		g := &UnknownGame{tt.g.complete, tt.g.attempts, tt.g.length, make(Grid, len(tt.g.grid))}
		for i, row := range tt.g.grid {
			copy(g.grid[i], row)
		}
		for i, word := range tt.tries {
			res, err := g.Play(word)

			// Make the assertions
			assert.NoErrorf(t, err, "Expected nil error for test '%s'", tn)
			assert.Equalf(t, tt.expected[i], res, "Expected play outcome to match for test '%s'", tn)
		}
	}
}

var unknownGameAddResultTests = map[string]struct {
	g            *UnknownGame
	results      [][]GridCell
	expected     []bool
	expectedErr  string
	expectedGrid Grid
}{
	"5-letter, won": {
		g: unknownGameTapirStart,
		results: [][]GridCell{
			unknownGameTapirFinished.grid[0],
			unknownGameTapirFinished.grid[1],
			unknownGameTapirFinished.grid[2],
			unknownGameTapirFinished.grid[3],
		},
		expected:     []bool{false, false, false, true},
		expectedGrid: unknownGameTapirFinished.grid,
	},
	"5-letter, try 4-letter word": {
		g:           unknownGameTapirStart,
		results:     [][]GridCell{{GridCell{}, GridCell{}, GridCell{}, GridCell{}}},
		expectedErr: "The entered word length is wrong, should be: 5",
	},
	"5-letter, try 6-letter word": {
		g:           unknownGameTapirStart,
		results:     [][]GridCell{{GridCell{}, GridCell{}, GridCell{}, GridCell{}, GridCell{}, GridCell{}}},
		expectedErr: "The entered word length is wrong, should be: 5",
	},
	"2-letter, lost": {
		g:            unknownGameAtStart,
		results:      [][]GridCell{unknownGameAtFinished.grid[0]},
		expected:     []bool{false},
		expectedGrid: unknownGameAtFinished.grid,
	},
}

func Test_UnkownGame_AddResult(t *testing.T) {
	for tn, tt := range unknownGameAddResultTests {
		// Copy first
		g := &UnknownGame{tt.g.complete, tt.g.attempts, tt.g.length, make(Grid, len(tt.g.grid))}
		for i, row := range tt.results {
			res, err := g.AddResult(row)

			// Make the assertions
			if tt.expectedErr != "" {
				assert.Falsef(t, res, "Expected res false for test '%s', result %d", tn, i)
				assert.Errorf(t, err, "Expected error to match for test '%s', result %d", tn, i)
			} else {
				assert.NoErrorf(t, err, "Expected nil error for test '%s', result %d", tn, i)
				assert.Equalf(t, tt.expected[i], res, "Expected play outcome to match for test '%s', result %d", tn, i)
				assert.Equalf(t, tt.expected[i], g.complete, "Expected complete to match for test '%s', result %d", tn, i)
				assert.Equalf(t, tt.expectedGrid[i], g.grid[i], "Expected grid row to match for test '%s', result %d", tn, i)
				assert.Equalf(t, i+1, g.attempts, "Expected attempts to match for test '%s', result %d", tn, i)
			}
		}
	}
}

var unknownGameHasEndedTests = map[string]struct {
	g        *UnknownGame
	expected bool
}{
	"5-letter start":    {unknownGameTapirStart, false},
	"5-letter finished": {unknownGameTapirFinished, true},
	"2-letter start":    {unknownGameAtStart, false},
	"2-letter finished": {unknownGameAtFinished, true},
}

func Test_UnknownGame_HasEnded(t *testing.T) {
	for tn, tt := range unknownGameHasEndedTests {
		assert.Equalf(t, tt.expected, tt.g.HasEnded(), "Expected result to match for test '%s'", tn)
	}
}

var unknownGameGetScoreTests = map[string]struct {
	g             *UnknownGame
	expectedScore int
	expectedOf    int
}{
	"5-letter start":    {unknownGameTapirStart, 0, 6},
	"5-letter finished": {unknownGameTapirFinished, 4, 6},
	"2-letter start":    {unknownGameAtStart, 0, 1},
	"2-letter finished": {unknownGameAtFinished, 1, 1},
}

func Test_UnknownGame_GetScore(t *testing.T) {
	for tn, tt := range unknownGameGetScoreTests {
		score, of := tt.g.GetScore()

		assert.Equalf(t, tt.expectedScore, score, "Expected score to match for test '%s'", tn)
		assert.Equalf(t, tt.expectedOf, of, "Expected of to match for test '%s'", tn)
	}
}

var unknownGameGetLastPlayTests = map[string]struct {
	g        *UnknownGame
	expected []GridCell
}{
	"5-letter start":    {unknownGameTapirStart, nil},
	"5-letter finished": {unknownGameTapirFinished, unknownGameTapirFinished.grid[3]},
	"2-letter start":    {unknownGameAtStart, nil},
	"2-letter finished": {unknownGameAtFinished, unknownGameAtFinished.grid[0]},
}

func Test_UnknownGame_GetLastPlay(t *testing.T) {
	for tn, tt := range unknownGameGetLastPlayTests {
		assert.Equalf(t, tt.expected, tt.g.GetLastPlay(), "Expected result to match for test '%s'", tn)
	}
}

var unknownGameOutputForConsoleTests = map[string]struct {
	g        *UnknownGame
	expected string
}{
	"5-letter start": {
		g:        unknownGameTapirStart,
		expected: "\n-------\n-------\n",
	},
	"5-letter finished": {
		g: unknownGameTapirFinished,
		expected: "\n-------\n" +
			"|G" + COLOUR_RESET + COLOUR_YELLOW + "R" + COLOUR_RESET + "O" + COLOUR_RESET + "U" + COLOUR_RESET + COLOUR_YELLOW + "P" + COLOUR_RESET + "|\n" +
			"|" + COLOUR_YELLOW + "P" + COLOUR_RESET + COLOUR_YELLOW + "R" + COLOUR_RESET + COLOUR_YELLOW + "A" + COLOUR_RESET + "N" + COLOUR_RESET + "K" + COLOUR_RESET + "|\n" +
			"|S" + COLOUR_RESET + COLOUR_YELLOW + "P" + COLOUR_RESET + COLOUR_YELLOW + "A" + COLOUR_RESET + COLOUR_YELLOW + "R" + COLOUR_RESET + "E" + COLOUR_RESET + "|\n" +
			"|" + COLOUR_GREEN + "T" + COLOUR_RESET + COLOUR_GREEN + "A" + COLOUR_RESET + COLOUR_GREEN + "P" + COLOUR_RESET + COLOUR_GREEN + "I" + COLOUR_RESET + COLOUR_GREEN + "R" + COLOUR_RESET + "|\n" +
			"-------\n",
	},
	"2-letter start": {
		g:        unknownGameAtStart,
		expected: "\n----\n----\n",
	},
	"2-letter finished": {
		g:        unknownGameAtFinished,
		expected: "\n----\n|" + COLOUR_YELLOW + "T" + COLOUR_RESET + COLOUR_YELLOW + "A" + COLOUR_RESET + "|\n----\n",
	},
}

func Test_UnknownGame_OutputForConsole(t *testing.T) {
	for tn, tt := range unknownGameOutputForConsoleTests {
		assert.Equalf(t, tt.expected, tt.g.OutputForConsole(), "Expected result to match for test '%s'", tn)
	}
}

var unknownGameOutputToShareTests = map[string]struct {
	g        *UnknownGame
	expected string
}{
	"5-letter start": {
		g:        unknownGameTapirStart,
		expected: "Wordle 0/6\n\n\n",
	},
	"5-letter finished": {
		g: unknownGameTapirFinished,
		expected: "Wordle 4/6\n\n" +
			"⬜🟨⬜⬜🟨\n" +
			"🟨🟨🟨⬜⬜\n" +
			"⬜🟨🟨🟨⬜\n" +
			"🟩🟩🟩🟩🟩\n\n",
	},
	"2-letter start": {
		g:        unknownGameAtStart,
		expected: "Wordle 0/1\n\n\n",
	},
	"2-letter finished": {
		g:        unknownGameAtFinished,
		expected: "Wordle X/1\n\n🟨🟨\n\n",
	},
}

func Test_UnknownGame_OutputToShare(t *testing.T) {
	for tn, tt := range unknownGameOutputToShareTests {
		assert.Equalf(t, tt.expected, tt.g.OutputToShare(), "Expected result to match for test '%s'", tn)
	}
}

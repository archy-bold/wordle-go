package game

import "strings"

var letters = map[string]int{
	"Q": 0, "W": 1, "E": 2, "R": 3, "T": 4, "Y": 5, "U": 6, "I": 7, "O": 8, "P": 9, "br1": 10,
	"A": 11, "S": 12, "D": 13, "F": 14, "G": 15, "H": 16, "J": 17, "K": 18, "L": 19, "br2": 20,
	"Z": 21, "X": 22, "C": 23, "V": 24, "B": 25, "N": 26, "M": 27,
}

// GridCell represents a cell within a game grid
type Key struct {
	Letter string
	Status CellStatus
}

// Keyboard is a representation of a game keyboard, where keys have state
type Keyboard interface {
	// GetKeyState gets the state of a specific key
	GetKeyState(key string) CellStatus
	// SetKeyState sets the state of a specific key
	SetKeyState(key string, status CellStatus)
	// OutputForConsole returns a string representation of the keyboard for the command line
	OutputForConsole() string
}

type keyboard struct {
	keys []Key
}

func (kb *keyboard) GetKeyState(keyChar string) CellStatus {
	keyChar = strings.ToUpper(keyChar)
	pos := letters[keyChar]
	return kb.keys[pos].Status
}

func (kb *keyboard) SetKeyState(keyChar string, status CellStatus) {
	keyChar = strings.ToUpper(keyChar)
	pos := letters[keyChar]
	key := kb.keys[pos]
	key.Status = status
	kb.keys[pos] = key
}

func (kb *keyboard) OutputForConsole() string {
	// Q W E R T Y U I O P
	//  A S D F G H J K L
	//   Z X C V B N M
	str := ""
	for _, key := range kb.keys {
		switch key.Status {
		case STATUS_CORRECT:
			str += COLOUR_GREEN
		case STATUS_INCORRECT:
			str += COLOUR_YELLOW
		case STATUS_WRONG:
			str += COLOUR_GREY
		}
		str += " " + key.Letter + COLOUR_RESET
	}
	str += "\n"
	return str
}

func newKeyboard() Keyboard {
	return &keyboard{
		keys: []Key{
			{Letter: "Q"}, {Letter: "W"}, {Letter: "E"}, {Letter: "R"}, {Letter: "T"}, {Letter: "Y"}, {Letter: "U"}, {Letter: "I"}, {Letter: "O"}, {Letter: "P"}, {Letter: "\n "},
			{Letter: "A"}, {Letter: "S"}, {Letter: "D"}, {Letter: "F"}, {Letter: "G"}, {Letter: "H"}, {Letter: "J"}, {Letter: "K"}, {Letter: "L"}, {Letter: "\n  "},
			{Letter: "Z"}, {Letter: "X"}, {Letter: "C"}, {Letter: "V"}, {Letter: "B"}, {Letter: "N"}, {Letter: "M"},
		},
	}
}

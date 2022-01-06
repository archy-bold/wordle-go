package game

import "errors"

// ErrWrongWordLength error to represent when a word of the incorrect length is entered
var ErrWrongWordLength = errors.New("The entered word length is wrong, should be")

// ErrInvalidWord error to represent when an invalid is entered
var ErrInvalidWord = errors.New("The entered word is invalid")

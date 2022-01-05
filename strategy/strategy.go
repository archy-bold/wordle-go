package strategy

import "github.com/archy-bold/wordle-go/game"

// Strategy represents an object that determine optimal next moves for the strategy
type Strategy interface {
	// GetNextMove will get the next move for the given strategy
	GetNextMove() string
	// SetMoveOutcome is to tell the strategy the outcome of the last move
	SetMoveOutcome(row []game.GridCell)
	// GetSuggestions will get the best n suggestions given the current state
	GetSuggestions(n int) PairList
}

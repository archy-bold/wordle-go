package strategy

import "github.com/archy-bold/wordle-go/game"

// Strategy represents an object that determine optimal next moves for the strategy
type Strategy interface {
	// GetNextMove will get the next move for the given strategy
	GetNextMove() string
	// SetMoveOutcome is to tell the strategy the outcome of the last move
	SetMoveOutcome(row []game.GridCell)
}

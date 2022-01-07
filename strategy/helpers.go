package strategy

import (
	"strings"

	"github.com/archy-bold/wordle-go/game"
)

type hintSummary struct {
	allIncorrect     []string
	answersCorrect   []string
	answersIncorrect []string
	allEliminated    map[string]bool
}

func getHintSummary(hint []game.GridCell) hintSummary {
	hs := hintSummary{
		allIncorrect:     make([]string, 0),
		answersCorrect:   make([]string, len(hint)),
		answersIncorrect: make([]string, len(hint)),
		allEliminated:    map[string]bool{},
	}
	for i, h := range hint {
		switch h.Status {
		case game.STATUS_INCORRECT:
			hs.allIncorrect = append(hs.allIncorrect, h.Letter)
			hs.answersIncorrect[i] = h.Letter
		case game.STATUS_WRONG:
			hs.allEliminated[h.Letter] = true
		case game.STATUS_CORRECT:
			hs.answersCorrect[i] = h.Letter
		}
	}

	return hs
}

func isWordValid(word string, hs hintSummary) bool {
	chrs := strings.Split(word, "")

	if len(hs.allIncorrect) > 0 {
		diff := difference(hs.allIncorrect, chrs)
		if len(diff) > 0 {
			return false
		}
	}

	for i, chr := range chrs {
		// Filter out if this char has been eliminated
		if eliminated, ok := hs.allEliminated[chr]; ok && eliminated {
			return false
		}

		// Filter out if this is position has been found to be correct and this char doesn't match
		if hs.answersCorrect[i] != "" && hs.answersCorrect[i] != chr {
			return false
		}

		// Also filter out if there's an incorrect answer in this position
		if hs.answersIncorrect[i] == chr {
			return false
		}
	}

	return true
}

func filterWordsList(words []string, hint []game.GridCell) []string {
	filteredWords := make([]string, 0)

	hs := getHintSummary(hint)

	for _, word := range words {
		// Add to the list if the word is still possible
		if isWordValid(word, hs) {
			filteredWords = append(filteredWords, word)
		}
	}

	return filteredWords
}

func filterWordsListCount(words []string, hint []game.GridCell) int {
	numMatches := 0
	hs := getHintSummary(hint)

	for _, word := range words {
		// Add to the list if the word is still possible
		if isWordValid(word, hs) {
			numMatches++
		}
	}

	return numMatches
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

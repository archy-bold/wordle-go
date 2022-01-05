package strategy

import (
	"sort"
	"strings"

	"github.com/archy-bold/wordle-go/game"
)

// HistogramEntry represents an entry in the character frequency histogram
type HistogramEntry struct {
	Occurences            int
	OccurrencesInPosition []int
}

// CharFrequencyStrategy a strategy that plays based on the frequency of characters in the solutions list
type CharFrequencyStrategy struct {
	wordLength          int
	letters             map[string]bool
	validWords          *[]string
	histogram           map[string]HistogramEntry
	rankedWords         PairList
	answersCorrect      []string
	answersIncorrect    [][]string
	answersIncorrectAll []string
}

func (s *CharFrequencyStrategy) buildHistogram() {
	s.histogram = make(map[string]HistogramEntry, len(s.letters))
	for l := range s.letters {
		s.histogram[l] = HistogramEntry{0, make([]int, s.wordLength)}
	}

	// Loop through each word and check which unique letters are in the word
	for _, word := range *s.validWords {
		chrs := strings.Split(word, "")

		// Loop through each char and update the histogram
		checkedChars := map[string]bool{}
		for i, chr := range chrs {
			// Ignore removed letters
			if !s.letters[chr] {
				continue
			}
			// If we've not already processed this letter
			if _, ok := checkedChars[chr]; !ok {
				checkedChars[chr] = true
				if entry, ok2 := s.histogram[chr]; ok2 {
					entry.Occurences++
					entry.OccurrencesInPosition[i]++
					s.histogram[chr] = entry
				}
			}
		}
	}
}

func (s *CharFrequencyStrategy) rankWords() {
	s.rankedWords = make(PairList, len(*s.validWords))

word:
	for _, word := range *s.validWords {
		chrs := strings.Split(word, "")
		// First set the score based on the letters that exist
		// TODO score based on letter position too
		checkedChars := map[string]bool{}
		score := 0

		// Check if any of the incorrect answers don't appear in the word
		if len(s.answersIncorrectAll) > 0 {
			diff := difference(s.answersIncorrectAll, chrs)
			if len(diff) > 0 {
				s.rankedWords = append(s.rankedWords, Pair{word, 0})
				continue word
			}
		}

		for i, chr := range chrs {
			// If this is an eliminated letter, score down
			if !s.letters[chr] {
				s.rankedWords = append(s.rankedWords, Pair{word, 0})
				continue word
			}

			// If there is an answer in this position, we can disregard words that don't have that letter in that position
			if s.answersCorrect[i] != "" && s.answersCorrect[i] != chr {
				s.rankedWords = append(s.rankedWords, Pair{word, 0})
				continue word
			}

			// Also check if there's an incorrect answer in this position
			if len(s.answersIncorrect[i]) > 0 {
				for _, ai := range s.answersIncorrect[i] {
					if ai == chr {
						s.rankedWords = append(s.rankedWords, Pair{word, 0})
						continue word
					}
				}
			}

			if _, ok := checkedChars[chr]; !ok {
				// Score based on occurences and occurences in the position
				scoreToAdd := s.histogram[chr].Occurences + s.histogram[chr].OccurrencesInPosition[i]
				// Increase score for incorrectly placed letters
				for _, aiChr := range s.answersIncorrectAll {
					if chr == aiChr {
						scoreToAdd *= 2
						break
					}
				}
				score += scoreToAdd
				checkedChars[chr] = true
			}
		}

		// Add to the ranked list
		s.rankedWords = append(s.rankedWords, Pair{word, score})
	}

	// Sort the
	sort.Sort(sort.Reverse(s.rankedWords))
}

// GetNextMove simply returns the top-ranked word
func (s *CharFrequencyStrategy) GetNextMove() string {
	return s.rankedWords[0].Key
}

func (s *CharFrequencyStrategy) SetMoveOutcome(row []game.GridCell) {
	// Update the internal state for the row
	numCorrect := 0
	rejected := make([]bool, s.wordLength)
	for i, cell := range row {
		switch cell.Status {
		case game.STATUS_WRONG:
			// If this letter has shown up before but not rejected, don't eliminate
			shouldEliminate := true
			for j := 0; j < i; j++ {
				if cell.Letter == row[j].Letter && !rejected[j] {
					shouldEliminate = false
				}
			}
			if shouldEliminate {
				s.letters[cell.Letter] = false
			}
			rejected[i] = true
		case game.STATUS_INCORRECT:
			s.answersIncorrect[i] = append(s.answersIncorrect[i], cell.Letter)
			s.answersIncorrectAll = append(s.answersIncorrectAll, cell.Letter)
		case game.STATUS_CORRECT:
			s.answersCorrect[i] = cell.Letter
			numCorrect++
		}
	}

	// Rebuild the histogram and ranking
	s.buildHistogram()
	s.rankWords()
}

// NewCharFrequencyStrategy create a char frequency-based strategy given the word list and letters list
func NewCharFrequencyStrategy(wordLength int, letters []string, validWords *[]string) Strategy {
	lettersMap := map[string]bool{}
	for _, l := range letters {
		lettersMap[l] = true
	}

	s := &CharFrequencyStrategy{
		wordLength:       wordLength,
		letters:          lettersMap,
		validWords:       validWords,
		answersCorrect:   make([]string, wordLength),
		answersIncorrect: make([][]string, wordLength),
	}

	// Initialise the histogram
	s.buildHistogram()
	// Rank the words
	s.rankWords()

	return s
}

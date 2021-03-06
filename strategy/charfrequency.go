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
	attempts            int
	starter             string
	wordLength          int
	letters             map[string]bool
	possibleAnswers     []string
	allWords            *[]string
	histogram           map[string]HistogramEntry
	rankedWords         PairList
	answersIncorrectAll []string
}

func (s *CharFrequencyStrategy) buildHistogram() {
	s.histogram = make(map[string]HistogramEntry, len(s.letters))
	for l := range s.letters {
		s.histogram[l] = HistogramEntry{0, make([]int, s.wordLength)}
	}

	// Loop through each word and check which unique letters are in the word
	for _, word := range s.possibleAnswers {
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
	s.rankedWords = make(PairList, len(s.possibleAnswers))
	words := s.possibleAnswers

	for i, word := range words {
		chrs := strings.Split(word, "")
		// First set the score based on the letters that exist
		checkedChars := map[string]bool{}
		score := 0

		for i, chr := range chrs {
			if _, ok := checkedChars[chr]; !ok {
				// Score based on occurences and occurences in the position
				scoreToAdd := s.histogram[chr].Occurences + (s.histogram[chr].OccurrencesInPosition[i] * 10)
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
		s.rankedWords[i] = Pair{word, score}
	}

	// Sort the
	sort.Sort(sort.Reverse(s.rankedWords))
}

// GetNextMove simply returns the top-ranked word
func (s *CharFrequencyStrategy) GetNextMove() string {
	if s.attempts == 0 && s.starter != "" {
		return s.starter
	}
	// else if s.attempts == 1 {
	// 	return "crony"
	// }
	return s.rankedWords[0].Key
}

func (s *CharFrequencyStrategy) SetMoveOutcome(row []game.GridCell) {
	// Update the internal state for the row
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
			s.answersIncorrectAll = append(s.answersIncorrectAll, cell.Letter)
		}
	}

	s.attempts++

	// Rebuild the histogram and ranking
	s.possibleAnswers = filterWordsList(s.possibleAnswers, row)
	s.buildHistogram()
	s.rankWords()
}

// GetSuggestions will get the best n suggestions given the current state
func (s *CharFrequencyStrategy) GetSuggestions(n int) PairList {
	if n >= len(s.rankedWords) {
		n = len(s.rankedWords) - 1
	}
	return s.rankedWords[0:n]
}

// NewCharFrequencyStrategy create a char frequency-based strategy given the word list and letters list
func NewCharFrequencyStrategy(wordLength int, letters []string, validAnswers []string, allAcceptedWords *[]string, starter string) Strategy {
	lettersMap := map[string]bool{}
	for _, l := range letters {
		lettersMap[l] = true
	}

	s := &CharFrequencyStrategy{
		starter:         starter,
		wordLength:      wordLength,
		letters:         lettersMap,
		possibleAnswers: validAnswers,
		allWords:        allAcceptedWords,
	}

	if starter == "" {
		// Initialise the histogram
		s.buildHistogram()
		// Rank the words
		s.rankWords()
	}

	return s
}

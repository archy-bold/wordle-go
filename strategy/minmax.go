package strategy

import (
	"math"

	"github.com/archy-bold/wordle-go/game"
)

type minMaxGuessRanking struct {
	maxScore     int
	averageScore int
	bestScore    int
	index        int
}

func compareRankings(a, b minMaxGuessRanking) bool {
	if a.maxScore < b.maxScore {
		return true
	}
	if a.maxScore > b.maxScore {
		return false
	}

	if a.averageScore < b.averageScore {
		return true
	}
	if a.averageScore > b.averageScore {
		return false
	}

	if a.bestScore < b.bestScore {
		return true
	}
	if a.bestScore > b.bestScore {
		return false
	}

	return a.index < b.index
}

type MinMaxStrategy struct {
	attempts         int
	starter          string
	wordLength       int
	possibleAnswers  []string
	allAcceptedWords *[]string
	bestGuessRanking minMaxGuessRanking
}

func (s *MinMaxStrategy) GetNextMove() string {
	if s.attempts == 0 && s.starter != "" {
		return s.starter
	}
	if len(s.possibleAnswers) == 1 {
		return s.possibleAnswers[0]
	}

	return (*s.allAcceptedWords)[s.bestGuessRanking.index]
}

func (s *MinMaxStrategy) SetMoveOutcome(row []game.GridCell) {
	// Simply filter down the possibleAnswers based on the row and build the ranked words list
	s.attempts++
	s.possibleAnswers = filterWordsList(s.possibleAnswers, row)

	// Re-rank words
	s.rankWords()
}

func (s *MinMaxStrategy) rankWords() {
	s.bestGuessRanking = minMaxGuessRanking{
		maxScore:     math.MaxInt32,
		averageScore: math.MaxInt32,
		bestScore:    math.MaxInt32,
		index:        0,
	}
	// s.rankedWords = make(PairList, len(*s.allAcceptedWords))

	for i, guess := range *s.allAcceptedWords {
		ranking := minMaxGuessRanking{
			maxScore:     0,
			averageScore: 0,
			bestScore:    0,
			index:        i,
		}

		for _, word := range s.possibleAnswers {
			if word != guess {
				hint := game.EvaluateGuess(guess, word)
				score := filterWordsListCount(s.possibleAnswers, hint)
				if score == 0 {
					score = len(s.possibleAnswers)
				}
				ranking.averageScore += score
				ranking.maxScore = max(score, ranking.maxScore)
				ranking.bestScore = min(score, ranking.bestScore)
			} else {
				ranking.bestScore = 0
			}

			if ranking.maxScore > s.bestGuessRanking.maxScore {
				break
			}
		}

		// fmt.Printf("%d ", i)

		if compareRankings(ranking, s.bestGuessRanking) {
			s.bestGuessRanking = ranking
		}
		// s.rankedWords[i] = Pair{guess, maxScore}
	}

	// Sort the ranked words
	// sort.Sort(s.rankedWords)

	// fmt.Printf("rankedWords %+v\n", s.rankedWords)
}

// GetSuggestions will get the best n suggestions given the current state
func (s *MinMaxStrategy) GetSuggestions(n int) PairList {
	if s.starter != "" && s.attempts == 0 {
		return PairList{{Key: s.starter, Value: 0}}
	}
	// if n >= len(s.rankedWords) {
	// 	n = len(s.rankedWords) - 1
	// }
	return PairList{{Key: (*s.allAcceptedWords)[s.bestGuessRanking.index], Value: s.bestGuessRanking.maxScore}}
}

// NewMinMaxStrategy create a MinMax-based strategy given the word lists
func NewMinMaxStrategy(wordLength int, possibleAnswers []string, allAcceptedWords *[]string, starter string) Strategy {
	s := &MinMaxStrategy{
		starter:          starter,
		wordLength:       wordLength,
		possibleAnswers:  possibleAnswers,
		allAcceptedWords: allAcceptedWords,
	}

	// Rank the words the first time
	if starter == "" {
		s.rankWords()
	}

	return s
}

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/archy-bold/wordle-go/game"
	"github.com/archy-bold/wordle-go/strategy"
)

const (
	NUM_LETTERS  = 5
	NUM_ATTEMPTS = 6
)

var ErrAllStartersFlagInvalid = errors.New("all-starters flag must be one of: valid, answers")
var ErrStrategyFlagInvalid = errors.New("strategy flag must be one of: minmax, charfreq")

var letters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var startDate = time.Date(2021, time.June, 19, 0, 0, 0, 0, time.UTC)
var validWords = []string{}
var allAcceptedWords = []string{}

// var dictionary = map[string]int{}

type wordAnalysisResult struct {
	word        string
	successes   int
	sumTries    int
	avgTries    float64
	longestGame int
}

func main() {
	wordPtr := flag.String("word", "", "The game's answer")
	cheatPtr := flag.Bool("cheat", false, "Whether to run the solver mode")
	autoPtr := flag.Bool("auto", false, "Play the game automatically")
	strategyPtr := flag.String("strategy", "charfreq", "Choose which strategy to use. One of 'minmax' (slow, effective) or 'charfreq' (fast, less effective).")
	randomPtr := flag.Bool("random", false, "Choose a random word, if none specified. Otherwise gets daily word")
	datePtr := flag.String("date", "", "If specified, will choose the word for this day")
	starterPtr := flag.String("starter", "", "The starter word to use in strategies")
	allPtr := flag.Bool("all", false, "Play all permutations of the answers")
	allStartersPtr := flag.String("all-starters", "", "Play all permutations of the answers, with all permutations of the chosen starter list. Starter options are: valid (12972 iterations), answers (2315 iterations)")
	flag.Parse()

	if *allStartersPtr != "" && *allStartersPtr != "valid" && *allStartersPtr != "answers" {
		fmt.Println(ErrAllStartersFlagInvalid.Error())
		return
	}
	if *strategyPtr != "" && *strategyPtr != "minmax" && *strategyPtr != "charfreq" {
		fmt.Println(ErrStrategyFlagInvalid.Error())
		return
	}

	auto := *autoPtr

	// Read the valid words
	var err error
	fmt.Print("Reading solutions... ")
	err = readWordList(&validWords, "words/5/solutions.txt")
	fmt.Printf("found %d\n", len(validWords))
	check(err)

	// Read the valid guesses
	fmt.Print("Reading valid guesses... ")
	err = readWordList(&allAcceptedWords, "words/5/guesses.txt")
	// Sort the valid guesses as we will be searching that array often
	sort.Strings(allAcceptedWords)
	fmt.Printf("found %d\n", len(allAcceptedWords))
	check(err)

	reader := bufio.NewReader(os.Stdin)

	var strat strategy.Strategy
	if auto {
		strat = getStrategy(*strategyPtr, *starterPtr)
	}

	// Play out all permutations
	if *allPtr {
		sumTries := 0
		numSuccesses := 0
		for i, answer := range validWords {
			strat = getStrategy(*strategyPtr, *starterPtr)
			g := game.CreateGame(answer, NUM_ATTEMPTS, &allAcceptedWords, i+1)

			for {
				word := strat.GetNextMove()
				success, _ := g.Play(word)
				strat.SetMoveOutcome(g.GetLastPlay())

				if success {
					score, of := g.GetScore()
					fmt.Printf("%s in %d/%d\n", answer, score, of)
					numSuccesses++
					sumTries += score

					break
				} else if g.HasEnded() {
					fmt.Printf("%s failed\n", answer)
					break
				}

			}
		}

		fmt.Printf("Completed %d/%d\n", numSuccesses, len(validWords))
		fmt.Printf("On average %f\n", float64(sumTries)/float64(numSuccesses))

		return
	} else if *allStartersPtr != "" {
		var results []wordAnalysisResult
		var words []string
		if *allStartersPtr == "answers" {
			words = validWords
		} else {
			words = allAcceptedWords
		}
		results = make([]wordAnalysisResult, len(words))

		for i, starter := range words {

			var wg sync.WaitGroup
			wg.Add(len(validWords))

			fmt.Printf("%4d %s ", i, starter)

			result := wordAnalysisResult{
				word:        starter,
				successes:   0,
				sumTries:    0,
				avgTries:    0,
				longestGame: 0,
			}
			var resMutex sync.Mutex

			for _, answer := range validWords {
				go func(result *wordAnalysisResult, resMutex *sync.Mutex, answer string) {
					strat := strategy.NewCharFrequencyStrategy(NUM_LETTERS, letters, validWords, &allAcceptedWords, starter)
					g := game.CreateGame(answer, 15, &allAcceptedWords, i+1)

					for {
						word := strat.GetNextMove()
						success, _ := g.Play(word)
						strat.SetMoveOutcome(g.GetLastPlay())

						if success {
							resMutex.Lock()
							score, _ := g.GetScore()
							result.sumTries += score
							if score <= NUM_ATTEMPTS {
								result.successes++
							}
							if result.longestGame < score {
								result.longestGame = score
							}
							resMutex.Unlock()
							wg.Done()
							return
						} else if g.HasEnded() {
							wg.Done()
							return
						}
					}
				}(&result, &resMutex, answer)
			}

			wg.Wait()

			result.avgTries = float64(result.sumTries) / float64(result.successes)
			fmt.Printf("%4d %d %f\n", result.successes, result.longestGame, result.avgTries)

			results[i] = result
		}

		// Sort by successes first
		sort.Slice(results, func(i, j int) bool {
			return results[i].successes > results[j].successes
		})
		fmt.Println("\nRanking: Num Successes")
		for i := 0; i < 50; i++ {
			res := results[i]
			fmt.Printf("%d. %s %d %d %f\n", i+1, res.word, res.successes, res.longestGame, res.avgTries)
		}
		// Then average tries
		sort.Slice(results, func(i, j int) bool {
			return results[i].avgTries < results[j].avgTries
		})
		fmt.Println("\nRanking: Average Tries")
		for i := 0; i < 50; i++ {
			res := results[i]
			fmt.Printf("%d. %s %d %d %f\n", i+1, res.word, res.successes, res.longestGame, res.avgTries)
		}
		// Finally the longest game
		sort.Slice(results, func(i, j int) bool {
			return results[i].longestGame > results[j].longestGame
		})
		fmt.Println("\nRanking: Longest Game")
		fmt.Printf(" %s %d %d %f\n", results[0].word, results[0].successes, results[0].longestGame, results[0].avgTries)

		return
	}

	// Cheat mode
	if *cheatPtr {
		strat = getStrategy(*strategyPtr, *starterPtr)
		ug := game.CreateUnknownGame(NUM_LETTERS, NUM_ATTEMPTS)
		for {

			// Print the top 10 answers
			fmt.Println("Top answers:")
			suggestions := strat.GetSuggestions(10)
			for i, suggestion := range suggestions {
				if suggestion.Key != "" {
					fmt.Printf("  %d: %s (%d)\n", i+1, suggestion.Key, suggestion.Value)
				}
			}

			// Read the entered word from stdin
			// TODO handle errors such as wrong sized word, wrong pattern for response
			fmt.Print("Enter number of entered word, or word itself: ")
			word, _ := reader.ReadString('\n')
			word = strings.TrimSpace(word)
			if idx, err := strconv.Atoi(word); err == nil && idx <= len(suggestions) {
				word = suggestions[idx-1].Key
			}
			wordParts := strings.Split(word, "")

			fmt.Print("Enter the result, where x is incorrect, o is wrong position, y is correct\n eg yxxox: ")
			input, _ := reader.ReadString('\n')
			parts := strings.Split(strings.TrimSpace((input)), "")
			row := make([]game.GridCell, NUM_LETTERS)
			for i, chr := range parts {
				if chr == "x" {
					row[i] = game.GridCell{
						Letter: wordParts[i],
						Status: game.STATUS_WRONG,
					}
				} else if chr == "y" {
					row[i] = game.GridCell{
						Letter: wordParts[i],
						Status: game.STATUS_CORRECT,
					}
				} else if chr == "o" {
					row[i] = game.GridCell{
						Letter: wordParts[i],
						Status: game.STATUS_INCORRECT,
					}
				}
			}

			// Update the game grid and strategy
			strat.SetMoveOutcome(row)
			complete, _ := ug.(*game.UnknownGame).AddResult(row)

			if complete {
				score, _ := ug.GetScore()
				fmt.Printf("Hooray! (%d/%d)\n", score, NUM_ATTEMPTS)
				return
			}
		}
	}

	// If no answer given in the word flag, choose
	answer := *wordPtr

	var gameNum int
	if answer == "" {
		if *randomPtr {
			rand.Seed(time.Now().Unix())
			gameNum = rand.Intn(len(validWords))
		} else {
			// Go by date
			var dt time.Time
			if *datePtr != "" {
				dt, err = time.Parse("2006-01-02", *datePtr)
				check(err)
				// TODO check the date isn't before the start date or after the end date
			} else {
				dt = time.Now().UTC()
				year, month, day := dt.Date()
				dt = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
			}
			gameNum = int(dt.Sub(startDate).Hours() / 24)
		}
		answer = validWords[gameNum]
	}
	g := game.CreateGame(answer, NUM_ATTEMPTS, &allAcceptedWords, gameNum)

	for {
		// Play based on whether a strategy is provided
		var word string
		if strat != nil {
			word = strat.GetNextMove()
		} else {
			fmt.Print("Enter your guess: ")
			input, _ := reader.ReadString('\n')
			word = strings.TrimSpace(input)
		}

		// TODO handle error here
		success, _ := g.Play(word)

		if strat != nil {
			strat.SetMoveOutcome(g.GetLastPlay())
		}

		if strat == nil || success || g.HasEnded() {
			fmt.Println(g.OutputForConsole())
		}

		if success {
			score, of := g.GetScore()
			fmt.Println(g.OutputToShare())
			fmt.Printf("Great work! %d/%d\n", score, of)
			return
		} else if g.HasEnded() {
			fmt.Println(g.OutputToShare())
			fmt.Printf("Better luck next time! The word was '%s'. X/%d\n", answer, NUM_ATTEMPTS)
			return
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getStrategy(strat string, starter string) strategy.Strategy {
	if strat == "minmax" {
		return strategy.NewMinMaxStrategy(NUM_LETTERS, validWords, &allAcceptedWords, starter)
	}
	return strategy.NewCharFrequencyStrategy(NUM_LETTERS, letters, validWords, &allAcceptedWords, starter)
}

func readWordList(arr *[]string, fname string) error {
	data, err := Asset(fname)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	i := 0
	for scanner.Scan() {
		word := scanner.Text()
		// dictionary[word] = i
		*arr = append(*arr, word)
		i++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

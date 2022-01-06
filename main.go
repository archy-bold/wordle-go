package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/archy-bold/wordle-go/game"
	"github.com/archy-bold/wordle-go/strategy"
)

const (
	NUM_LETTERS  = 5
	NUM_ATTEMPTS = 6
)

var letters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var startDate = time.Date(2021, time.June, 19, 0, 0, 0, 0, time.UTC)
var validWords = []string{}
var allAcceptedWords = []string{}

// var dictionary = map[string]int{}

func main() {
	wordPtr := flag.String("word", "", "The game's answer")
	cheatPtr := flag.Bool("cheat", false, "Whether to run the solver mode")
	autoPtr := flag.Bool("auto", false, "Play the game automatically")
	randomPtr := flag.Bool("random", false, "Choose a random word, if none specified. Otherwise gets daily word")
	datePtr := flag.String("date", "", "If specified, will choose the word for this day")
	starterPtr := flag.String("starter", "", "The starter word to use in strategies")
	allPtr := flag.Bool("all", false, "Play all permutations")
	flag.Parse()

	auto := *autoPtr

	// Read the valid words
	var err error
	fmt.Print("Reading solutions... ")
	err = readWordList(&validWords, "./words/5/solutions.txt")
	fmt.Printf("found %d\n", len(validWords))
	check(err)

	// Read the valid guesses
	fmt.Print("Reading valid guesses... ")
	err = readWordList(&allAcceptedWords, "./words/5/guesses.txt")
	// Sort the valid guesses as we will be searching that array often
	sort.Strings(allAcceptedWords)
	fmt.Printf("found %d\n", len(allAcceptedWords))
	check(err)

	reader := bufio.NewReader(os.Stdin)

	var strat strategy.Strategy
	if auto {
		strat = strategy.NewCharFrequencyStrategy(NUM_LETTERS, letters, &validWords, *starterPtr)
	}

	// Play out all permutations
	if *allPtr {
		sumTries := 0
		numSuccesses := 0
		for _, answer := range validWords {
			strat = strategy.NewCharFrequencyStrategy(NUM_LETTERS, letters, &validWords, *starterPtr)
			g := game.CreateGame(answer, NUM_ATTEMPTS, &allAcceptedWords)

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
	}

	// Cheat mode
	if *cheatPtr {
		strat = strategy.NewCharFrequencyStrategy(NUM_LETTERS, letters, &validWords, *starterPtr)
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

			fmt.Print("Enter the result, where x is incorrect, o is wrong position, y is correct eg yxxox: ")
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

	if answer == "" {
		var pos int
		if *randomPtr {
			rand.Seed(time.Now().Unix())
			pos = rand.Intn(len(validWords))
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
			pos = int(dt.Sub(startDate).Hours() / 24)
		}
		answer = validWords[pos]
	}
	g := game.CreateGame(answer, NUM_ATTEMPTS, &allAcceptedWords)

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

func readWordList(arr *[]string, fname string) error {
	file, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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

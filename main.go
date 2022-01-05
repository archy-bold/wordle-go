package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
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
var validWords = []string{}

// var dictionary = map[string]int{}

func main() {
	wordPtr := flag.String("word", "", "The game's answer")
	cheatPtr := flag.Bool("cheat", false, "Whether to run the solver mode")
	autoPtr := flag.Bool("auto", false, "Play the game automatically")
	starterPtr := flag.String("starter", "", "The starter word to use in strategies")
	allPtr := flag.Bool("all", false, "Play all permutations")
	flag.Parse()

	auto := *autoPtr

	// Read the valid words
	fmt.Print("Reading words... ")
	err := readValidWords()
	fmt.Printf("found %d\n", len(validWords))
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
			g := game.CreateGame(answer, NUM_ATTEMPTS)

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
		rand.Seed(time.Now().Unix())
		answer = validWords[rand.Intn(len(validWords))]
	}
	g := game.CreateGame(answer, NUM_ATTEMPTS)

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

func readValidWords() error {
	file, err := os.Open("./solutions.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		word := scanner.Text()
		// dictionary[word] = i
		validWords = append(validWords, word)
		i++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

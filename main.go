package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	NUM_LETTERS  = 5
	NUM_ATTEMPTS = 6

	COLOUR_RESET  = "\033[0m"
	COLOUR_GREEN  = "\033[32m"
	COLOUR_YELLOW = "\033[33m"
)

type HistogramEntry struct {
	Occurences            int
	OccurrencesInPosition []int
}

var letters = map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true, "f": true, "g": true, "h": true, "i": true, "j": true, "k": true, "l": true, "m": true, "n": true, "o": true, "p": true, "q": true, "r": true, "s": true, "t": true, "u": true, "v": true, "w": true, "x": true, "y": true, "z": true}
var validWords = []string{}
var dictionary = map[string]int{}
var histogram = map[string]HistogramEntry{}
var rankedWords PairList
var answersCorrect []string
var answersIncorrect [][]string
var answersIncorrectAll []string
var board []string

func main() {
	// Read the valid words
	fmt.Println("Reading words...")
	err := readValidWords()
	check(err)

	answersCorrect = make([]string, NUM_LETTERS)
	answersIncorrect = make([][]string, NUM_LETTERS)

	reader := bufio.NewReader(os.Stdin)
	for {
		// Generate histogram
		fmt.Println("Building histogram...")
		buildHistogram()

		// Rank words based on frequency
		fmt.Println("Ranking words...")
		rankWords()

		// Print the top 10 answers
		fmt.Println("Top answers:")
		ln := 10
		if len(rankedWords) < ln {
			ln = len(rankedWords)
		}
		for i := 0; i < ln; i++ {
			rank := rankedWords[i]
			if rank.Key != "" {
				fmt.Printf("  %d: %s (%d)\n", i+1, rank.Key, rank.Value)
			}
		}

		// Read the entered word from stdin
		answersIncorrectAll = make([]string, 0)
		// TODO handle errors such as wrong sized word, wrong pattern for response
		fmt.Print("Enter number of entered word, or word itself: ")
		word, _ := reader.ReadString('\n')
		word = strings.TrimSpace(word)
		if idx, err := strconv.Atoi(word); err == nil && idx <= len(rankedWords) {
			word = rankedWords[idx-1].Key
		}
		wordParts := strings.Split(word, "")

		fmt.Print("Enter the result, where x is incorrect, o is wrong position, y is correct eg yxxox: ")
		input, _ := reader.ReadString('\n')
		parts := strings.Split(strings.TrimSpace((input)), "")
		boardRow := ""
		numCorrect := 0
		rejected := make([]bool, NUM_LETTERS)
		for i, chr := range parts {
			if chr == "x" {
				// If this letter has shown up before but not rejected, don't eliminate
				shouldEliminate := true
				for j := 0; j < i; j++ {
					if chr == parts[j] && !rejected[j] {
						shouldEliminate = false
					}
				}
				if shouldEliminate {
					letters[wordParts[i]] = false
				}
				rejected[i] = true
			} else if chr == "y" {
				boardRow += COLOUR_GREEN
				answersCorrect[i] = wordParts[i]
				numCorrect++
			} else if chr == "o" {
				boardRow += COLOUR_YELLOW
				answersIncorrect[i] = append(answersIncorrect[i], wordParts[i])
				answersIncorrectAll = append(answersIncorrectAll, wordParts[i])
			}
			boardRow += wordParts[i] + COLOUR_RESET
		}
		board = append(board, boardRow)

		outputBoard()

		if numCorrect == NUM_LETTERS {
			fmt.Printf("Hooray! (%d/%d)\n", len(board), NUM_ATTEMPTS)
			return
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
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
		dictionary[word] = i
		validWords = append(validWords, word)
		i++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func buildHistogram() {
	histogram = make(map[string]HistogramEntry, len(letters))
	for l := range letters {
		histogram[l] = HistogramEntry{0, make([]int, NUM_LETTERS)}
	}

	// Loop through each word and check which unique letters are in the word
	for _, word := range validWords {
		chrs := strings.Split(word, "")

		// Loop through each char and update the histogram
		checkedChars := map[string]bool{}
		for i, chr := range chrs {
			// Ignore removed letters
			if !letters[chr] {
				continue
			}
			// If we've not already processed this letter
			if _, ok := checkedChars[chr]; !ok {
				checkedChars[chr] = true
				if entry, ok2 := histogram[chr]; ok2 {
					entry.Occurences++
					entry.OccurrencesInPosition[i]++
					histogram[chr] = entry
				}
			}
		}
	}
}

func rankWords() {
	rankedWords = make(PairList, len(validWords))

word:
	for _, word := range validWords {
		chrs := strings.Split(word, "")
		// First set the score based on the letters that exist
		// TODO score based on letter position too
		checkedChars := map[string]bool{}
		score := 0

		// Check if any of the incorrect answers don't appear in the word
		if len(answersIncorrectAll) > 0 {
			diff := difference(answersIncorrectAll, chrs)
			if len(diff) > 0 {
				rankedWords = append(rankedWords, Pair{word, 0})
				continue word
			}
		}

		for i, chr := range chrs {
			// If this is an eliminated letter, score down
			if !letters[chr] {
				rankedWords = append(rankedWords, Pair{word, 0})
				continue word
			}

			// If there is an answer in this position, we can disregard words that don't have that letter in that position
			if answersCorrect[i] != "" && answersCorrect[i] != chr {
				rankedWords = append(rankedWords, Pair{word, 0})
				continue word
			}

			// Also check if there's an incorrect answer in this position
			if len(answersIncorrect[i]) > 0 {
				for _, ai := range answersIncorrect[i] {
					if ai == chr {
						rankedWords = append(rankedWords, Pair{word, 0})
						continue word
					}
				}
			}

			if _, ok := checkedChars[chr]; !ok {
				// Score based on occurences and occurences in the position
				scoreToAdd := histogram[chr].Occurences + histogram[chr].OccurrencesInPosition[i]
				// Increase score for incorrectly placed letters
				for _, aiChr := range answersIncorrectAll {
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
		rankedWords = append(rankedWords, Pair{word, score})
	}

	// Sort the
	sort.Sort(sort.Reverse(rankedWords))
}

func outputBoard() {
	fmt.Println("")
	fmt.Println(strings.Repeat("-", NUM_LETTERS+2))
	for _, row := range board {
		fmt.Printf("|%s|\n", row)
	}
	fmt.Println(strings.Repeat("-", NUM_LETTERS+2))
	fmt.Println("")
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const NUM_LETTERS = 5

var letters = map[string]bool{"a": true, "b": true, "c": true, "d": true, "e": true, "f": true, "g": true, "h": true, "i": true, "j": true, "k": true, "l": true, "m": true, "n": true, "o": true, "p": true, "q": true, "r": true, "s": true, "t": true, "u": true, "v": true, "w": true, "x": true, "y": true, "z": true}
var validWords = []string{}
var dictionary = map[string]int{}
var histogram = map[string]int{}
var rankedWords PairList
var answers []string

func main() {
	// Read the valid words
	fmt.Println("Reading words...")
	err := readValidWords()
	check(err)

	answers = make([]string, NUM_LETTERS)

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
		for i := 0; i < 10; i++ {
			rank := rankedWords[i]
			fmt.Printf("  %s: %d\n", rank.Key, rank.Value)
		}

		// Read from stdin letters to eliminate
		fmt.Print("Enter letters to eliminate eg afb: ")
		input, _ := reader.ReadString('\n')
		parts := strings.Split(input, "")
		for _, chr := range parts {
			// Ignore anything not in the map
			if _, ok := letters[chr]; ok {
				letters[chr] = false
			}
		}

		// Read from stdin correct letters
		fmt.Print("Enter correctly positioned letters eg --a-b: ")
		input, _ = reader.ReadString('\n')
		parts = strings.Split(input, "")
		for i, chr := range parts {
			// Ignore anything not in the map
			if _, ok := letters[chr]; ok && i < NUM_LETTERS {
				answers[i] = chr
			}
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
	histogram = make(map[string]int, len(letters))
	for l := range letters {
		histogram[l] = 0
	}

	// Loop through each word and check which unique letters are in the word
	for _, word := range validWords {
		chrs := strings.Split(word, "")

		// Loop through each char and update the histogram
		checkedChars := map[string]bool{}
		for _, chr := range chrs {
			// Ignore removed letters
			if !letters[chr] {
				continue
			}
			// If we've not already processed this letter
			if _, ok := checkedChars[chr]; !ok {
				histogram[chr]++
				checkedChars[chr] = true
			}
		}
	}
}

func rankWords() {
	rankedWords = make(PairList, len(validWords))

	for _, word := range validWords {
		chrs := strings.Split(word, "")
		// First set the score based on the letters that exist
		// TODO score based on letter position too
		checkedChars := map[string]bool{}
		score := 0
		for i, chr := range chrs {
			// If there is an answer in this position, we can disregard words that don't have that letter in that position
			if answers[i] != "" && answers[i] != chr {
				score = 0
				break
			}

			if _, ok := checkedChars[chr]; !ok {
				score += histogram[chr]
				checkedChars[chr] = true
			}
		}

		// Add to the ranked list
		rankedWords = append(rankedWords, Pair{word, score})
	}

	// Sort the
	sort.Sort(sort.Reverse(rankedWords))
}

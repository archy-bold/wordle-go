package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var letters = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
var validWords = []string{}
var dictionary = map[string]int{}
var histogram = map[string]int{}
var rankedWords PairList

func main() {
	// Read the valid words
	fmt.Println("Reading words...")
	err := readValidWords()
	check(err)

	// Generate histogram
	fmt.Println("Building histogram...")
	buildHistogram()

	// Rank words based on frequency
	fmt.Println("Ranking words...")
	rankWords()

	fmt.Printf("%+v\n", histogram)
	for i := 0; i < 10; i++ {
		rank := rankedWords[i]
		fmt.Printf("%s: %d\n", rank.Key, rank.Value)
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
	for _, l := range letters {
		histogram[l] = 0
	}

	// Loop through each word and check which unique letters are in the word
	for _, word := range validWords {
		chrs := strings.Split(word, "")

		// Loop through each char and update the histogram
		checkedChars := map[string]bool{}
		for _, chr := range chrs {
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
		checkedChars := map[string]bool{}
		score := 0
		for _, chr := range chrs {
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

# Wordle Command Line

![wordle command line tool running](https://media.giphy.com/media/35BYFJGLLiRze1NVkZ/giphy.gif)

A Wordle implementation for the command line, written in go.

## Get It

Download the latest version for your system from the [releases](https://github.com/archy-bold/wordle-go/releases) page.

Or clone the repo and run from there. You may need to re-build for your system first.

```bash
git clone git@github.com:archy-bold/wordle-go.git
cd wordle-go
./wordle
```

## Build

```bash
go build -o wordle
```

## Install

To install to your `$GOPATH` bin folder, run the following;

```bash
go install && mv $GOPATH/bin/wordle-go $GOPATH/bin/wordle
wordle
```

## Run

Run the command with no arguments to play a game with today's word.

```bash
./wordle
```

### Options

- `-cheat` - Runs in solve mode to work out an existing wordle. Follow the instructions to enter your results and receive suggested words to play.
- `-random` - Choose a random word instead of today's
- `-auto` - Automatically completes the puzzle
- `-strategy=[minmax|charfreq]` - Choose a strategy to solve games.
  - `charfreq` - *[default]* The character frequency strategy will analyse the character frequency of the answers to determine the best word to use. It plays in hard mode. This is much quicker, but far less effective at solving games and may fail on some words, depending on the starter word. 
  - `minmax` - The MinMax strategy aims to eliminate as many possible solutions in as few turns as possible. This strategy is very slow as it checks each possible word against each possible answer. It's recommended using this with a `starter` specified.
- `-starter=[word]` - Specify the starter word for strategies
- `-date=[2021-12-31]` - Set the winning word from a specific date
- `-word=[answer]` - Set the winning word with this argument.
- `-all` - Runs the auto-solver through every permutation, giving results when complete.
- `-all-staters=[answers|valid]` - Run the auto-solver through all game permutations in turn with each starter. Starters list comprises either the answers list (2315 iterations) or the valid words list (12972 words).

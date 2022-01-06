package game

import (
	"fmt"
	"sort"
	"strings"
)

const (
	NUM_LETTERS  = 5
	NUM_ATTEMPTS = 6

	COLOUR_RESET  = "\033[0m"
	COLOUR_GREEN  = "\033[32m"
	COLOUR_YELLOW = "\033[33m"
	COLOUR_GREY   = "\u001b[30;1m"
)

func outputGridForConsole(grid [][]GridCell, length int, numSpaces int) string {
	spacing := strings.Repeat(" ", numSpaces)
	str := "\n" + spacing + strings.Repeat("-", length+2) + "\n"
	for _, row := range grid {
		if len(row) == 0 {
			break
		}

		str += spacing + "|"
		for _, cell := range row {
			switch cell.Status {
			case STATUS_CORRECT:
				str += COLOUR_GREEN
			case STATUS_INCORRECT:
				str += COLOUR_YELLOW
			}
			str += strings.ToUpper(cell.Letter) + COLOUR_RESET
		}
		str += "|\n"
	}
	str += spacing + strings.Repeat("-", length+2) + "\n"

	return str
}

func outputGridToShare(grid [][]GridCell, score string, of int) string {
	// TODO output game number here
	str := fmt.Sprintf("Wordle %s/%d\n\n", score, of)
	for _, row := range grid {
		if len(row) == 0 {
			break
		}

		for _, cell := range row {
			switch cell.Status {
			case STATUS_CORRECT:
				str += "🟩"
			case STATUS_INCORRECT:
				str += "🟨"
			case STATUS_WRONG:
				str += "⬜"
			}
		}
		str += "\n"
	}
	str += "\n"

	return str
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func stringInSortedSlice(a string, list *[]string) bool {
	i := sort.SearchStrings(*list, a)
	return i < len(*list) && (*list)[i] == a

}

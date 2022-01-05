package game

const (
	NUM_LETTERS  = 5
	NUM_ATTEMPTS = 6

	COLOUR_RESET  = "\033[0m"
	COLOUR_GREEN  = "\033[32m"
	COLOUR_YELLOW = "\033[33m"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

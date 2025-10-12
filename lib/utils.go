package lib

import "strings"

func CleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	text = strings.ToLower(text)
	words := strings.Fields(text)

	return words
}

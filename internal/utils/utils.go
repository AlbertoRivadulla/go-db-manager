package utils

import (
	"strings"
)

func WrapText(text string, allowedWidth int) (string, int) {
	if len(text) <= allowedWidth {
		return text, len(text)
	}

	width := 0

	var lines []string
	words := strings.Fields(text)
	currLine := words[0]

	for _, word := range words[1:] {
		if len(currLine) + 1 + len(word) <= allowedWidth {
			currLine += " " + word
		} else {
			lines = append(lines, currLine)

			if len(currLine) > width {
				width = len(currLine)
			}

			currLine = word
		}
	}

	lines = append(lines, currLine)
	return strings.Join(lines, "\n"), width
}

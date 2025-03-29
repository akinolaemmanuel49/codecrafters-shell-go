package parser

import "strings"

func ParseInput(input string) []string {
	var tokens []string
	var current strings.Builder
	var inQuotes rune = 0 // Track if inside quotes (' or ")

	for _, ch := range strings.TrimSpace(input) {
		switch ch {
		case ' ', '\t': // Space or tab
			if inQuotes != 0 {
				current.WriteRune(ch) // Keep spaces inside quotes
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case '\'', '"': // Handle quotes
			if inQuotes == ch {
				inQuotes = 0 // Closing quote
			} else if inQuotes == 0 {
				inQuotes = ch // Opening quote
			} else {
				current.WriteRune(ch) // Nested quotes (e.g., "it's")
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String()) // Add last token
	}

	return tokens
}

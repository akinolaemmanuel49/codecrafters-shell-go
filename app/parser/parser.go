package parser

import (
	"errors"
	"strings"
)

// ParseInput parses a shell-like command string into arguments, handling quotes and escapes.
func ParseInput(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var inQuotes rune = 0 // Track quote type: ' or "
	escaped := false      // Track if previous character was a backslash

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		// Handle escape sequences
		if escaped {
			current.WriteRune(ch) // Keep the escaped character
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			escaped = true // Mark next character as escaped
		case ' ', '\t':
			if inQuotes != 0 {
				current.WriteRune(ch) // Keep spaces inside quotes
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String()) // Add completed token
				current.Reset()
			}
		case '\'', '"':
			if inQuotes == ch {
				inQuotes = 0 // Closing quote
			} else if inQuotes == 0 {
				inQuotes = ch // Opening quote
			} else {
				current.WriteRune(ch) // Keep the quote if it's inside different quotes
			}
		default:
			current.WriteRune(ch)
		}
	}

	if escaped {
		return nil, errors.New("syntax error: trailing backslash")
	}

	if inQuotes != 0 {
		return nil, errors.New("syntax error: unmatched quote")
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String()) // Add last token
	}

	return tokens, nil
}

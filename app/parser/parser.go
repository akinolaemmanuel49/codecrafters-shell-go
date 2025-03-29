package parser

import (
	"errors"
	"strings"
)

// ParseInput parses a shell-like command string into arguments, handling quotes and escapes.
func ParseInput(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var inQuotes rune = 0 // Tracks the active quote type (' or ")
	escaped := false      // Tracks if the last character was '\'

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		if escaped {
			// Handle escaped characters correctly
			if inQuotes != 0 { // Only expand \n inside quotes
				switch ch {
				case 'n':
					current.WriteRune('\n')
				case '\\', '\'', '"', ' ':
					current.WriteRune(ch)
				default:
					current.WriteRune('\\') // Keep unknown backslash escapes
					current.WriteRune(ch)
				}
			} else {
				// Keep `\n` as-is outside of quotes
				current.WriteRune('\\')
				current.WriteRune(ch)
			}
			escaped = false
			continue
		}

		switch ch {
		case '\\':
			escaped = true
		case ' ', '\t':
			if inQuotes != 0 {
				current.WriteRune(ch) // Keep spaces inside quotes
			} else if current.Len() > 0 {
				tokens = append(tokens, current.String()) // End of token
				current.Reset()
			}
		case '\'', '"':
			if inQuotes == ch {
				inQuotes = 0 // Closing quote
			} else if inQuotes == 0 {
				inQuotes = ch // Opening quote
			} else {
				current.WriteRune(ch) // Keep quotes inside different quote types
			}
		default:
			current.WriteRune(ch)
		}
	}

	// Handle trailing errors
	if escaped {
		return nil, errors.New("syntax error: trailing backslash")
	}
	if inQuotes != 0 {
		return nil, errors.New("syntax error: unmatched quote")
	}

	// Add the last token if present
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

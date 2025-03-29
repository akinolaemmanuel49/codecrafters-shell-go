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
	escaped := false      // Track escape sequences

	for _, ch := range strings.TrimSpace(input) {
		if escaped {
			// If previous char was '\', handle escaping
			if ch == ' ' {
				current.WriteRune(' ') // Convert `\ ` to a space
			} else if inQuotes == '"' && (ch == '"' || ch == '$' || ch == '`' || ch == '\\') {
				current.WriteRune(ch) // Escape valid double-quote characters
			} else {
				current.WriteRune('\\') // Preserve the backslash for other cases
				current.WriteRune(ch)
			}
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
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case '\'', '"':
			if inQuotes == ch {
				inQuotes = 0 // Closing quote
			} else if inQuotes == 0 {
				inQuotes = ch // Opening quote
			} else {
				current.WriteRune(ch) // Nested quotes allowed
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
		tokens = append(tokens, current.String())
	}

	return tokens, nil
}

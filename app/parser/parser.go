package parser

import (
	"errors"
	"strings"
)

// ParseInput parses a shell-like command string into arguments, handling quotes and backslashes.
func ParseInput(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var inQuotes rune = 0 // Track quote type: ' or "
	escaped := false      // Track if previous character was a backslash

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		// Handle escape sequences
		if escaped {
			if ch == '\n' {
				escaped = false // Ignore \newline (line continuation)
				continue
			}

			// Inside double quotes, only certain characters are escaped
			if inQuotes == '"' && (ch == '"' || ch == '$' || ch == '`' || ch == '\\') {
				current.WriteRune(ch)
			} else if inQuotes == 0 { // Outside quotes, backslash escapes any character
				current.WriteRune(ch)
			} else {
				current.WriteRune('\\') // Inside single quotes, backslash is preserved
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
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case '\'', '"':
			if inQuotes == ch {
				inQuotes = 0 // Closing quote
			} else if inQuotes == 0 {
				inQuotes = ch // Opening quote
			} else {
				current.WriteRune(ch) // Nested quote inside different type
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

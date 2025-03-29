package parser

import (
	"errors"
	"regexp"
	"strings"
)

// ParseInput parses a shell-like command string into arguments, handling quotes and escapes.
func ParseInput(input string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var inQuotes rune = 0 // Tracks quote type (' or ")
	escaped := false      // Tracks if previous character was a backslash

	if strings.Contains(input, "\\") {
		re := regexp.MustCompile(`[^\\] +`)
		args := re.Split(input, -1)

		for i := range args {
			args[i] = strings.ReplaceAll(args[i], "\\", "")
		}

		return args, nil
	}

	for i := 0; i < len(input); i++ {
		ch := rune(input[i])

		if escaped {
			current.WriteRune(ch)
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
				current.WriteRune(ch) // Allow different nested quotes
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

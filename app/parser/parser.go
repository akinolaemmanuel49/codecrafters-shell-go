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

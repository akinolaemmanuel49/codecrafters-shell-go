package parser

import (
	"regexp"
	"strings"
)

// ParseInput extracts arguments, handling quotes and backslashes correctly.
func ParseInput(input string) []string {
	var result []string

	// Match quoted strings or unquoted words
	re := regexp.MustCompile(`"([^"\\]*(\\.[^"\\]*)*)"|'([^'\\]*(\\.[^'\\]*)*)'|\\?(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	for _, match := range matches {
		if match[1] != "" { // Double-quoted string
			result = append(result, strings.ReplaceAll(match[1], `\"`, `"`))
		} else if match[3] != "" { // Single-quoted string
			result = append(result, match[3])
		} else if match[5] != "" { // Unquoted word with optional backslash prefix
			result = append(result, strings.ReplaceAll(match[5], "\\", ""))
		}
	}

	return result
}

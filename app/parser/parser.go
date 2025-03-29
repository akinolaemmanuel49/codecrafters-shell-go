package parser

import (
	"strings"
)

func ParseInput(input string) []string {
	// Trim any carriage returns or newlines
	_input := strings.Trim(input, "\r\n")
	var tokens []string
	var currentToken string
	var inWord bool = false

	// Process the input until we've consumed it all
	for len(_input) > 0 {
		// Skip leading whitespace when not in a word
		if !inWord {
			_input = strings.TrimLeft(_input, " \t")
			if len(_input) == 0 {
				break
			}
		}

		// Check if we start with a quote
		if _input[0] == '"' || _input[0] == '\'' {
			inWord = true
			// Get the quote character
			quote := _input[0]
			// Remove the opening quote
			_input = _input[1:]

			// Find the matching closing quote, handling escaped quotes
			for len(_input) > 0 {
				// If we find a backslash inside a quote, check if it's escaping something
				if len(_input) > 1 && _input[0] == '\\' {
					currentToken += string(_input[1])
					_input = _input[2:]
				} else if _input[0] == quote {
					// End of quoted section
					_input = _input[1:]
					break
				} else {
					// Add the character to our token
					currentToken += string(_input[0])
					_input = _input[1:]
				}
			}
		} else if _input[0] == '\\' && len(_input) > 1 {
			// Handle backslash escape outside quotes
			inWord = true
			// Preserve the literal value of the next character, including space
			currentToken += string(_input[1])
			_input = _input[2:]
		} else if strings.ContainsRune(" \t", rune(_input[0])) {
			// Whitespace outside quotes means end of current token
			if inWord {
				tokens = append(tokens, currentToken)
				currentToken = ""
				inWord = false
			}
			_input = _input[1:]
		} else {
			// Regular character outside quotes
			inWord = true
			currentToken += string(_input[0])
			_input = _input[1:]
		}
	}

	// Add the final token if there is one
	if currentToken != "" {
		tokens = append(tokens, currentToken)
	}

	return tokens
}

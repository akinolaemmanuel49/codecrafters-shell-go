package parser

import (
	"strings"
)

func ParseInput(input string) []string {
	// Trim any carriage returns or newlines
	_input := strings.Trim(input, "\r\n")
	var tokens []string

	// Process the input until we've consumed it all
	for len(_input) > 0 {
		// Skip leading whitespace
		_input = strings.TrimLeft(_input, " \t")
		if len(_input) == 0 {
			break
		}

		// Check if we start with a quote
		if _input[0] == '"' || _input[0] == '\'' {
			// Get the quote character
			quote := _input[0]
			// Remove the opening quote
			_input = _input[1:]

			// Find the matching closing quote, handling escaped quotes
			token := ""
			for len(_input) > 0 {
				// If we find a backslash inside a quote, check if it's escaping the quote
				if len(_input) > 1 && _input[0] == '\\' && _input[1] == quote {
					token += string(quote)
					_input = _input[2:]
				} else if _input[0] == quote {
					// End of quoted section
					_input = _input[1:]
					break
				} else {
					// Add the character to our token
					token += string(_input[0])
					_input = _input[1:]
				}
			}

			tokens = append(tokens, token)
		} else {
			// Handle non-quoted tokens, which might contain backslash escapes
			token := ""
			for len(_input) > 0 && !strings.ContainsRune(" \t", rune(_input[0])) {
				// Handle backslash escape outside quotes
				if _input[0] == '\\' && len(_input) > 1 {
					// Preserve the literal value of the next character, including space
					token += string(_input[1])
					_input = _input[2:]
				} else {
					token += string(_input[0])
					_input = _input[1:]
				}
			}

			if token != "" {
				tokens = append(tokens, token)
			}
		}
	}

	return tokens
}

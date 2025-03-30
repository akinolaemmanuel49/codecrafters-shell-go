package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// RedirectionImpl processes tokens for output redirection.
// It detects redirection operators (">") and returns:
// 1. Modified tokens list with redirection operators and targets removed
// 2. A file pointer if redirection is requested (nil otherwise)
// 3. Any error encountered during processing
func RedirectionImpl(tokens []string) ([]string, *os.File, error) {
	// If we have no tokens, return early
	if len(tokens) == 0 {
		return tokens, nil, nil
	}

	// Look for ">" redirection operator in the tokens
	for i := 0; i < len(tokens); i++ {
		if tokens[i] == ">" {
			// We found a redirection operator
			if i == len(tokens)-1 {
				// ">" is the last token, which is invalid (no target file)
				return nil, nil, errors.New("syntax error: no target for redirection")
			}

			// Get filename from the token after ">"
			filename := tokens[i+1]

			// Ensure directory exists
			dir := filepath.Dir(filename)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, nil, err
			}

			// Create or truncate the output file
			file, err := os.Create(filename)
			if err != nil {
				return nil, nil, err
			}

			// Return tokens before the redirection operator
			// This removes both the ">" and the filename from the command
			if i > 0 {
				return tokens[:i], file, nil
			}
			// If ">" is the first token, return an empty slice but valid file
			return []string{}, file, nil
		}
	}

	// No redirection found
	return tokens, nil, nil
}

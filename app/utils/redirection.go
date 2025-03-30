package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// RedirectionImpl processes tokens for output redirection.
// It supports truncating (`>`), appending (`>>`), stderr+stdout redirection (`&>`, `>&`), and appending both (`&>>`).
func RedirectionImpl(tokens []string) ([]string, *os.File, error) {
	if len(tokens) == 0 {
		return tokens, nil, nil
	}

	for i := 0; i < len(tokens); i++ {
		if strings.HasPrefix(tokens[i], ">") || strings.HasPrefix(tokens[i], "&>") {
			var filename string
			var flags int
			redirectStdErr := false // Track whether stderr should be redirected too

			if tokens[i] == ">>" || tokens[i] == "&>>" {
				// Appending mode
				flags = os.O_APPEND | os.O_CREATE | os.O_WRONLY
				redirectStdErr = (tokens[i] == "&>>")
			} else if tokens[i] == ">" || tokens[i] == "&>" || tokens[i] == ">&" {
				// Truncate mode
				flags = os.O_TRUNC | os.O_CREATE | os.O_WRONLY
				redirectStdErr = (tokens[i] == "&>" || tokens[i] == ">&")
			} else {
				continue
			}

			// Ensure there's a filename
			if i == len(tokens)-1 {
				return nil, nil, errors.New("syntax error: no target for redirection")
			}

			filename = tokens[i+1]

			// Ensure directory exists
			dir := filepath.Dir(filename)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return nil, nil, err
			}

			// Open file with proper flags
			file, err := os.OpenFile(filename, flags, 0o644)
			if err != nil {
				return nil, nil, err
			}

			// Redirect stderr if required
			if redirectStdErr {
				os.Stderr = file
			}

			// Remove redirection operators and filename from tokens
			if i > 0 {
				return tokens[:i], file, nil
			}
			return []string{}, file, nil
		}
	}

	// No redirection found
	return tokens, nil, nil
}

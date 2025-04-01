package utils

import (
	"errors"
	"os"
	"path/filepath"
)

// RedirectionImpl processes tokens for output and error redirection.
// It detects redirection operators (>, >>, 2>, 2>>, &>, &>>) and returns:
// 1. Modified tokens list with redirection operators and targets removed
// 2. File descriptors for stdout and stderr (nil if not redirected)
// 3. Any error encountered during processing
func RedirectionImpl(tokens []string) ([]string, *os.File, *os.File, error) {
	var stdoutFile, stderrFile *os.File
	var err error

	if len(tokens) == 0 {
		return tokens, nil, nil, nil
	}

	cleanTokens := []string{}
	i := 0
	for i < len(tokens) {
		if tokens[i] == ">" || tokens[i] == "1>" || tokens[i] == ">|" {
			// Standard output redirection (truncate)
			if i+1 >= len(tokens) {
				return nil, nil, nil, errors.New("syntax error: no target for redirection")
			}
			stdoutFile, err = createFile(tokens[i+1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
			i += 2
		} else if tokens[i] == ">>" || tokens[i] == "1>>" {
			// Standard output redirection (append)
			if i+1 >= len(tokens) {
				return nil, nil, nil, errors.New("syntax error: no target for appending redirection")
			}
			stdoutFile, err = createFile(tokens[i+1], os.O_CREATE|os.O_WRONLY|os.O_APPEND)
			i += 2
		} else if tokens[i] == "2>" {
			// Standard error redirection (truncate)
			if i+1 >= len(tokens) {
				return nil, nil, nil, errors.New("syntax error: no target for stderr redirection")
			}
			stderrFile, err = createFile(tokens[i+1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
			i += 2
		} else if tokens[i] == "2>>" {
			// Standard error redirection (append)
			if i+1 >= len(tokens) {
				return nil, nil, nil, errors.New("syntax error: no target for stderr appending redirection")
			}
			stderrFile, err = createFile(tokens[i+1], os.O_CREATE|os.O_WRONLY|os.O_APPEND)
			i += 2
		} else if tokens[i] == "&>" || tokens[i] == "&>|" {
			// Redirect both stdout and stderr (truncate)
			if i+1 >= len(tokens) {
				return nil, nil, nil, errors.New("syntax error: no target for &> redirection")
			}
			stdoutFile, err = createFile(tokens[i+1], os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
			stderrFile = stdoutFile // Redirect stderr to same file
			i += 2
		} else if tokens[i] == "&>>" {
			// Redirect both stdout and stderr (append)
			if i+1 >= len(tokens) {
				return nil, nil, nil, errors.New("syntax error: no target for &>> redirection")
			}
			stdoutFile, err = createFile(tokens[i+1], os.O_CREATE|os.O_WRONLY|os.O_APPEND)
			stderrFile = stdoutFile // Redirect stderr to same file
			i += 2
		} else {
			cleanTokens = append(cleanTokens, tokens[i])
			i++
		}
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return cleanTokens, stdoutFile, stderrFile, nil
}

// createFile ensures the target file exists and opens it with the given flag
func createFile(filename string, flag int) (*os.File, error) {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return os.OpenFile(filename, flag, 0o644)
}

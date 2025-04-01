package prompt

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"golang.org/x/term"
)

// PromptConfig stores configuration for the shell prompt
type PromptConfig struct {
	Prompt     string
	HistoryMax int
}

// Prompter manages the terminal input and history
type Prompter struct {
	Config   PromptConfig
	Term     *term.Terminal
	History  []string
	OldState *term.State
}

// NewPrompter creates a new prompter with the given configuration
func NewPrompter(config PromptConfig) (*Prompter, error) {
	// Check if stdin is a terminal
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return nil, errors.New("stdin is not a terminal")
	}

	// Save terminal state for restoration when exiting
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}

	// Create terminal
	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	terminal := term.NewTerminal(screen, config.Prompt)

	// Set up history with configured size
	history := make([]string, 0, config.HistoryMax)

	p := &Prompter{
		Config:   config,
		Term:     terminal,
		History:  history,
		OldState: oldState,
	}

	// Set up tab completion
	terminal.AutoCompleteCallback = p.handleAutoComplete

	return p, nil
}

// handleAutoComplete implements tab completion for commands
func (p *Prompter) handleAutoComplete(line string, pos int, key rune) (string, int, bool) {
	// Only process tab key
	if key != '\t' {
		return line, pos, false
	}

	// Get the word to complete (everything from the last space to the cursor position)
	var wordToComplete string
	if pos == 0 {
		wordToComplete = ""
	} else if pos > 0 {
		beforeCursor := line[:pos]
		lastSpaceIdx := strings.LastIndex(beforeCursor, " ")
		if lastSpaceIdx == -1 {
			// No space found, the whole line up to the cursor is the word to complete
			wordToComplete = beforeCursor
		} else {
			// There's a space, so the word to complete starts after that space
			wordToComplete = beforeCursor[lastSpaceIdx+1:]
		}
	}

	// Check if we're completing the first word (command) or an argument
	isFirstWord := !strings.Contains(line[:pos], " ")

	// Get completion candidates
	var candidates []string
	if isFirstWord {
		// Complete command (builtin or executable)
		candidates = getCommandCompletions(wordToComplete)
	} else {
		// Complete file path
		candidates = getFileCompletions(wordToComplete)
	}

	// No completions found
	if len(candidates) == 0 {
		return line, pos, false
	}

	// Single completion match
	if len(candidates) == 1 {
		// Replace the word to complete with the full match
		before := ""
		if pos > len(wordToComplete) {
			before = line[:pos-len(wordToComplete)]
		}
		after := ""
		if pos < len(line) {
			after = line[pos:]
		}

		// Add a space if this is a command completion
		suffix := " "
		if !isFirstWord && strings.HasSuffix(candidates[0], "/") {
			// Don't add space after directory name
			suffix = ""
		}

		newLine := before + candidates[0] + suffix + after
		newPos := pos - len(wordToComplete) + len(candidates[0]) + len(suffix)
		return newLine, newPos, true
	}

	// Multiple matches - find common prefix
	commonPrefix := candidates[0]
	for _, candidate := range candidates[1:] {
		commonPrefix = findCommonPrefix(commonPrefix, candidate)
	}

	// If common prefix is longer than word to complete, use that
	if len(commonPrefix) > len(wordToComplete) {
		before := ""
		if pos > len(wordToComplete) {
			before = line[:pos-len(wordToComplete)]
		}
		after := ""
		if pos < len(line) {
			after = line[pos:]
		}

		newLine := before + commonPrefix + after
		newPos := pos - len(wordToComplete) + len(commonPrefix)
		return newLine, newPos, true
	}

	// Temporarily restore terminal to normal mode for proper output
	if err := term.Restore(int(os.Stdin.Fd()), p.OldState); err != nil {
		// If we can't restore, just try to print directly
		fmt.Println()
		for _, candidate := range candidates {
			fmt.Println(candidate)
		}
	} else {
		// Terminal is now in normal mode, print completions
		fmt.Println()
		displayCompletionsInColumns(candidates)

		// Put terminal back in raw mode
		newState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err == nil {
			p.OldState = newState
		}
	}

	// Reprint prompt and current line
	fmt.Print(p.Config.Prompt + line)

	// Move cursor back to original position
	if len(line) > pos {
		// Move cursor back by the number of characters after the cursor position
		fmt.Print(strings.Repeat("\b", len(line)-pos))
	}

	return line, pos, false
}

// displayCompletionsInColumns prints completions in a neatly formatted grid
func displayCompletionsInColumns(completions []string) {
	// Sort completions for better readability
	sort.Strings(completions)

	// Get terminal width
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to standard 80 chars if we can't get terminal width
		width = 80
	}

	// Determine max item length
	maxLen := 0
	for _, item := range completions {
		if len(item) > maxLen {
			maxLen = len(item)
		}
	}

	// Add padding between columns
	colWidth := maxLen + 2

	// Calculate number of columns that can fit
	numCols := width / colWidth
	if numCols < 1 {
		numCols = 1
	}

	// Calculate number of rows needed
	numRows := (len(completions) + numCols - 1) / numCols

	// Create a 2D grid of completions
	grid := make([][]string, numRows)
	for i := range grid {
		grid[i] = make([]string, numCols)
	}

	// Fill the grid (column-major order)
	for i, item := range completions {
		col := i / numRows
		row := i % numRows
		if col < numCols {
			grid[row][col] = item
		}
	}

	// Print the grid
	for _, row := range grid {
		for colIdx, item := range row {
			if item == "" {
				continue
			}

			// Print item with padding
			fmt.Print(item)

			// Add space padding (except for last column)
			if colIdx < numCols-1 && len(item) < colWidth {
				padding := colWidth - len(item)
				fmt.Print(strings.Repeat(" ", padding))
			}
		}
		fmt.Println()
	}
}

// findCommonPrefix determines the common prefix between two strings
func findCommonPrefix(a, b string) string {
	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}

	for i := 0; i < minLen; i++ {
		if a[i] != b[i] {
			return a[:i]
		}
	}

	return a[:minLen]
}

// getCommandCompletions returns possible completions for built-in commands and executables
func getCommandCompletions(prefix string) []string {
	var completions []string

	// Add built-in commands
	for _, cmd := range commands.COMMANDS {
		if strings.HasPrefix(cmd, prefix) {
			completions = append(completions, cmd)
		}
	}

	// Add executables from PATH
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	for _, dir := range pathDirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			// Check if file is executable
			info, err := file.Info()
			if err != nil {
				continue
			}

			// Check executable bit
			if info.Mode()&0o111 != 0 && strings.HasPrefix(file.Name(), prefix) {
				completions = append(completions, file.Name())
			}
		}
	}

	return completions
}

// getFileCompletions returns possible completions for file paths
func getFileCompletions(prefix string) []string {
	// Default to current directory if no directory in prefix
	dir := "."
	basename := prefix

	// If prefix contains a path separator, split into dir and basename
	if lastSlash := strings.LastIndex(prefix, "/"); lastSlash >= 0 {
		dir = prefix[:lastSlash+1]
		if dir == "" {
			dir = "/"
		}
		basename = prefix[lastSlash+1:]
	}

	// Handle ~ expansion
	if strings.HasPrefix(dir, "~/") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			dir = homeDir + dir[1:]
		}
	} else if dir == "~" {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			dir = homeDir
		}
	}

	// Normalize path
	absDir, err := filepath.Abs(dir)
	if err == nil {
		dir = absDir
	}

	// Read directory contents
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var completions []string
	for _, file := range files {
		// Skip files that don't match prefix
		if !strings.HasPrefix(file.Name(), basename) {
			continue
		}

		// Build completion
		completion := filepath.Join(dir, file.Name())

		// For directories, add trailing slash
		if file.IsDir() {
			completion += "/"
		}

		// If directory was in prefix, preserve it in completion
		if lastSlash := strings.LastIndex(prefix, "/"); lastSlash >= 0 {
			completion = prefix[:lastSlash+1] + file.Name()
			if file.IsDir() {
				completion += "/"
			}
		} else {
			completion = file.Name()
			if file.IsDir() {
				completion += "/"
			}
		}

		completions = append(completions, completion)
	}

	return completions
}

// ReadLine reads a line from the terminal with history support
func (p *Prompter) ReadLine() (string, error) {
	line, err := p.Term.ReadLine()
	if err != nil {
		return "", err
	}

	// Add non-empty lines to history
	if line = strings.TrimSpace(line); line != "" {
		// Add to history if different from last entry
		if len(p.History) == 0 || p.History[len(p.History)-1] != line {
			// If history is full, remove oldest entry
			if len(p.History) >= p.Config.HistoryMax {
				p.History = p.History[1:]
			}
			p.History = append(p.History, line)
		}
	}

	return line, nil
}

// Close restores the terminal to its original state
func (p *Prompter) Close() error {
	if p.OldState != nil {
		err := term.Restore(int(os.Stdin.Fd()), p.OldState)
		p.OldState = nil // Prevent double-restoration
		return err
	}
	return nil
}

package main

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/prompt"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
	"golang.org/x/term"
)

// Now update main() to integrate with eval() correctly
func main() {
	// Setup signal handling first to ensure we always restore terminal state
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Create a prompter with default configuration
	prompter, err := prompt.NewPrompter(prompt.PromptConfig{
		Prompt:     "$ ",
		HistoryMax: 100,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	// Setup cleanup to happen in any exit case
	cleanup := func() {
		// First, restore the terminal to normal mode
		if prompter != nil {
			prompter.Close()
		}
		// Add a newline for cleaner transition back to normal shell
		fmt.Println()
	}

	// Handle signals for graceful shutdown
	go func() {
		<-sigs
		cleanup()
		os.Exit(1)
	}()

	// Ensure cleanup happens on normal exit
	defer cleanup()

	// Main shell loop
	for {
		// Get input from user
		input, err := prompter.ReadLine()
		if err != nil {
			if err == io.EOF {
				// Handle Ctrl+D gracefully
				return
			}
			// For other errors, restore terminal and print error
			prompter.Close()
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// Before evaluation, restore terminal state
		// This is critical when passing control to external commands
		prompter.Close()

		// Process the command
		output, err := eval(input)

		// After evaluation, reset to raw mode for our prompter
		oldState, err2 := term.MakeRaw(int(os.Stdin.Fd()))
		if err2 != nil {
			fmt.Fprintln(os.Stderr, "Failed to reset terminal mode:", err2)
			return
		}
		prompter.OldState = oldState

		// Now handle command results
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else if output != "" {
			fmt.Println(output)
		}
	}
}

// Modified version of eval to work with our prompter approach
func eval(input string) (string, error) {
	tokens, outputFile, _ := parser.ParseInput(input)
	originalStdout := os.Stdout
	if outputFile != nil {
		// Replace stdout with our file
		os.Stdout = outputFile
		defer func() {
			outputFile.Close()
			os.Stdout = originalStdout
		}()
	}
	output, err := utils.ExecuteCommand(tokens)
	return output, err
}

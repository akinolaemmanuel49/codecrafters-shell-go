package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	EXIT_SIGNAL_DEFAULT = "exit"
	EXIT_SIGNAL_0       = "exit 0"
)

func eval(command string) (string, error) {
	_command := strings.TrimSuffix(command, "\n")

	if _command == EXIT_SIGNAL_DEFAULT {
		os.Exit(0)
	}
	if _command == EXIT_SIGNAL_0 {
		os.Exit(0)
	}

	return _command + ": command not found", nil
}

func main() {
	// Wait for user input
	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprint(os.Stderr, "Error reading input:", err)
			return
		}

		output, err := eval(command)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		} else {
			fmt.Println(output)
		}
	}
}

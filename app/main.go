package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	EXIT = "exit"
	ECHO = "echo"
)

func exit(codeStr *string) {
	if codeStr == nil {
		os.Exit(0)
	}
	codeInt, err := strconv.Atoi(*codeStr)
	if err != nil {
		fmt.Fprint(os.Stderr, "Invalid exit code:", *codeStr)
		os.Exit(1)
	}
	os.Exit(codeInt)
}

func echo(args []string) {
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	fmt.Println()
}

func eval(input string) (string, error) {
	_input := strings.Split(strings.TrimSuffix(input, "\n"), " ")

	_command := _input[0]
	_args := _input[1:]

	switch _command {
	case EXIT:
		if len(_args) == 0 {
			exit(nil)
		}
		exit(&_args[0])
	case ECHO:
		echo(_args)
		return "", nil
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
			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

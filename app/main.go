package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
	PWD  = "pwd"
)

var COMMANDS = []string{
	EXIT,
	ECHO,
	TYPE,
	PWD,
}

func exitImpl(codeStr *string) {
	fmt.Fprint(os.Stdout, "exit\n")
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

func echoImpl(args []string) {
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	fmt.Println()
}

func typeImpl(args []string) {
	for i, cmd := range args {
		if slices.Contains(COMMANDS, cmd) {
			fmt.Println(args[i] + " is a shell builtin")
		} else if path, err := exec.LookPath(args[i]); err == nil {
			fmt.Printf("%s is %s\n", args[i], path)
		} else {
			fmt.Println(args[i] + ": not found")
		}
	}
}

func pwdImpl() {
	ex, err := os.Executable()
	if err != nil {
		return
	}
	exPath := filepath.Dir(ex)
	fmt.Fprint(os.Stdout, exPath)
}

func execImpl(command string, args []string) {
	if command == "" {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
		return
	}

	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
	}
}

func eval(input string) (string, error) {
	_input := strings.Split(strings.TrimSuffix(input, "\n"), " ")

	_command := _input[0]
	_args := _input[1:]

	switch _command {
	case EXIT:
		if len(_args) == 0 {
			exitImpl(nil)
		} else if len(_args) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", EXIT)
			return "", nil
		}
		exitImpl(&_args[0])
	case ECHO:
		echoImpl(_args)
		return "", nil
	case TYPE:
		typeImpl(_args)
		return "", nil
	case PWD:
		pwdImpl()
		return "", nil
	default:
		execImpl(_command, _args)
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
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return
		}

		output, err := eval(command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

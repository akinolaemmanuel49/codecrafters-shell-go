package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func eval(input string) (string, error) {
	_input := parser.ParseInput(input)

	_command := _input[0]
	_args := _input[1:]

	switch _command {
	case commands.EXIT:
		if len(_args) == 0 {
			commands.ExitImpl(nil)
		} else if len(_args) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", commands.EXIT)
			return "", nil
		}
		commands.ExitImpl(&_args[0])
	case commands.ECHO:
		commands.EchoImpl(_args)
		return "", nil
	case commands.TYPE:
		commands.TypeImpl(_args)
		return "", nil
	case commands.PWD:
		commands.PwdImpl()
		return "", nil
	case commands.CD:
		if len(_args) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", commands.CD)
			return "", nil
		} else {
			commands.CdImpl(&_args)
		}
		return "", nil
	default:
		utils.ExecImpl(_command, _args)
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

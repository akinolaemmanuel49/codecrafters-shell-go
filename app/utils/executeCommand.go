package utils

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/commands"
)

func ExecuteCommand(tokens []string) (string, error) {
	command := tokens[0]
	commandArgs := tokens[1:]

	switch command {
	case commands.EXIT:
		if len(commandArgs) == 0 {
			commands.ExitImpl(nil)
		} else if len(commandArgs) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", commands.EXIT)
			return "", nil
		}
		commands.ExitImpl(&commandArgs[0])
	case commands.ECHO:
		commands.EchoImpl(commandArgs)
		return "", nil
	case commands.TYPE:
		commands.TypeImpl(commandArgs)
		return "", nil
	case commands.PWD:
		commands.PwdImpl()
		return "", nil
	case commands.CD:
		if len(commandArgs) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", commands.CD)
			return "", nil
		} else {
			commands.CdImpl(&commandArgs)
		}
		return "", nil
	default:
		ExecImpl(command, commandArgs)
		return "", nil
	}
	return command + ": command not found", nil
}

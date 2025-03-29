package utils

import (
	"fmt"
	"os"
	"os/exec"
)

func ExecImpl(command string, args []string) {
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
		if _, ok := err.(*exec.ExitError); ok {
			return
		}
		fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
	}
}

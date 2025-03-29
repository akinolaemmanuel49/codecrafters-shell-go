package commands

import (
	"fmt"
	"os/exec"
	"slices"
)

func TypeImpl(args []string) {
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

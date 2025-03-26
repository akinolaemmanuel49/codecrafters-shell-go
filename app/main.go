package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

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
		fmt.Println(command[:len(command)-1] + ": command not found")
	}
}

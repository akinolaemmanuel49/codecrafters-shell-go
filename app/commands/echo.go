package commands

import "fmt"

func EchoImpl(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ") // Print space *before* next argument (not after)
		}
		fmt.Print(arg)
	}
	fmt.Println() // Add newline at the end (like Bash)
}

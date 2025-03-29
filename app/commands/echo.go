package commands

import "fmt"

func EchoImpl(args []string) {
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	fmt.Println()
}

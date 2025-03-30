package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/parser"
	"github.com/codecrafters-io/shell-starter-go/app/utils"
)

func eval(input string) (string, error) {
	tokens, outputFile, _ := parser.ParseInput(input)

	originalStdout := os.Stdout
	if outputFile != nil {
		// Replace stdout with our file
		os.Stdout = outputFile
		defer func() {
			outputFile.Close()
			os.Stdout = originalStdout
		}()
	}

	output, err := utils.ExecuteCommand(tokens)

	return output, err
}

func main() {
	// Wait for user input
	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

		output, err := eval(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

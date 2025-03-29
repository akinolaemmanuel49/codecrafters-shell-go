package commands

import (
	"fmt"
	"os"
	"strconv"
)

func ExitImpl(codeStr *string) {
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

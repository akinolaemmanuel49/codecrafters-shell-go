package commands

import (
	"fmt"
	"os"
)

func PwdImpl() {
	dir, _ := os.Getwd()

	fmt.Fprintln(os.Stdout, dir)
}

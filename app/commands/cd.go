package commands

import (
	"fmt"
	"os"
)

func CdImpl(args *[]string) {
	pwd, _ := os.Getwd()
	os.Setenv("OLDPWD", os.Getenv("PWD"))
	os.Setenv("PWD", pwd)
	// ENVIRONMENT_VARIABLES["OLDPWD"] = ENVIRONMENT_VARIABLES["PWD"]
	// ENVIRONMENT_VARIABLES["PWD"] = pwd

	if args == nil || len(*args) == 0 {
		os.Chdir(os.Getenv("HOME"))
		return
	}

	dir := (*args)[0]

	if dir == "~" {
		dir = os.Getenv("HOME")
	}

	if dir == "-" {
		// dir = ENVIRONMENT_VARIABLES["OLDPWD"]
		dir = os.Getenv("OLDPWD")
		if dir == "" {
			fmt.Fprintln(os.Stderr, "cd: OLDPWD not set")
			return
		}
		err := os.Chdir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
			return
		}
		return
	}

	err := os.Chdir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cd: %s: No such file or directory\n", dir)
		return
	}
}

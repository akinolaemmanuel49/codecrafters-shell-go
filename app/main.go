package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
	PWD  = "pwd"
	CD   = "cd"
)

var COMMANDS = []string{
	EXIT,
	ECHO,
	TYPE,
	PWD,
	CD,
}

var ENVIRONMENT_VARIABLES = map[string]string{}

func exitImpl(codeStr *string) {
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

func echoImpl(args []string) {
	for _, arg := range args {
		fmt.Print(arg, " ")
	}
	fmt.Println()
}

func typeImpl(args []string) {
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

func pwdImpl() {
	dir, _ := os.Getwd()

	fmt.Fprintln(os.Stdout, dir)
}

func cdImpl(args *[]string) {
	pwd, _ := os.Getwd()
	ENVIRONMENT_VARIABLES["OLDPWD"] = ENVIRONMENT_VARIABLES["PWD"]
	ENVIRONMENT_VARIABLES["PWD"] = pwd

	if args == nil || len(*args) == 0 {
		os.Chdir(os.Getenv("HOME"))
		return
	}

	dir := (*args)[0]

	if dir == "~" {
		dir = os.Getenv("HOME")
	}

	if dir == "-" {
		dir = ENVIRONMENT_VARIABLES["OLDPWD"]
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

func execImpl(command string, args []string) {
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
		fmt.Fprintf(os.Stderr, "%s: command not found\n", command)
	}
}

func eval(input string) (string, error) {
	_input := strings.Split(strings.TrimSuffix(input, "\n"), " ")

	_command := _input[0]
	_args := _input[1:]

	switch _command {
	case EXIT:
		if len(_args) == 0 {
			exitImpl(nil)
		} else if len(_args) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", EXIT)
			return "", nil
		}
		exitImpl(&_args[0])
	case ECHO:
		echoImpl(_args)
		return "", nil
	case TYPE:
		typeImpl(_args)
		return "", nil
	case PWD:
		pwdImpl()
		return "", nil
	case CD:
		if len(_args) > 1 {
			fmt.Fprintf(os.Stderr, "%s: too many arguments\n", CD)
			return "", nil
		} else {
			cdImpl(&_args)
		}
		return "", nil
	default:
		execImpl(_command, _args)
		return "", nil
	}
	return _command + ": command not found", nil
}

func main() {
	// Wait for user input
	for {
		fmt.Fprint(os.Stdout, "$ ")
		reader := bufio.NewReader(os.Stdin)
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return
		}

		output, err := eval(command)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		} else {
			if output != "" {
				fmt.Println(output)
			}
		}
	}
}

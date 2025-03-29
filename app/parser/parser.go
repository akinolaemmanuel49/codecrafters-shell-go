package parser

import (
	"regexp"
	"strings"
)

// ParseInput parses a shell-like command string into arguments, handling quotes and escapes.
func ParseInput(input string) ([]string, string, error) {
	input = strings.Trim(input, "\r\n")

	if input == "" {
		return nil, "", nil
	}

	// Extract command and argument string
	command, argstr, _ := strings.Cut(input, " ")

	var args []string

	if strings.Contains(input, "\"") {
		re := regexp.MustCompile("\"(.*?)\"")
		args = re.FindAllString(input, -1)
		for i := range args {
			args[i] = strings.Trim(args[i], "\"")
		}
	} else if strings.Contains(input, "'") {
		re := regexp.MustCompile("'(.*?)'")
		args = re.FindAllString(input, -1)
		for i := range args {
			args[i] = strings.Trim(args[i], "'")
		}
	} else {
		if strings.Contains(argstr, "\\") {
			re := regexp.MustCompile(`[^\\] +`)
			args = re.Split(argstr, -1)
			for i := range args {
				args[i] = strings.ReplaceAll(args[i], "\\", "")
			}
		} else {
			args = strings.Fields(argstr)
		}
	}

	return args, command, nil
}

package common

import (
	"os/exec"
	"strings"
)

type Command struct {
	Name string
	Args []string
	Dir  string
}

func RunCommand(c Command) []string {
	// Command to run. In this example, we'll run the "ls" command.
	cmd := exec.Command(c.Name, c.Args...)
	cmd.Dir = c.Dir

	DebugPrintLn(strings.Join(cmd.Args, " "), c.Dir)

	// Run the command and capture its output
	output, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	// Convert the output to a string and print it
	outputStr := string(output)

	// If you want to split the output into lines, you can use strings.Split
	lines := strings.Split(outputStr, "\n")

	return lines
}

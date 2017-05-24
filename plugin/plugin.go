package plugin

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecutePlugin takes a plugin command and executes it
func ExecutePlugin(pluginCommand string) {
	//	Split the entire command up using ' -' as the delimeter
	parts := strings.Split(pluginCommand, " -")

	//	The first part is the command, the rest are the args:
	head := parts[0]
	args := parts[1:len(parts)]

	//	Format the command
	cmd := exec.Command(head, args...)

	/*
		//	Sanity check -- just print out the detected args:
		for _, arg := range cmd.Args {
			log.Println(arg)
		}
	*/

	//	Sanity check -- capture stdout and stderr:
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	//	Run the command
	cmd.Run()

	//	Output our results
	fmt.Printf("Result: %v / %v", out.String(), stderr.String())
}

// FormatTokenizedString will format a string containing tokens by replacing
// the tokens with their actual values and returning the new string
func FormatTokenizedString(originalString string, tokens map[string]string) string {
	retval := originalString

	//	Replace each token with its value:
	for token, value := range tokens {
		retval = strings.Replace(retval, token, value, -1)
	}

	return retval
}

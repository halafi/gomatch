// Package unixpipe provides input funcionality for STDIN.
package file

import "io/ioutil"
import "log"
import "strings"

// ReadLog attempts to read Log data from STDIN if it's possible, if not
// it tries reading from a FilePath given in a single command line
// argument.
func ReadLine(filePath string) (logLines []string) {
	logFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	logLines = lineSplit(string(logFile))
	return logLines
}

// Function that parses a mutli-line string into single lines (array of
// strings).
func lineSplit(input string) []string {
	inputSplit := make([]string, 1)
	inputSplit[0] = input                // default single pattern, no line break
	if strings.Contains(input, "\r\n") { //CR+LF
		inputSplit = strings.Split(input, "\r\n")
	} else if strings.Contains(input, "\n") { //LF
		inputSplit = strings.Split(input, "\n")
	} else if strings.Contains(input, "\r") { //CR
		inputSplit = strings.Split(input, "\r")
	}
	return inputSplit
}

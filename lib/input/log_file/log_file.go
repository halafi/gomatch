// Package file provides input funcionality for a single file.
package file

import "io/ioutil"
import "log"
import "strings"

// ReadLog attempts to read from a file located at 'filePath' and then
// parses it into an array of strings (single lines).
func ReadLog(filePath string) (logLines []string) {
	logFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return lineSplit(string(logFile))
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

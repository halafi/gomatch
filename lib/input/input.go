// Package input provides input funcionality - read of STDIN or command
// line argument.
package input

import "io/ioutil"
import "os"
import "code.google.com/p/go.crypto/ssh/terminal"
import "log"
import "strings"

//FilePath of Pattern definitions
const patternsFilePath = "Patterns" 

//FilePath of Token definitions 
const tokensFilePath = "Tokens"

// ReadLog attempts to read Log data from STDIN if it's possible, if not
// it tries reading from a FilePath given in a single command line
// argument.
func ReadLog() (logLines []string) {
	if ! terminal.IsTerminal(0) {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		logLines = lineSplit(string(bytes))
	} else {
		if len(os.Args) == 2 { 
			logFile, err := ioutil.ReadFile(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			logLines = lineSplit(string(logFile))
		} else {
			log.Fatal("No standard input or FilePath argument given.")
		}
	}
	return logLines
}

// ReadPatterns reads a single file of patterns located at
// 'patternsFilePath' constant location.
func ReadPatterns() ([]string) {
	return lineSplit(fileToString(patternsFilePath))
}

// ReadTokens reads a single file of tokens (regex definitions) located
// at 'tokensFilePath' constant location.
func ReadTokens() ([]string) {
	return lineSplit(fileToString(tokensFilePath))
}

// Simple file reader that returns a string content of a given 
// 'filePath' file location.
func fileToString(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return string(file)
}

// Function that parses a mutli-line string into single lines (array of
// strings).
func lineSplit(input string) []string {
	inputSplit := make([]string, 1) 
	inputSplit[0] = input //default single pattern, no line break
	if strings.Contains(input, "\r\n") { //CR+LF
		inputSplit = strings.Split(input, "\r\n")
	} else if strings.Contains(input, "\n") { //LF
		inputSplit = strings.Split(input, "\n")
	} else if strings.Contains(input, "\r") { //CR
		inputSplit = strings.Split(input, "\r")
	}
	return inputSplit
}

// Package patterns provides funcionality for reading of patterns.
package patterns

import "log"
import "strings"
import "io/ioutil"

// ReadPatterns reads a single file of patterns located at
// 'filePath' argument location.
func ReadPatterns(filePath string) (output []string) {
	patternsFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	patterns := lineSplit(string(patternsFile))
	output = make([]string, 0)
	
	for i := range patterns {
		if patterns[i] == "" || patterns[i][0] == '#' {
			// skip empty lines and comments
		} else {
			patternsNameSplit := strings.Split(patterns[i], "##") //separate pattern name from its definition
			if len(patternsNameSplit) != 2 {
				log.Fatal("Error with pattern number ", i+1, " name, use [NAME##<token> word ...].")
			}
			if len(patternsNameSplit[0]) == 0 {
				log.Fatal("Error with pattern number ", i+1, ": name cannot be empty.")
			}
			if len(patternsNameSplit[1]) == 0 {
				log.Fatal("Error with pattern number ", i+1, ": pattern cannot be empty.")
			}
			newOutput := make([]string, cap(output)+1)
			copy(newOutput, output)
			output = newOutput
			output[len(output)-1] = patterns[i]
		}
	}
	return output
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

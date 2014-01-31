package main

import "io/ioutil"
import "log"
import "strings"

// ReadTokens reads a single file of tokens (regex definitions) located
// at 'filePath' argument location into map of key=token, value=regex.
func readTokens(filePath string) (output map[string]string) {
	tokensFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	tokens := lineSplit(string(tokensFile))
	output = make(map[string]string)
	for t := range tokens {
		if tokens[t] == "" || tokens[t][0] == '#' {
			// skip empty lines and comments
		} else {
			currentTokenLine := strings.Split(tokens[t], " ")
			if len(currentTokenLine) == 2 {
				output[currentTokenLine[0]] = currentTokenLine[1]
			} else {
				log.Printf("invalid token definition: \"" + tokens[t] +"\", ignoring")
			}
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

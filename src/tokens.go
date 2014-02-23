package main

import (
	"log"
	"regexp"
	"strings"
)

// addRegex validates the given string. If it's ok, adds it into given
// map of token referencing names and their compiled regexes.
func addRegex(line string, regexes map[string]*regexp.Regexp) {
	if line == "" || line[0] == '#' { // empty lines and comments
		return
	}
	lineSplit := strings.Split(line, " ") // separate token name and regex
	if len(lineSplit) != 2 {
		log.Fatal("invalid token definition: \"", line, "\"")
	}
	compiled, err := regexp.Compile(lineSplit[1])
	if err != nil {
		log.Fatal(err)
	}
	regexes[lineSplit[0]] = compiled
}

// parseTokensFile reads file at given filePath into map of token
// referencing names and their compiled regexes.
func parseTokensFile(filePath string) map[string]*regexp.Regexp {
	tokensReader := openFile(filePath)
	regexes := make(map[string]*regexp.Regexp)
	for {
		line, eof := readLine(tokensReader)
		if eof {
			break
		}
		addRegex(line, regexes)
	}
	return regexes
}

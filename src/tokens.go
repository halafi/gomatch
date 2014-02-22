// tokens.go provides funcionality for handling file with tokens.
package main

import (
	"log"
	"strings"
	"regexp"
)

// addToken takes a token line, validates it and if it's ok, adds it to
// a given map of token names and their compiled regullar expressions.
func addToken(token string, tokens map[string]*regexp.Regexp) {
	if token == "" || token[0] == '#' { // empty lines and comments
		return
	}
	tokenSplit := strings.Split(token, " ") // separate name and regex
	if len(tokenSplit) != 2 {
		log.Fatal("invalid token definition: \"", token, "\"")
	}
	compiled, err := regexp.Compile(tokenSplit[1])
	if err != nil {
		log.Fatal(err)
	}
	tokens[tokenSplit[0]] = compiled
}

// readTokens reasds all tokens at given filePath into map.
func readTokens(filePath string) map[string]*regexp.Regexp {
	tokenReader := openFile(filePath)
	tokens := make(map[string]*regexp.Regexp)
	for {
		token, eof := readLine(tokenReader)
		if eof {
			break
		}
		addToken(token, tokens)
	}
	return tokens
}

// matchToken returns true if a word matches regex, false otherwise.
func matchToken(tokens map[string]*regexp.Regexp, regex, word Token) bool {
	if tokens[regex.Value] == nil {
		log.Fatal("<", regex.Value, "> undefined")
	}
	return tokens[regex.Value].MatchString(word.Value)
}

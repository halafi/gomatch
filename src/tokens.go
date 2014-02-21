// tokens.go provides funcionality for handling file with tokens.
package main

import (
	"log"
	"strings"
	"regexp"
)

// readTokens reasds all tokens at given filePath into map.
func readTokens(filePath string) map[string]*regexp.Regexp {
	tokenReader := openFile(filePath)
	tokens := make(map[string]*regexp.Regexp)
	for {
		line, eof := readLine(tokenReader)
		if eof {
			break
		}
		token := checkToken(line)
		if token != "" {
			addToken(token, tokens)
		}
	}
	return tokens
}

// checkToken validates the given token line, if it passes the token
// line is returned.
// If empty line or comment is encountered empty string is returned.
func checkToken(token string) string {
	if token == "" || token[0] == '#' {
		return ""
	}
	tokenSplit := strings.Split(token, " ") // separate name and regex
	if len(tokenSplit) != 2 {
		log.Fatal("invalid token definition: \"", token, "\"")
	}
	return token
}

// addToken takes a token line and adds it to a given map (key = token
// name; value = regular expression for that token).
func addToken(token string, tokens map[string]*regexp.Regexp) {
	tokenSplit := strings.Split(token, " ")
	compiled, err := regexp.Compile(tokenSplit[1])
	if err != nil {
		log.Fatal(err)
	}
	tokens[tokenSplit[0]] = compiled
}

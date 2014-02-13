// tokens.go provides funcionality for handling file with tokens.
package main

import (
	"log"
	"strings"
)

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
func addToken(token string, tokens map[string]string) {
	tokenSplit := strings.Split(token, " ")
	tokens[tokenSplit[0]] = tokenSplit[1]
}

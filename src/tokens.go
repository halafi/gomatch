package main

import "log"
import "strings"

// Function checkPattern validates given token, if it passes the given
// token is returned, if empty line or comment is encounteres, then
// empty string is returned.
func checkToken(token string) string {
	if token == "" {
		return ""
	}
	if token[0] == '#' {
		return ""
	}
	tokenSplit := strings.Split(token, " ") // separate name and regex
	if len(tokenSplit) != 2 {
		log.Fatal("invalid token definition: \"", token, "\"")
	}
	return token
}

// Performs addition of token to map of tokens.
func addToken(token string, tokens map[string]string) {
	tokenSplit := strings.Split(token, " ")
	tokens[tokenSplit[0]] = tokenSplit[1]
}

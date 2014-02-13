// patterns.go provides funcionality for handling events/patterns.
package main

import (
	"log"
	"strings"
)

// checkPattern validates given pattern line, if it passes the pattern
// line is returned.
// If empty line or comment is encountered empty string is returned.
// If the pattern is invalid, error is logged and empty string returned.
func checkPattern(pattern string) string {
	if pattern == "" || pattern[0] == '#' {
		return ""
	}
	patternNameSplit := separatePatternFromName(pattern)
	if len(patternNameSplit) != 2 {
		log.Println("invalid pattern: \"", pattern+"\"")
		return ""
	}
	if len(patternNameSplit[0]) == 0 {
		log.Println("invalid pattern \"", pattern, "\": empty name")
		return ""
	}
	if len(patternNameSplit[1]) == 0 {
		log.Println("invalid pattern \"", pattern, "\": is empty")
		return ""
	}
	return pattern
}

// separatePatternFromName does what it says, splits a single line by
// separator that distinguishes pattern/event name form the rest.
func separatePatternFromName(pattern string) []string {
	return strings.Split(pattern, "##")
}

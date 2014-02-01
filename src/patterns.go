package main

import "log"
import "strings"

// Function checkPattern validates given pattern, if it passes the given
// pattern is returned, otherwise empty string is returned.
func checkPattern(pattern string) string {
	if pattern == "" {
		return ""
	}
	patternNameSplit := strings.Split(pattern, "##") //separate pattern name from its definition
	if len(patternNameSplit) != 2 {
		log.Println("invalid pattern: \"", pattern+"\"")
		return ""
	}
	if len(patternNameSplit[0]) == 0 {
		log.Println("invalid pattern \"", pattern, "\": name cannot be empty")
		return ""
	}
	if len(patternNameSplit[1]) == 0 {
		log.Println("invalid pattern \"", pattern, "\": cannot be empty")
		return ""
	}
	return pattern
}

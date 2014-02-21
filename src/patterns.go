// patterns.go provides funcionality for handling events/patterns.
package main

import (
	"log"
	"strings"
)

// Pattern is the representation of a single pattern.
// Pattern consists of: name and an array of words to match.
type Pattern struct {
	Name string
	Body []string
}

// addPattern validates the given pattern string and if its ok, adds it
// into given array of patterns and returns it. 
// Otherwise returns the old array.
func addPattern(pattern string, patterns []Pattern) []Pattern {
	if pattern != "" && pattern[0] != '#' {
		split := strings.Split(pattern, "##")
		if len(split) != 2 {
			log.Println("invalid pattern: \"", pattern+"\"")
			return patterns
		}
		if split[0] == "" {
			log.Println("invalid pattern \"", pattern, "\": empty name")
			return patterns
		}
		if split[1] == "" {
			log.Println("invalid pattern \"", pattern, "\": is empty")
			return patterns
		}
		newPattern := Pattern{split[0], strings.Split(split[1], " ")}
		newA := make([]Pattern, cap(patterns)+1)
		copy(newA, patterns)
		patterns = newA
		patterns[len(patterns)-1] = newPattern
	}
	return patterns
}

// readPatterns reads all patterns at given filePath into an array of 
// strings.
func readPatterns(filePath string) []Pattern {
	patternReader := openFile(filePath)
	patternsArr := make([]Pattern, 0)
	for {
		line, eof := readLine(patternReader)
		if eof {
			break
		}
		patternsArr = addPattern(line, patternsArr)
	}
	return patternsArr
}

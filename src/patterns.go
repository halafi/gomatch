// patterns.go provides funcionality for handling events/patterns.
package main

import (
	"log"
	"strings"
)

// Pattern is the representation of a single pattern.
// Pattern consists of: name and an array of Tokens to match.
type Pattern struct {
	Name string
	Body []Token
}

// Token represents a single thing to match. Can be regex or a word.
type Token struct {
	IsRegex bool
	Value string // i.e.: IP
	OutputName string // i.e.: ipAddress
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
		patternBody := strings.Split(split[1], " ")
		body := make([]Token, len(patternBody))
		for n := range patternBody {
			if patternBody[n][0] == '<' && patternBody[n][len(patternBody[n])-1] == '>' { 
				// add as regex
				tokenWithoutBrackets := cutWord(1, len(patternBody[n])-2, patternBody[n])
				tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
				switch len(tokenWithoutBracketsSplit) {
				case 2:
					{ // token + name, i.e. <IP:ipAddress>, OutputName ipAddress
						 body[n] = Token{true, tokenWithoutBracketsSplit[0], tokenWithoutBracketsSplit[1]}
					}
				case 1:
					{ // token only, i.e.: <IP>, OutputName IP
						body[n] = Token{true, tokenWithoutBrackets, tokenWithoutBrackets}
					}
				default:
					log.Fatal("invalid token definition: \"<" + patternBody[n] + ">\"")
				}
			} else { 
				// add as word
				body[n] = Token{false, patternBody[n], ""}
			}
		}
		newPattern := Pattern{split[0], body}
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

package main

import (
	"log"
	"regexp"
	"strings"
)

// Pattern is the representation of a single pattern/event.
// Pattern consists of: name and an array of Tokens to match.
type Pattern struct {
	Name string
	Body []Token
}

// Token represents a single thing to match. Can be regex or a word.
type Token struct {
	IsRegex    bool
	Value      string // i.e.: IP
	OutputName string // i.e.: ipAddress
}

// addPattern validates the given string. If it's ok, adds it into given
// array of patterns.
func addPattern(pattern string, patterns []Pattern, regexes map[string]*regexp.Regexp) []Pattern {
	if pattern != "" && pattern[0] != '#' { // ignore empty lines and comments
		split := strings.Split(pattern, "##")

		// check for errors
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

		// convert pattern words into Tokens
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
						if regexes[tokenWithoutBracketsSplit[0]] == nil {
							log.Printf(patternBody[n] + " undefined, failed to load event: \"" + split[0] + "\"\n")
							return patterns
						}
						body[n] = Token{true, tokenWithoutBracketsSplit[0], tokenWithoutBracketsSplit[1]}
					}
				case 1:
					{ // token only, i.e.: <IP>, OutputName IP
						if regexes[tokenWithoutBrackets] == nil {
							log.Printf(patternBody[n] + " undefined, failed to load event: \"" + split[0] + "\"\n")
							return patterns
						}
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

		// add new pattern
		newArr := make([]Pattern, cap(patterns)+1)
		copy(newArr, patterns)
		patterns = newArr
		patterns[len(patterns)-1] = Pattern{split[0], body}
	}
	return patterns
}

// readPatterns reads all patterns at given filePath into an array of
// strings.
func readPatterns(filePath string, regexes map[string]*regexp.Regexp) []Pattern {
	patternReader := openFile(filePath)
	patternsArr := make([]Pattern, 0)
	for {
		line, eof := readLine(patternReader)
		if eof {
			break
		}
		patternsArr = addPattern(line, patternsArr, regexes)
	}
	return patternsArr
}

// cutWord for a given word performs a cut (both prefix and sufix).
func cutWord(begin, end int, word string) string {
	if end >= len(word) {
		return ""
	}
	d := make([]uint8, end-begin+1)
	for j, i := 0, begin; i <= end; i, j = i+1, j+1 {
		d[j] = word[i]
	}
	return string(d)
}

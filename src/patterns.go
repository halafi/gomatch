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

// Regex represents a single string regex with it's compiled value.
type Regex struct {
	Expression string         // i.e.: ^\w+$
	Compiled   *regexp.Regexp // nil or compiled regex, if it was used
}

// readPatterns reads all patterns at given filePath into an array of
// strings.
func readPatterns(patternsFilePath, tokensFilePath string) (map[string]Regex, []Pattern) {
	regexes := parseTokensFile(tokensFilePath)
	patternReader := openFile(patternsFilePath)
	patternsArr := make([]Pattern, 0)
	for {
		line, eof := readLine(patternReader)
		if eof {
			break
		}
		regexes, patternsArr = addPattern(line, patternsArr, regexes)
	}
	return regexes, patternsArr
}

// addPattern validates the given string. If it's ok, adds it into given
// array of patterns.
func addPattern(pattern string, patterns []Pattern, regexes map[string]Regex) (map[string]Regex, []Pattern) {
	if pattern != "" && pattern[0] != '#' { // ignore empty lines and comments
		split := strings.Split(pattern, "##")

		// check for errors
		if len(split) != 2 {
			log.Println("invalid pattern: \"", pattern+"\"")
			return regexes, patterns
		}
		if split[0] == "" {
			log.Println("invalid pattern \"", pattern, "\": empty name")
			return regexes, patterns
		}
		if split[1] == "" {
			log.Println("invalid pattern \"", pattern, "\": is empty")
			return regexes, patterns
		}

		// convert pattern words into Tokens
		patternBody := strings.Split(split[1], " ")
		body := make([]Token, len(patternBody))
		for n := range patternBody {
			if patternBody[n][0] == '<' && patternBody[n][len(patternBody[n])-1] == '>' {

				regexName := cutWord(1, len(patternBody[n])-2, patternBody[n])
				outputName := regexName
				regexNameSplit := strings.Split(regexName, ":")
				if len(regexNameSplit) == 2 { // token + name, i.e. <IP:ipAddress>, OutputName ipAddress
					regexName = regexNameSplit[0]
					outputName = regexNameSplit[1]
				} else if len(regexNameSplit) != 1 { // !(token only, i.e.: <IP>, OutputName IP)
					log.Fatal("invalid token definition: \"<" + patternBody[n] + ">\"")
				}

				if regexes[regexName].Expression == "" { // missing regex in Tokens file check
					log.Printf(patternBody[n] + " undefined, failed to load event: \"" + split[0] + "\"\n")
					return regexes, patterns
				}

				if regexes[regexName].Compiled == nil { // compile regex if it wasn't compiled yet
					compiled, err := regexp.Compile(regexes[regexName].Expression)
					if err != nil {
						log.Fatal(err)
					}
					regexes[regexName] = Regex{regexes[regexName].Expression, compiled}
				}

				body[n] = Token{true, regexName, outputName} // add regex Token
			} else {
				body[n] = Token{false, patternBody[n], ""} // add word Token
			}
		}

		// add new pattern
		newArr := make([]Pattern, cap(patterns)+1)
		copy(newArr, patterns)
		patterns = newArr
		patterns[len(patterns)-1] = Pattern{split[0], body}
	}
	return regexes, patterns
}

// parseTokensFile reads file at given filePath into map of token
// referencing names and their compiled regexes (default ./Tokens).
func parseTokensFile(filePath string) map[string]Regex {
	tokensReader := openFile(filePath)
	regexes := make(map[string]Regex)
	for {
		line, eof := readLine(tokensReader)
		if eof {
			break
		}
		addRegex(line, regexes)
	}
	return regexes
}

// addRegex validates the given string. If it's ok, adds it into given
// map of token referencing names and their regexes (takes lines from
// Tokens file).
func addRegex(line string, regexes map[string]Regex) {
	if line == "" || line[0] == '#' { // empty lines and comments
		return
	}
	lineSplit := strings.Split(line, " ") // separate token name and regex
	if len(lineSplit) != 2 {
		log.Fatal("invalid token definition: \"", line, "\"")
	}
	regexes[lineSplit[0]] = Regex{lineSplit[1], nil}
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

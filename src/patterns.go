package main

import (
	"log"
	"regexp"
	"strings"
)

// Pattern represents a single pattern/event.
type Pattern struct {
	Name string  // pattern name
	Body []Token // an array of Tokens to match
}

// Token represents a single regex or a word to match.
type Token struct {
	IsRegex    bool
	Value      string // token name (i.e.: IP)
	OutputName string // name used in output (i.e.: ipAddress)
}

// Regex represents a single regular expression.
type Regex struct {
	Expression string         // regex string (i.e.: ^\w+$)
	Compiled   *regexp.Regexp // nil until first used
}

// readPatterns reads every pattern from a file.
func readPatterns(patternsFilePath, tokensFilePath string) (map[string]Regex, []Pattern) {
	regexMap := parseTokensFile(tokensFilePath)
	patternReader := openFile(patternsFilePath)
	patternsArr := make([]Pattern, 0)
	for {
		line, eof := readLine(patternReader)
		if eof {
			break
		}
		patternsArr = addPattern(string(line), patternsArr, regexMap)
	}
	return regexMap, patternsArr
}

// addPattern validates the given string. If it's ok, appends it.
func addPattern(pattern string, patterns []Pattern, regexMap map[string]Regex) []Pattern {
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
				regexName := cutWord(1, len(patternBody[n])-2, patternBody[n])

				// default: token only, i.e.: <IP>
				outputName := regexName

				regexNameSplit := strings.Split(regexName, ":")
				if len(regexNameSplit) == 2 {
					// token + name, i.e. <IP:ipAddress>
					regexName = regexNameSplit[0]
					outputName = regexNameSplit[1]
				} else if len(regexNameSplit) != 1 {
					log.Fatal("invalid token definition: \"<" + patternBody[n] + ">\"")
				}

				// check for missing regex in Tokens file
				if regexMap[regexName].Expression == "" {
					log.Printf(patternBody[n] + " undefined, failed to load event: \"" + split[0] + "\"\n")
					return patterns
				}
				
				if regexMap[regexName].Compiled == nil {
					compiled, err := regexp.Compile(regexMap[regexName].Expression)
					if err != nil {
						log.Fatal(err)
					}
					regexMap[regexName] = Regex{regexMap[regexName].Expression, compiled}
				}
				
				// check for duplicate token (i.e. pattern_name##... <MONTH> ... <MONTH> ...)
				if n > 0 {
					for i := range body {
						if body[i].OutputName == outputName {
							log.Fatal("event: \"", split[0], "\" cannot use same token name multiple times (", body[i].OutputName,")")
						}
					}
				}
				
				// add regex Token
				body[n] = Token{true, regexName, outputName}
			} else {
				// add word Token
				body[n] = Token{false, patternBody[n], ""}
			}
		}

		// add new pattern
		patterns = append(patterns, Pattern{split[0], body})
	}
	return patterns
}

// parseTokensFile reads all regex strings from a file.
func parseTokensFile(filePath string) map[string]Regex {
	r := openFile(filePath)
	regexMap := make(map[string]Regex)
	for {
		line, eof := readLine(r)
		if eof && len(line) == 0 {
			break
		}
		addRegex(string(line), regexMap)
	}
	return regexMap
}

// addRegex validates the given string. If it's ok, appends it.
func addRegex(line string, regexMap map[string]Regex) {
	if line == "" || line[0] == '#' { // empty lines and comments
		return
	}
	lineSplit := strings.Split(line, " ") // separate token name and regex
	if len(lineSplit) != 2 {
		log.Fatal("invalid token definition: \"", line, "\"")
	}
	regexMap[lineSplit[0]] = Regex{lineSplit[1], nil}
}

// cutWord removes initial and ending characters from a word.
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

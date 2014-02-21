// match.go provides the core for match handling.
package main

import (
	"log"
	"strings"
	"regexp"
)

// Match is the representation of a single event matched.
// Match consists of: type (name), and an array of matched token(s)
// followed with value that was matched.
type Match struct {
	Type string
	Body []string
}

// getMatch returns match for a given log line.
func getMatch(logLine string, patterns []Pattern, tokens map[string]*regexp.Regexp, tree map[int]map[string]int, finalFor []int) Match {
	match, matchBody := Match{}, make([]string, 0)
	current := 0
	logWords := logLineSplit(logLine)
	for i := range logWords {
		transitionTokens := getTransitionTokens(current, tree)
		validTokens := 0
		if getTransition(current, logWords[i], tree) != -1 {
			// we move by word
			current = getTransition(current, logWords[i], tree)
		} else if len(transitionTokens) > 0 {
			// we can move by some regex
			for t := range transitionTokens {
				tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
				tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
				switch len(tokenWithoutBracketsSplit) {
				case 2:
					{ // token + name, i.e. <IP:ipAddress>
						if matchToken(tokens, tokenWithoutBracketsSplit[0], logWords[i]) {
							validTokens++
							current = getTransition(current, transitionTokens[t], tree)
							matchBody = stringArraySizeUp(matchBody, 2)
							matchBody[len(matchBody)-2] = tokenWithoutBracketsSplit[1]
							matchBody[len(matchBody)-1] = logWords[i]
						}
					}
				case 1:
					{ // token only, i.e.: <IP>
						if matchToken(tokens, tokenWithoutBrackets, logWords[i]) {
							validTokens++
							current = getTransition(current, transitionTokens[t], tree)
							matchBody = stringArraySizeUp(matchBody, 2)
							matchBody[len(matchBody)-2] = tokenWithoutBrackets
							matchBody[len(matchBody)-1] = logWords[i]
						}
					}
				default:
					log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
				}
			}
			if validTokens > 1 {
				log.Fatal("multiple acceptable tokens for one word at log line:\n" + logLine + "\nfor word: \"" + logWords[i] + "\"")
			}
		} else {
			break
		}
		// leaf node - got match
		if finalFor[current] != 0 && i == len(logWords)-1 {
			if len(matchBody) > 0 {
				match = Match{patterns[finalFor[current]-1].Name, matchBody}
			} else {
				match = Match{patterns[finalFor[current]-1].Name, nil}
			}
		}
	}
	return match
}

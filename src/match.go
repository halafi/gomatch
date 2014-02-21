// match.go provides the core for match handling.
package main

import (
	"log"
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
func getMatch(logLine string, patterns []Pattern, tokens map[string]*regexp.Regexp, trie map[int]map[Token]int, finalFor []int) Match {
	match, matchBody := Match{}, make([]string, 0)
	current := 0
	logWords := logLineSplit(logLine)
	for i := range logWords {
		transitionTokens := getTransitionRegexes(current, trie)
		validTokens := 0
		if getTransition(current, Token{false, logWords[i], ""}, trie) != -1 {
			// we move by word
			current = getTransition(current, Token{false, logWords[i], ""}, trie)
		} else if len(transitionTokens) > 0 {
			// we can move by some regex
			for t := range transitionTokens {
				if matchToken(tokens, transitionTokens[t].Value, logWords[i]) {
					validTokens++
					current = getTransition(current, transitionTokens[t], trie)
					matchBody = stringArraySizeUp(matchBody, 2)
					matchBody[len(matchBody)-2] = transitionTokens[t].OutputName
					matchBody[len(matchBody)-1] = logWords[i]
				}
			}
			if validTokens > 1 {
				log.Fatal("multiple acceptable tokens for one word: \"" + logWords[i] + "\"")
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

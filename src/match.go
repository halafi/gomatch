package main

import (
	"log"
	"regexp"
)

// Match is the representation of a single event matched.
type Match struct {
	Type string   // matched event name
	Body []string // token name followed by matched value
}

// getMatch returns match for a given log line.
func getMatch(logLine string, patterns []Pattern, regexes map[string]*regexp.Regexp, trie map[int]map[Token]int, finalFor []int) Match {
	match := Match{}
	matchBody := make([]string, 0)

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
				if regexes[transitionTokens[t].Value].MatchString(logWords[i]) {
					validTokens++
					current = getTransition(current, transitionTokens[t], trie)
					matchBody = append(matchBody, transitionTokens[t].OutputName)
					matchBody = append(matchBody, logWords[i])
				}
			}
			if validTokens > 1 {
				log.Fatal("multiple acceptable tokens for one word: \"" + logWords[i] + "\"")
			}
		} else {
			break
		}

		// leaf node (ending state) -> got match
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

// logLineSplit splits a single log line string into words, words can
// only be separated by any ammount of spaces.
func logLineSplit(line string) []string {
	words := make([]string, 0)
	if line == "" {
		return words
	}
	words = append(words, "")
	wordIndex := 0
	chars := []uint8(line)
	for c := range chars {
		if chars[c] == ' ' && c < len(chars)-1 {
			if words[wordIndex] != "" {
				words = append(words, "")
				wordIndex++
			}
		} else if chars[c] != ' ' {
			words[wordIndex] = words[wordIndex] + string(chars[c])
		}
	}
	return words
}

package main

import "log"
import "strings"

// Structure used for storing a single match. Event type and a map of
// matched token(s) and their matched values (1 to 1).
type Match struct {
	Type string
	Body []string
}

// Function getMatch finds and returns match for a given log line.
func getMatch(logLine string, patterns []string, tokens map[string]string, tree map[int]map[string]int, finalFor []int) Match {
	inputMatch := Match{}
	inputMatchBody := make([]string, 0)
	words := strings.Split(logLine, " ")
	current := 0
	for w := range words {
		transitionTokens := getTransitionTokens(current, tree)
		validTokens := make([]string, 0)
		if getTransition(current, words[w], tree) != -1 { // we move by word
			current = getTransition(current, words[w], tree)
		} else if len(transitionTokens) > 0 { // we can move by some regex
			for t := range transitionTokens { // for each token leading from 'current' state
				tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
				tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
				switch len(tokenWithoutBracketsSplit) {
				case 2:
					{ // token + name, i.e. <IP:ipAddress>
						if matchToken(tokens, tokenWithoutBracketsSplit[0], words[w]) {
							validTokens = addWord(validTokens, transitionTokens[t])
						}
					}
				case 1:
					{ // token only, i.e.: <IP>
						if matchToken(tokens, tokenWithoutBrackets, words[w]) {
							validTokens = addWord(validTokens, transitionTokens[t])
						}
					}
				default:
					log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
				}
			}
			if len(validTokens) > 1 {
				log.Fatal("multiple acceptable tokens for one word at log line:\n" + logLine + "\nword: \"" + words[w] + "\"")
			} else if len(validTokens) == 1 { // we move by regex
				current = getTransition(current, validTokens[0], tree)
				inputMatchBody = stringArraySizeUp(inputMatchBody, 2)
				inputMatchBody[len(inputMatchBody)-2] = validTokens[0]
				inputMatchBody[len(inputMatchBody)-1] = words[w]
			}
		} else {
			break
		}
		if finalFor[current] != 0 && w == len(words)-1 { // leaf node - match
			patternSplit := strings.Split(patterns[finalFor[current]-1], "##")
			if len(inputMatchBody) > 0 { // body with some tokens
				inputMatch = Match{patternSplit[0], inputMatchBody}
			} else { // empty body
				inputMatch = Match{patternSplit[0], nil}
			}
		}
	}
	return inputMatch
}

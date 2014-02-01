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
			}
		} else {
			break
		}
		if finalFor[current] != 0 && w == len(words)-1 { // leaf node - match
			patternSplit := strings.Split(patterns[finalFor[current]-1], "##")
			body := getMatchBody(logLine, patternSplit[1], tokens)
			if len(body) >= 1 { // body with some tokens
				inputMatch = Match{patternSplit[0], body}
			} else { // empty body
				inputMatch = Match{patternSplit[0], nil}
			}
		}
	}
	return inputMatch
}

// Function getMatchBody returns a Match Body - map of matched token(s) 
// and their matched values (1 to 1 relation).
func getMatchBody(logLine, pattern string, tokens map[string]string) (output []string) {
	logLineWords := strings.Split(logLine, " ")
	patternWords := strings.Split(pattern, " ")
	output = make([]string, 0)
	for i := range patternWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := cutWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			switch len(tokenWithoutBracketsSplit) {
			case 2:
				{
					if matchToken(tokens, tokenWithoutBracketsSplit[0], logLineWords[i]) {
						newOutput := make([]string, cap(output)+2) // array size +2
						copy(newOutput, output)
						newOutput[len(newOutput)-2] = tokenWithoutBracketsSplit[1]
						newOutput[len(newOutput)-1] = logLineWords[i]
						output = newOutput
					}
				}
			case 1:
				{
					if matchToken(tokens, tokenWithoutBrackets, logLineWords[i]) {
						newOutput := make([]string, cap(output)+2) // array size +2
						copy(newOutput, output)
						newOutput[len(newOutput)-2] = tokenWithoutBrackets
						newOutput[len(newOutput)-1] = logLineWords[i]
						output = newOutput
					}
				}
			default:
				log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
			}
		}
	}
	return output
}

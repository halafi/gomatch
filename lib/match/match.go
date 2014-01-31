// Package match provides funcionality to access and store matches of
// log lines and patterns.
package match

import "log"
import "strings"
import "./trie"
import "../util"

// Structure used for storing a single match. Event type and a map of
// matched token(s) and their matched values (1 to 1).
type Match struct {
	Type string
	Body map[string]string
}

// GetMatch finds and returns match for a given log line.
func GetMatch(logLine string, patterns []string, tokens map[string]string, tree map[int]map[string]int, finalFor []int) Match {
	inputMatch := Match{}
	words := strings.Split(logLine, " ")
	current := 0
	for w := range words {
		transitionTokens := trie.GetTransitionTokens(current, tree)
		validTokens := make([]string, 0)
		if trie.GetTransition(current, words[w], tree) != -1 { // we move by word
			current = trie.GetTransition(current, words[w], tree)
		} else if len(transitionTokens) > 0 { // we can move by some regex
			for t := range transitionTokens { // for each token leading from 'current' state
				tokenWithoutBrackets := util.CutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
				tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
				switch len(tokenWithoutBracketsSplit) {
				case 2:
					{ // token + name, i.e. <IP:ipAddress>
						if util.MatchToken(tokens, tokenWithoutBracketsSplit[0], words[w]) {
							validTokens = util.AddWord(validTokens, transitionTokens[t])
						}
					}
				case 1:
					{ // token only, i.e.: <IP>
						if util.MatchToken(tokens, tokenWithoutBrackets, words[w]) {
							validTokens = util.AddWord(validTokens, transitionTokens[t])
						}
					}
				default:
					log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
				}
			}
			if len(validTokens) > 1 {
				log.Fatal("multiple acceptable tokens for one word at log line:\n" + logLine + "\nword: \"" + words[w] + "\"")
			} else if len(validTokens) == 1 { // we move by regex
				current = trie.GetTransition(current, validTokens[0], tree)
			}
		} else {
			break
		}
		if finalFor[current] != 0 && w == len(words)-1 { // leaf node - match
			patternSplit := strings.Split(patterns[finalFor[current]-1], "##")
			body := GetMatchBody(logLine, patternSplit[1], tokens)
			if len(body) >= 1 { // body with some tokens
				inputMatch = Match{patternSplit[0], body}
			} else { // empty body
				inputMatch = Match{patternSplit[0], nil}
			}
		}
	}
	return inputMatch
}

// GetMatchBody returns a Match Body, map of matched token(s) and their
// matched values (1 to 1 relation).
func GetMatchBody(logLine, pattern string, tokens map[string]string) (output map[string]string) {
	logLineWords := strings.Split(logLine, " ")
	patternWords := strings.Split(pattern, " ")
	output = make(map[string]string)
	for i := range patternWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := util.CutWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			switch len(tokenWithoutBracketsSplit) {
			case 2:
				{
					if util.MatchToken(tokens, tokenWithoutBracketsSplit[0], logLineWords[i]) {
						output[tokenWithoutBracketsSplit[1]] = logLineWords[i]
					}
				}
			case 1:
				{
					if util.MatchToken(tokens, tokenWithoutBrackets, logLineWords[i]) {
						output[tokenWithoutBrackets] = logLineWords[i]
					}
				}
			default:
				log.Fatal("invalid token definition: \"<" + tokenWithoutBrackets + ">\"")
			}
		}
	}
	return output
}

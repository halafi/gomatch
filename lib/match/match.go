// Package match provides funcionality to access and store matches of
// log lines and patterns.
package match

import "log"
import "strings"
import "strconv"
import "./trie"
import "./string_util"
import "./token_util"

// Structure used for storing a single match. Event type and a map of
// matched token(s) and their matched values (1 to 1).
type Match struct {
	Type string
	Body map[string]string
}

// GetMatches performs the matching function in a given Log text lines
// against the pattern definitions (pattern lines).
// Returns matches for each log line in an array of 'Match'.
func GetMatches(logLines, patterns []string, tokens map[string]string) (matchPerLine []Match) {
	tree, finalFor, stateIsTerminal := trie.ConstructPrefixTree(tokens, patterns)
	matchPerLine = make([]Match, len(logLines))
	for n := range logLines {
		words := strings.Split(logLines[n], " ")
		current := 0
		for w := range words {
			transitionTokens := trie.GetTransitionTokens(current, tree)
			validTokens := make([]string, 0)
			if trie.GetTransition(current, words[w], tree) != -1 { // we move by word
				current = trie.GetTransition(current, words[w], tree)
			} else if len(transitionTokens) > 0 { // we can move by some regex
				for t := range transitionTokens { // for each token leading from 'current' state
					tokenWithoutBrackets := string_util.CutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
					case 2:
						{ // token + name, i.e. <IP:ipAddress>
							if token_util.MatchToken(tokens, tokenWithoutBracketsSplit[0], words[w]) {
								validTokens = string_util.AddWord(validTokens, transitionTokens[t])
							}
						}
					case 1:
						{ // token only, i.e.: <IP>
							if token_util.MatchToken(tokens, tokenWithoutBrackets, words[w]) {
								validTokens = string_util.AddWord(validTokens, transitionTokens[t])
							}
						}
					default:
						log.Fatal("Problem in token definition: <" + tokenWithoutBrackets + ">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
				if len(validTokens) > 1 {
					log.Fatal("Multiple acceptable tokens for one word at log line: " + strconv.Itoa(n+1) + ", position: " + strconv.Itoa(w+1) + ".")
				} else if len(validTokens) == 1 { // we move by regex
					current = trie.GetTransition(current, validTokens[0], tree)
				}
			} else {
				break
			}
			if stateIsTerminal[current] && w == len(words)-1 { // leaf node - match
				patternSplit := strings.Split(patterns[finalFor[current]], "##")
				body := GetMatchBody(logLines[n], patternSplit[1], tokens)
				if len(body) > 1 { // body with some tokens
					matchPerLine[n] = Match{patternSplit[0], body}
				} else { // empty body
					matchPerLine[n] = Match{patternSplit[0], nil}
				}
			}
		}
	}
	return matchPerLine
}

// GetMatchBody returns a Match Body, map of matched token(s) and their
// matched values (1 to 1 relation).
func GetMatchBody(logLine, pattern string, tokens map[string]string) (output map[string]string) {
	logLineWords := strings.Split(logLine, " ")
	patternWords := strings.Split(pattern, " ")
	output = make(map[string]string)
	for i := range patternWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := string_util.CutWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			switch len(tokenWithoutBracketsSplit) {
			case 2:
				{
					if token_util.MatchToken(tokens, tokenWithoutBracketsSplit[0], logLineWords[i]) {
						output[tokenWithoutBracketsSplit[1]] = logLineWords[i]
					}
				}
			case 1:
				{
					if token_util.MatchToken(tokens, tokenWithoutBrackets, logLineWords[i]) {
						output[tokenWithoutBrackets] = logLineWords[i]
					}
				}
			default:
				log.Fatal("Problem in token definition: <" + tokenWithoutBrackets + ">, use only <TOKEN> or <TOKEN:name>.")
			}
		}
	}
	return output
}

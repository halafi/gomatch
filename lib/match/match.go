// Package match provides funcionality to access and store matches of
// log lines and patterns.
package match

import "log"
import "strings"
import "regexp"
import "strconv"

// Structure used for storing a single match. Type = event type, Body =
// map of matched token(s) and their matched values (1 to 1 relation).
type Match struct { 
	Type string
	Body map[string]string
}

// GetMatches performs the matching function in a given log lines 
// against given patterns (pattern lines).
// Returns matches for each log line in an array of 'Match'.
func GetMatches (logLines, patterns []string, tokenDefinitions map[string]string) (matchPerLine []Match) {
	trie, finalFor, stateIsTerminal := constructPrefixTree(tokenDefinitions, patterns) 
	matchPerLine = make([]Match, len(logLines))
	for n := range logLines {
		words, current := strings.Split(logLines[n], " "), 0
		for w := range words {
			transitionTokens := getTransitionTokens(current, trie)
			validTokens := make([]string, 0)
			if getTransition(current, words[w], trie) != -1 { // we can move by a word: 'words[w]'
				current = getTransition(current, words[w], trie)
			} else if len(transitionTokens) > 0 { // we can move by some regex
				for t := range transitionTokens { // for each token leading from 'current' state
					tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
						case 2: { // token + name, i.e. <IP:ipAddress>
							if MatchToken(tokenDefinitions, tokenWithoutBracketsSplit[0], words[w]) {
								validTokens = addWord(validTokens, transitionTokens[t])
							}
						}
						case 1: { // token only, i.e.: <IP>
							if MatchToken(tokenDefinitions, tokenWithoutBrackets, words[w]) {
								validTokens = addWord(validTokens, transitionTokens[t])
							}
						}
						default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
				if len(validTokens) > 1 { // we got i.e. string "user" that matches both <WORD> and i.e. <USERNAME>...
					log.Fatal("Multiple acceptable tokens for one word at log line: "+strconv.Itoa(n+1)+", position: "+strconv.Itoa(w+1)+".")	
				} else if len(validTokens) == 1 { // we can move exactly by one regex/token
					current = getTransition(current, validTokens[0], trie)
				}
			} else {
				break
			}
			if stateIsTerminal[current] && w == len(words)-1 { // we have reached leaf node in prefix tree and end of log line - got match
				patternSplit := strings.Split(patterns[finalFor[current]], "##")
				body := GetMatchBody(logLines[n], patternSplit[1], tokenDefinitions)
				
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
func GetMatchBody (logLine, pattern string, tokens map[string]string) (output map[string]string) {
	logLineWords := strings.Split(logLine, " ")
	patternWords := strings.Split(pattern, " ")
	output = make(map[string]string)
	for i := range patternWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := cutWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			switch len(tokenWithoutBracketsSplit) {
				case 2: {
					if MatchToken(tokens, tokenWithoutBracketsSplit[0], logLineWords[i]) {
						output[tokenWithoutBracketsSplit[1]] = logLineWords[i]
					}
				}
				case 1: {
					if MatchToken(tokens, tokenWithoutBrackets, logLineWords[i]) {
						output[tokenWithoutBrackets] = logLineWords[i]
					}
				} 
				default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
			}
		}
	}
	return output
}

// Function that returns constructed prefix tree/automaton for a set of 
// strings 'p' (array of patterns beginning with event name separated by
// ## and after that containing words and tokens separated by single 
// spaces each).
func constructPrefixTree (tokenDefinitions map[string]string, p []string) (trie map[int]map[string]int, finalFor []int, stateIsTerminal []bool) {
	trie = make(map[int]map[string]int)
	stateIsTerminal = make([]bool, 1)
	finalFor = make([]int, 1) 
	state := 1
	for i := range p {
		patternsNameSplit := strings.Split(p[i], "##")
		words := strings.Split(patternsNameSplit[1], " ")
		current, j := 0, 0
		for j < len(words) && getTransition(current, words[j], trie) != -1 {
			current = getTransition(current, words[j], trie)
			j++
		}
		for j < len(words) {
			newStateIsTerminal := make([]bool, cap(stateIsTerminal)+1)
			copy(newStateIsTerminal, stateIsTerminal)
			stateIsTerminal = newStateIsTerminal
			
			newFinalFor := make([]int, cap(finalFor)+1)
			copy(newFinalFor, finalFor) 
			finalFor = newFinalFor
			
			stateIsTerminal[state] = false
			if len(getTransitionWords(current, trie)) > 0 && words[j][0] == '<' && words[j][len(words[j])-1] == '>' { //conflict check when adding regex transition
				transitionWords := getTransitionWords(current, trie)
				for w := range transitionWords {
					tokenWithoutBrackets := cutWord(1, len(words[j])-2, words[j])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
						case 2: {
							if MatchToken(tokenDefinitions, tokenWithoutBracketsSplit[0], transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitionWords[w]+".")	
							}
						}
						case 1: {
							if MatchToken(tokenDefinitions, tokenWithoutBrackets, transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitionWords[w]+".")	
							}
						}
						default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
			} else if len(getTransitionTokens(current, trie)) > 0 && words[j][0] != '<' && words[j][len(words[j])-1] != '>' { //conflict check when adding word
				transitionTokens := getTransitionTokens(current, trie)
				for t := range transitionTokens {
					tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
						case 2: {
							if MatchToken(tokenDefinitions, tokenWithoutBracketsSplit[0], words[j]) {
								log.Fatal("Conflict in patterns definition, token "+transitionTokens[t]+" matches word "+words[j]+".")	
							}
						}
						case 1: {
							if MatchToken(tokenDefinitions, tokenWithoutBrackets, words[j]) {
								log.Fatal("Conflict in patterns definition, token "+transitionTokens[t]+" matches word "+words[j]+".")	
							}
						}
						default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
			}
			createTransition(current, words[j], state, trie)
			current = state
			j++
			state++
		}
		if stateIsTerminal[current] {
			log.Fatal("Duplicate pattern definition detected, pattern number: ",i+1,".")
		} else {
			stateIsTerminal[current] = true
			finalFor[current] = i
		}
	}
	return trie, finalFor, stateIsTerminal
}

// MatchToken first finds a corresponding regex for a given token,
// then attempts to match the token against given word.
// Returns true if token matches, false otherwise.
func MatchToken (tokens map[string]string, token, word string) bool {
	regex := regexp.MustCompile(tokens[token])
	if regex.MatchString(word) {
		return true
	} else {
		return false
	}
}

// Function checks if a word 'word 'exist in an array of strings, if not
// then it is added. Returns an array of strings containing 'word' and
// all of the old values
func addWord(s []string, word string) []string {
	for i := range s {
		if s[i] == word {
			return s
		}
	}
	newS := make([]string, cap(s)+1)
	copy(newS, s)
	newS[len(newS)-1] = word
	return newS
}

// Function cutWord for a given 'word' performs a cut, so that the new
// word (returned) starts at 'begin' position of the old word, and ends
// at 'end' position of the old word.
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

// If there is no state 'fromState', this function creates it, after
// that transitionion 'σ(fromState,overString) = toState' is created in
// an automaton 'at' (map with stored states and their transitions).
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	if stateExists(fromState, at) {
		at[fromState][overString]= toState
	} else {
		at[fromState] = make(map[string]int)
		at[fromState][overString]= toState
	}
}

// Returns an ending state for transition 'σ(fromState,overString)' in 
// an automaton 'at' (map with stored states and their transitions).
// Returns '-1' if there is no transition.
func getTransition(fromState int, overString string, at map[int]map[string]int) int {
	if (!stateExists(fromState, at)) {
		return -1
	}
	toState, ok := at[fromState][overString]
	if (ok == false) {
		return -1	
	}
	return toState
}

// Checks if a state 'state' exists in an automaton (map with stored 
// states and their transitions) 'at'. Returns 'true' if it does, 
// 'false' otherwise.
func stateExists(state int, at map[int]map[string]int) bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}


// Returns all transitioning tokens (without words) for a given 'state'
// in an automaton 'at' (map with stored states and their transitions)
// as an array of strings.
func getTransitionTokens(state int, at map[int]map[string]int) []string {
	transitionTokens := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] == '<' && s[len(s)-1] == '>' {
			transitionTokens = addWord(transitionTokens, s)
		}
	}
	return transitionTokens
}

// Returns all transitioning words (without tokens) for a given 'state'
// in an automaton 'at' (map with stored states and their transitions)
// as an array of strings.
func getTransitionWords(state int, at map[int]map[string]int) []string {
	transitionWords := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] != '<' && s[len(s)-1] != '>' {
			transitionWords = addWord(transitionWords, s)
		}
	}
	return transitionWords
}

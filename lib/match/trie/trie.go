// Package trie provides construction of perfix tree and supporting
// funcionality.
package trie

import "strings"
import "log"
import "../string_util"
import "../token_util"

// Function that returns constructed prefix tree/automaton for a set of
// strings 'p' (array of patterns beginning with event name separated by
// ## and after that containing words and tokens separated by single
// spaces each).
func ConstructPrefixTree(tokens map[string]string, p []string) (trie map[int]map[string]int, finalFor []int, stateIsTerminal []bool) {
	trie = make(map[int]map[string]int)
	stateIsTerminal = make([]bool, 1)
	finalFor = make([]int, 1)
	state := 1
	for i := range p {
		patternsNameSplit := strings.Split(p[i], "##")
		words := strings.Split(patternsNameSplit[1], " ")
		current, j := 0, 0
		for j < len(words) && GetTransition(current, words[j], trie) != -1 {
			current = GetTransition(current, words[j], trie)
			j++
		}
		for j < len(words) {
			newStateIsTerminal := make([]bool, cap(stateIsTerminal)+1)
			copy(newStateIsTerminal, stateIsTerminal)
			stateIsTerminal = newStateIsTerminal // array size +1
			newFinalFor := make([]int, cap(finalFor)+1)
			copy(newFinalFor, finalFor)
			finalFor = newFinalFor // array size +1

			stateIsTerminal[state] = false
			if len(GetTransitionWords(current, trie)) > 0 && words[j][0] == '<' && words[j][len(words[j])-1] == '>' { // conflict check when adding regex transition
				transitionWords := GetTransitionWords(current, trie)
				for w := range transitionWords {
					tokenWithoutBrackets := string_util.CutWord(1, len(words[j])-2, words[j])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
					case 2:
						{
							if token_util.MatchToken(tokens, tokenWithoutBracketsSplit[0], transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token " + words[j] + " matches word " + transitionWords[w] + ".")
							}
						}
					case 1:
						{
							if token_util.MatchToken(tokens, tokenWithoutBrackets, transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token " + words[j] + " matches word " + transitionWords[w] + ".")
							}
						}
					default:
						log.Fatal("Problem in token definition: <" + tokenWithoutBrackets + ">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
			} else if len(GetTransitionTokens(current, trie)) > 0 && words[j][0] != '<' && words[j][len(words[j])-1] != '>' { //conflict check when adding word transition
				transitionTokens := GetTransitionTokens(current, trie)
				for t := range transitionTokens {
					tokenWithoutBrackets := string_util.CutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
					case 2:
						{
							if token_util.MatchToken(tokens, tokenWithoutBracketsSplit[0], words[j]) {
								log.Fatal("Conflict in patterns definition, token " + transitionTokens[t] + " matches word " + words[j] + ".")
							}
						}
					case 1:
						{
							if token_util.MatchToken(tokens, tokenWithoutBrackets, words[j]) {
								log.Fatal("Conflict in patterns definition, token " + transitionTokens[t] + " matches word " + words[j] + ".")
							}
						}
					default:
						log.Fatal("Problem in token definition: <" + tokenWithoutBrackets + ">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
			}
			createTransition(current, words[j], state, trie)
			current = state
			j++
			state++
		}
		if stateIsTerminal[current] {
			log.Fatal("Duplicate pattern definition detected, pattern number: ", i+1, ".")
		} else {
			stateIsTerminal[current] = true
			finalFor[current] = i
		}
	}
	return trie, finalFor, stateIsTerminal
}

// Returns all transitioning tokens (without words) for a given 'state'
// in an automaton 'at' (map with stored states and their transitions)
// as an array of strings.
func GetTransitionTokens(state int, at map[int]map[string]int) []string {
	transitionTokens := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] == '<' && s[len(s)-1] == '>' {
			transitionTokens = string_util.AddWord(transitionTokens, s)
		}
	}
	return transitionTokens
}

// Returns all transitioning words (without tokens) for a given 'state'
// in an automaton 'at' (map with stored states and their transitions)
// as an array of strings.
func GetTransitionWords(state int, at map[int]map[string]int) []string {
	transitionWords := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] != '<' && s[len(s)-1] != '>' {
			transitionWords = string_util.AddWord(transitionWords, s)
		}
	}
	return transitionWords
}

// Returns an ending state for transition 'σ(fromState,overString)' in
// an automaton 'at' (map with stored states and their transitions).
// Returns '-1' if there is no transition.
func GetTransition(fromState int, overString string, at map[int]map[string]int) int {
	if !stateExists(fromState, at) {
		return -1
	}
	toState, ok := at[fromState][overString]
	if ok == false {
		return -1
	}
	return toState
}

// If there is no state 'fromState', this function creates it, after
// that transitionion 'σ(fromState,overString) = toState' is created in
// an automaton 'at' (map with stored states and their transitions).
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	if stateExists(fromState, at) {
		at[fromState][overString] = toState
	} else {
		at[fromState] = make(map[string]int)
		at[fromState][overString] = toState
	}
}

// Checks if a state 'state' exists in an automaton (map with stored
// states and their transitions) 'at'. Returns 'true' if it does,
// 'false' otherwise.
func stateExists(state int, at map[int]map[string]int) bool {
	_, ok := at[state]
	if !ok || state == -1 || at[state] == nil {
		return false
	}
	return true
}

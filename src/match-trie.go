// match-trie.go contains funcionality for prefix tree (trie).
package main

import (
	"log"
	"regexp"
)

// createNewTrie initializes a new prefix tree. State is the number of
// first state to be created, i is the number of first pattern to be
// added.
func createNewTrie() (trie map[int]map[Token]int, finalFor []int, state int, i int) {
	return make(map[int]map[Token]int), make([]int, 1), 1, 1
}

// appendPattern creates all the necessary transitions for a single
// pattern to the given trie.
func appendPattern(tokens map[string]*regexp.Regexp, pattern Pattern, trie map[int]map[Token]int, finalFor []int, state int, i int) ([]int, int, int) {
	current := 0
	j := 0
	for j < len(pattern.Body) && getTransition(current, pattern.Body[j], trie) != -1 {
		current = getTransition(current, pattern.Body[j], trie)
		j++
	}
	for j < len(pattern.Body) {
		finalFor = intArraySizeUp(finalFor, 1)
		finalFor[state] = 0
		
		// conflict check when adding regex transition or word transition
		if len(getTransitionWords(current, trie)) > 0 && pattern.Body[j].IsRegex {
			// conflict check 
			transitionWords := getTransitionWords(current, trie)
			for w := range transitionWords {
				if matchToken(tokens, pattern.Body[j].Value, transitionWords[w].Value) {
					log.Fatal("pattern conflict: token \"" + pattern.Body[j].Value + "\" matches word \"" + transitionWords[w].Value + "\"")
				}
			}
		} else if len(getTransitionRegexes(current, trie)) > 0 && !pattern.Body[j].IsRegex {
			transitionTokens := getTransitionRegexes(current, trie)
			for t := range transitionTokens {
				if matchToken(tokens, transitionTokens[t].Value, pattern.Body[j].Value) {
					log.Fatal("pattern conflict: token \"" + transitionTokens[t].Value + "\" matches word \"" + pattern.Body[j].Value + "\"")
				}
			}
		}
		
		createTransition(current, pattern.Body[j], state, trie)
		current = state
		j++
		state++
	}
	if finalFor[current] != 0 {
		log.Fatal("duplicate pattern detected: \"", pattern.Name, "\"")
	} else {
		// mark current state as terminal for pattern number i
		finalFor[current] = i
	}
	i++ // increment pattern number
	return finalFor, state, i
}

// match-trie.go contains funcionality for prefix tree (trie).
package main

import (
	"log"
	"regexp"
)

// createNewTrie initializes a new prefix tree.
// State is the number of first state to be created, i is the number of
// first pattern to be added.
func createNewTrie() (trie map[int]map[Token]int, finalFor []int, state int, i int) {
	return make(map[int]map[Token]int), make([]int, 1), 1, 1
}

// appendPattern creates all the necessary transitions for a single
// pattern to the given trie.
func appendPattern(tokens map[string]*regexp.Regexp, pattern Pattern, trie map[int]map[Token]int, finalFor []int, state int, i int) ([]int, int, int) {
	current := 0
	j := 0
	
	// read current pattern for as long as there are transitions
	for j < len(pattern.Body) && getTransition(current, pattern.Body[j], trie) != -1 {
		current = getTransition(current, pattern.Body[j], trie)
		j++
	}
	
	// create missing transitions
	for j < len(pattern.Body) {
		// set current state as terminal for nothing
		finalFor = append(finalFor, 0) 
		
		// iterate over all current transitions and check for conflicts
		transitions := getAllTransitions(current, trie)
		if len(transitions) > 0 {
			for t := range transitions {
				if transitions[t].IsRegex && !pattern.Body[j].IsRegex {
					if matchToken(tokens, transitions[t], pattern.Body[j]) {
						log.Fatal("pattern conflict: <" + transitions[t].Value + "> matches word " + pattern.Body[j].Value)
					}
				} else if !transitions[t].IsRegex && pattern.Body[j].IsRegex {
					if matchToken(tokens, pattern.Body[j], transitions[t]) {
						log.Fatal("pattern conflict: <" + pattern.Body[j].Value + "> matches word " + transitions[t].Value)
					}
				}
			}
		}
		
		createTransition(current, pattern.Body[j], state, trie)
		current = state
		j++
		state++
	}
	
	if finalFor[current] != 0 {
		log.Fatal("duplicate pattern detected: ", pattern.Name)
	} else { // set current state as terminal for pattern number i
		finalFor[current] = i
	}
	
	return finalFor, state, i+1 // increment pattern number and return
}

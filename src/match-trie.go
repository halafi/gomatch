package main

import (
	"log"
)

// initTrie initializes a new prefix tree.
// State is the number of first state to be created, patternNumber is
// the number of first pattern to be added, finalFor is an array of
// states with number of pattern that they are final for.
func initTrie() (trie map[int]map[Token]int, finalFor []int, state int, patternNumber int) {
	return make(map[int]map[Token]int), make([]int, 1), 1, 1
}

// appendPattern creates all the necessary transitions for a single
// pattern to the given trie.
func appendPattern(pattern Pattern, trie map[int]map[Token]int, finalFor []int, state int, patternNumber int, regexes map[string]Regex) ([]int, int, int) {
	current := 0
	j := 0

	// read current pattern for as long as there are transitions
	for j < len(pattern.Body) && getTransition(current, pattern.Body[j], trie) != -1 {
		current = getTransition(current, pattern.Body[j], trie)
		j++
	}

	// create missing transitions
	for j < len(pattern.Body) {
		finalFor = append(finalFor, 0) // current state not terminal

		// iterate over all current transitions and check for conflicts
		transitions := getAllTransitions(current, trie)
		if len(transitions) > 0 {
			for t := range transitions {
				if transitions[t].IsRegex && !pattern.Body[j].IsRegex {
					if regexes[transitions[t].Value].Compiled.MatchString(pattern.Body[j].Value) {
						log.Fatal("pattern conflict: <" + transitions[t].Value + "> matches word " + pattern.Body[j].Value)
					}
				} else if !transitions[t].IsRegex && pattern.Body[j].IsRegex {
					if regexes[pattern.Body[j].Value].Compiled.MatchString(transitions[t].Value) {
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
	} else {
		finalFor[current] = patternNumber // mark current state terminal
	}

	return finalFor, state, patternNumber + 1
}

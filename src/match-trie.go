package main

import "log"

// initTrie initializes a new prefix tree - trie.
// state is the number of first state to be created.
// patternNumber is the number of first pattern to be added.
// finalFor is an array, where index is state and value is the number of
// pattern that the state is final for. 0 if none.
// patternNumber is the number of next pattern to be added to trie.
func initTrie() (trie map[int]map[Token]int, finalFor []int, state int, patternNumber int) {
	return make(map[int]map[Token]int), make([]int, 1), 1, 1
}

// appendPattern creates all the necessary transitions for a single
// pattern to the given trie.
// First it reads the pattern for as long as there are transitions,
// after that it creates all the missing transitions while checking for
// conflicts.
func appendPattern(pattern Pattern, trie map[int]map[Token]int, finalFor []int, state int, patternNumber int, regexes map[string]Regex) ([]int, int, int) {
	current := 0
	j := 0

	for j < len(pattern.Body) && getTransition(current, pattern.Body[j], trie) != -1 {
		current = getTransition(current, pattern.Body[j], trie)
		j++
	}

	for j < len(pattern.Body) {
		finalFor = append(finalFor, 0) // current state not terminal
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

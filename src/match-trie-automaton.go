// match-trie-automaton.go provides funcions for handling an automaton
// (finite state machine) with stored transitions in a double map, 
// transitions are over struct Token (patterns.go).
package main


// getTransition returns an ending state for transition function
// σ(fromState,overString).
// Returns -1 if there is no transition.
func getTransition(fromState int, overToken Token, at map[int]map[Token]int) int {
	if !stateExists(fromState, at) {
		return -1
	}
	toState, ok := at[fromState][overToken]
	if ok == false {
		return -1
	}
	return toState
}

// createTransition creates an ending state if there isnt one yet.
// After that transitionion function σ(fromState,overToken) = toState
// is created in finite automaton given.
func createTransition(fromState int, overToken Token, toState int, at map[int]map[Token]int) {
	if stateExists(fromState, at) {
		at[fromState][overToken] = toState
	} else {
		at[fromState] = make(map[Token]int)
		at[fromState][overToken] = toState
	}
}

// stateExists returns true if a given state exists, false otherwise.
func stateExists(state int, at map[int]map[Token]int) bool {
	_, ok := at[state]
	if !ok || state == -1 || at[state] == nil {
		return false
	}
	return true
}

// getTransitionRegexes returns all transition tokens that are regexes.
func getTransitionRegexes(state int, at map[int]map[Token]int) []Token {
	transitionRegexes := make([]Token, 0)
	for s, _ := range at[state] {
		if s.IsRegex {
			transitionRegexes = append(transitionRegexes, s)
		}
	}
	return transitionRegexes
}

// getTransitionWords returns all transition Tokens that aren't regexes.
func getTransitionWords(state int, at map[int]map[Token]int) []Token {
	transitionWords := make([]Token, 0)
	for s, _ := range at[state] {
		if !s.IsRegex {
			transitionWords = append(transitionWords, s)
		}
	}
	return transitionWords
}

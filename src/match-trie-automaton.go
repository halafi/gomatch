// match-trie-automaton.go provides funcions for handling an automaton
// (finite state machine) with stored transitions in a double map, 
// transition from one state to another is over string.
package main

// getTransition returns an ending state for transition function
// σ(fromState,overString).
// Returns -1 if there is no transition.
func getTransition(fromState int, overString string, at map[int]map[string]int) int {
	if !stateExists(fromState, at) {
		return -1
	}
	toState, ok := at[fromState][overString]
	if ok == false {
		return -1
	}
	return toState
}

// createTransition creates an ending state if there isnt one yet.
// After that transitionion function σ(fromState,overString) = toState
// is created.
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	if stateExists(fromState, at) {
		at[fromState][overString] = toState
	} else {
		at[fromState] = make(map[string]int)
		at[fromState][overString] = toState
	}
}

// stateExists returns true if a given state exists, false otherwise.
func stateExists(state int, at map[int]map[string]int) bool {
	_, ok := at[state]
	if !ok || state == -1 || at[state] == nil {
		return false
	}
	return true
}

// getTransitionTokens returns all transitions begining with '<' and
// ending with '>'.
func getTransitionTokens(state int, at map[int]map[string]int) []string {
	transitionTokens := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] == '<' && s[len(s)-1] == '>' {
			transitionTokens = append(transitionTokens, s)
		}
	}
	return transitionTokens
}

// getTransitionWords returns all transitions words that don't begin
// with '<' and end with '>'.
func getTransitionWords(state int, at map[int]map[string]int) []string {
	transitionWords := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] != '<' && s[len(s)-1] != '>' {
			transitionWords = append(transitionWords, s)
		}
	}
	return transitionWords
}

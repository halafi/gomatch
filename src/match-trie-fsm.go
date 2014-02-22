// match-trie-fsm.go provides funcions for handling an automaton
// (finite state machine) with stored transitions in a double map, 
// transitions are over struct Token (patterns.go).
package main


// getTransition returns an ending state for transition function
// σ(fromState,overString).
// Returns -1 if there is no transition.
func getTransition(fromState int, overToken Token, fsm map[int]map[Token]int) int {
	if !stateExists(fromState, fsm) {
		return -1
	}
	toState, ok := fsm[fromState][overToken]
	if ok == false {
		return -1
	}
	return toState
}

// createTransition creates an ending state if there isnt one yet, after
// that creates transitionion function σ(fromState,overToken)=toState.
func createTransition(fromState int, overToken Token, toState int, fsm map[int]map[Token]int) {
	if stateExists(fromState, fsm) {
		fsm[fromState][overToken] = toState
	} else {
		fsm[fromState] = make(map[Token]int)
		fsm[fromState][overToken] = toState
	}
}

// stateExists returns true if a given state exists, false otherwise.
func stateExists(state int, fsm map[int]map[Token]int) bool {
	_, ok := fsm[state]
	if !ok || state == -1 || fsm[state] == nil {
		return false
	}
	return true
}

// getAllTransitions returns all transitions for a given state.
func getAllTransitions(state int, fsm map[int]map[Token]int) []Token {
	tansitions := make([]Token, 0)
	for s, _ := range fsm[state] {
		tansitions = append(tansitions, s)
	}
	return tansitions
}

// getTransitionRegex returns all transition tokens that are regexes.
func getTransitionRegexes(state int, fsm map[int]map[Token]int) []Token {
	tansitions := make([]Token, 0)
	for s, _ := range fsm[state] {
		if s.IsRegex {
			tansitions = append(tansitions, s)
		}
	}
	return tansitions
}

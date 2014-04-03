package main

// getTransition returns an ending state for transition σ(from,over).
// Returns -1 if there is no transition.
func getTransition(from int, over Token, fsm map[int]map[Token]int) int {
	if !stateExists(from, fsm) {
		return -1
	}
	toState, exists := fsm[from][over]
	if !exists {
		return -1
	}
	return toState
}

// createTransition creates an ending state if there isn't one yet,
// after that creates transition σ(from,over) = toState.
func createTransition(from int, over Token, toState int, fsm map[int]map[Token]int) {
	if !stateExists(from, fsm) {
		fsm[from] = make(map[Token]int)
	}
	fsm[from][over] = toState
}

// stateExists returns true if state exists, false otherwise.
func stateExists(state int, fsm map[int]map[Token]int) bool {
	if _, notEmpty := fsm[state]; !notEmpty || state == -1 || fsm[state] == nil {
		return false
	}
	return true
}

// getAllTransitions returns all the transitions leading from a state.
func getAllTransitions(state int, fsm map[int]map[Token]int) []Token {
	transitions := make([]Token, 0)
	for s, _ := range fsm[state] {
		transitions = append(transitions, s)
	}
	return transitions
}

// getTransitionRegexes returns all the transitions from a state, that
// are regular expressions.
func getTransitionRegexes(state int, fsm map[int]map[Token]int) []Token {
	transitions := make([]Token, 0)
	for s, _ := range fsm[state] {
		if s.IsRegex {
			transitions = append(transitions, s)
		}
	}
	return transitions
}

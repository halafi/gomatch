package main

// getTransition returns an ending state for transition σ(from,over).
// Returns -1 if there is no transition.
func getTransition(from int, over Token, fsm map[int]map[Token]int) int {
	if !stateExists(from, fsm) {
		return -1
	}
	toState, ok := fsm[from][over]
	if ok == false {
		return -1
	}
	return toState
}

// createTransition creates an ending state if there isn't one yet,
// after that creates transitionion σ(fromState,overToken) = toState.
func createTransition(fromState int, overToken Token, toState int, fsm map[int]map[Token]int) {
	if stateExists(fromState, fsm) {
		fsm[fromState][overToken] = toState
	} else {
		fsm[fromState] = make(map[Token]int)
		fsm[fromState][overToken] = toState
	}
}

// stateExists returns true if state exists, false otherwise.
func stateExists(state int, fsm map[int]map[Token]int) bool {
	_, ok := fsm[state]
	if !ok || state == -1 || fsm[state] == nil {
		return false
	}
	return true
}

// getAllTransitions returns all transitions leading from state.
func getAllTransitions(state int, fsm map[int]map[Token]int) []Token {
	tansitions := make([]Token, 0)
	for s, _ := range fsm[state] {
		tansitions = append(tansitions, s)
	}
	return tansitions
}

// getTransitionRegexes returns all transition tokens from state, that
// are regullar expressions.
func getTransitionRegexes(state int, fsm map[int]map[Token]int) []Token {
	tansitions := make([]Token, 0)
	for s, _ := range fsm[state] {
		if s.IsRegex {
			tansitions = append(tansitions, s)
		}
	}
	return tansitions
}

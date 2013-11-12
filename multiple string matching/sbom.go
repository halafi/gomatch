package main
import ("fmt"; "log"; /*"os";*/"strings"; "io/ioutil"; "time")

/**
 	Implementation of Set Backward Oracle Matching algorithm (Factor based aproach).
	Searches for a set of strings in file "patterns.txt" in text file text.txt.
	
	Requires two files in the same folder as the algorithm
	@file patterns.txt containing the patterns to be searched for separated by ", " 
		  !!!cannot end with ", " followed by nothing
	@file text.txt containing the text to be searched in
*/
func main() {
	patFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil {
		log.Fatal(err)
	}
	textFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	patterns := strings.Split(string(patFile), ", ")
	fmt.Printf("\nRunning: Set Backward Oracle Matching algorithm.\n\n")
	fmt.Printf("Searching for %d patterns/words:\n",len(patterns))
	for i := 0; i < len(patterns); i++ {
		if (len(patterns[i]) > len(textFile)) {
			log.Fatal("There is a pattern that is longer than text! Pattern number:", i+1)
		}
		fmt.Printf("%q ", patterns[i])
	}
	fmt.Printf("\n\nIn text (%d chars long): \n%q\n\n",len(textFile), textFile)
	startTime := time.Now()
	sbom(string(textFile), patterns)
	elapsed := time.Since(startTime)
	fmt.Printf("\nElapsed %f secs\n", elapsed.Seconds())
}

/**
	Function sbom performing the Set Backward Oracle Matching alghoritm.
	
	@param t string/text to be searched in
	@param p list of patterns to be serached for
*/  
func sbom(t string, p []string) {
	//preprocessing
	lmin := computeMinLength(p)
	or, orF, isTerminalForP := buildOracleMultiple(reverseAll(trimToLength(p, lmin)))
	f := make([]int, len(orF)) //used for storing a set of states
	fmt.Printf("\n\nSBOM: \n")
	for q := range orF {
		f[q] = -1
	}
	for i := range p {
		f[isTerminalForP[i]] = isTerminalForP[i]
		fmt.Printf("\n%q has terminal state %d", p[i], isTerminalForP[i])
	}
	//searching
	or=or
	isTerminalForP = isTerminalForP
	
	return
}

func buildOracleMultiple(p []string) (map[int]map[uint8]int, []bool, []int) {
	var parent, down int
	var o uint8
	orTrie, orTrieF, isTerminalForP := constructTrie(p)
	supply := make([]int, len(orTrieF))
	i := 0 //root of trie
	supply[i] = -1
	fmt.Printf("\n\nOracle construction: \n")
	for current := 1; current < len(orTrieF); current++ {
		o, parent = getParent(current, orTrie)
		//fmt.Printf("current %d", current)
		down = supply[parent]
		for stateExists(down, orTrie) && getTransition(down, o, orTrie) == -1 {
			createTransition(down, o, current, orTrie)
			down = supply[down]
		}
		if stateExists(down, orTrie) {
			supply[current] = getTransition(down, o, orTrie)
		} else {
			supply[current] = i
		}
		 
	}
	return orTrie, orTrieF, isTerminalForP
}

/**
	Function that constructs Trie as an automaton for a set of strings .
	Returns built triematon + array of terminal states
*/
func constructTrie(p []string) (map[int]map[uint8]int, []bool, []int) {
	var current, j int
	state := 1
	trie := make(map[int]map[uint8]int)
	isTerminal := make([]bool, 1)
	isTerminalForP := make([]int, 1)
	f := make([]int, 1)
	fmt.Printf("\n\nTrie construction: \n")
	createNewState(0, trie)
	for i:=0; i<len(p); i++ {
		current = 0
		j = 0
		for j < len(p[i]) && getTransition(current, p[i][j], trie)!=-1 {
			current = getTransition(current, p[i][j], trie)
			j++
		}
		for j < len(p[i]) {
			if state==len(isTerminal) { //dynamic array size
				newIsTerminal := make([]bool, cap(isTerminal)+1)
				copy(newIsTerminal, isTerminal) //copy(dst, src)
				isTerminal = newIsTerminal
				newF := make([]int, cap(f)+1)
				copy(newF, f)
				f = newF
			}
			createNewState(state, trie)
			isTerminal[state]=false
			createTransition(current, p[i][j], state, trie)
			current = state
			j++
			state++
		}
		if isTerminal[current] {
			f[current] = f[current] + i
		} else {
			isTerminal[current] = true
			fmt.Printf("\n%d is terminal for %q.", current, p[i])  //
			if i==len(isTerminalForP) { //dynamic array size
				newIsTerminalForP := make([]int, cap(isTerminalForP)+1)
				copy(newIsTerminalForP, isTerminalForP) //copy(dst, src)
				isTerminalForP = newIsTerminalForP
			}
			isTerminalForP[i] = current
			f[current] = i
		}
	}
	return trie, isTerminal, isTerminalForP
}

/**
	Function that takes a set of strings, desired length and trims the set of strings to that length.
*/
func trimToLength(p []string, minLength int) []string {
	for i := 0; i < len(p); i++ {
		r := []rune(p[i])
		newR := make([]rune, minLength)
		for j := 0; j < minLength; j++ {
			newR[j] = r[j]
		}
		p[i]=string(newR)
	}
	return p
}

/**
	Function that computes minimal length of a set of strings.
*/
func computeMinLength(p []string) int{
	min := len(p[0])
	for i:=1; i<len(p); i++ {
		if (len(p[i])<min) {
			min = len(p[i])
		}
	}
	return min
}

/**	
	Function that takes an array of strings and reverses it.
*/
func reverseAll(s []string) []string {
	for i := 0; i < len(s); i++ {
		s[i] = reverse(s[i])
	}
	return s
}

/**	
	Function that takes a single string and reverses it.
	
	@author Walter http://stackoverflow.com/a/10043083
	@param s string to be reversed
	@return reversed string
*/
func reverse(s string) string {
    l := len(s)
    m := make([]rune, l)
    for _, c := range s {
        l--
        m[l] = c
    }
    return string(m)
}

////Follows some AUTOMATON FUNCTIONS.
////Automaton states are stored in map[int]map[uint8]int:
//// - for each initial state(key) there is a 'value':
////   set of unique characters(keywords) with their destination states (values).
//// - lets assume, that state 0 is always the inital state of the automaton
//// - state -1 is given by some functions as a non-existing state

/**
	Function that should return previous state of a state (only works for trie (finds the first previous state in automaton).
*/
func getParent(state int, oracle map[int]map[uint8]int) (uint8, int) {
	for key, value := range oracle {
		for subkey, subvalue := range value {
			if subvalue == state {
				//fmt.Printf("\nPARENT of %d is %d", state, key)
				return subkey, key
			}
		}
	}
	//fmt.Printf("\nPARENT of %d is 0", state)
	return 'f', 0
}

/**
	Automaton function for creating a new state.
	@param state state number of state to be created
	@param oracle factor oracle to add state to
*/
func createNewState(state int, oracle map[int]map[uint8]int) {
	emptyMap := make(map[uint8]int)
	oracle[state] = emptyMap
	fmt.Printf("\ncreated state %d", state) //
}

/**
 	Automaton function for creating a transition σ(state,letter)=end.
	@usage createTransition(fromSate, overChar, toState, oracle)
*/
func createTransition(fromState int, overChar uint8, toState int, oracle map[int]map[uint8]int) {
	stateMap := oracle[fromState]
	stateMap[overChar]= toState
	oracle[fromState] = stateMap
	fmt.Printf("\n    σ(%d,%c)=%d;",fromState,overChar,toState) //

}

/**
	Returns toState from 'σ(fromState,overChar)=toState'.
	@return toState state for the desired transition function σ, -1 if there is nothing to return
*/
func getTransition(fromState int, overChar uint8, oracle map[int]map[uint8]int)(toState int) {
	var ok bool
	if (!stateExists(fromState, oracle)) {
		return -1
	}
	stateMap := oracle[fromState]
	toState, ok = stateMap[overChar]
	if (ok == false) {
		return -1	
	}
	return toState
}

/**
	Checks if state 'state' exists.
	@param state state to check for
	@param oracle oracle to check state for in
	@return true if state exists
	@return false if state doesn't exist
*/
func stateExists(state int, oracle map[int]map[uint8]int)bool {
	_, ok := oracle[state]
	if (!ok || state == -1 || oracle[state] == nil) {
		return false
	} else {
		return true
	}
}
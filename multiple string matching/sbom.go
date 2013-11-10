package main
import ("fmt"; "log"; /*"os";*/"strings"; "io/ioutil")

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
	sbom(string(textFile), patterns)
}

/**
	Function sbom performing the Set Backward Oracle Matching alghoritm.
	
	@param t string/text to be searched in
	@param p list of patterns to be serached for
*/  
func sbom(t string, p []string) {
	if(len(p) == 0) {
		return
	}
	//preprocessing
	minLength := computeMinLength(p)
	p = trimToLength(p, minLength)
	p = reverseAll(p)
	//print
		fmt.Printf("Minimum length of one pattern is: %d.\n Trimmed and then reversed patterns: ", minLength)
		for i := 0; i < len(p); i++ {
			fmt.Printf("%q ", p[i])
		}
	//print
	//or := buildOracleMultiple(p)
	//searching
	return
}

func buildOracleMultiple(p []string) map[int]map[uint8]int {
	oracle := make(map[int]map[uint8]int)
	//orTrie := map[int]map[uint8]int
	//supply := make([]int, len(p)+2) //supply function
	return oracle
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
	Automaton function for creating a new state.
	@param state state number of state to be created
	@param oracle factor oracle to add state to
*/
func createNewState(state int, oracle map[int]map[uint8]int) {
	emptyMap := make(map[uint8]int)
	oracle[state] = emptyMap
	//if(runInSilentMode==false) {
		fmt.Printf("\ncreated state %d", state)
	//}
}

/**
 	Automaton function for creating a transition σ(state,letter)=end.
	@usage createTransition(fromSate, overChar, toState, oracle)
*/
func createTransition(fromState int, overChar uint8, toState int, oracle map[int]map[uint8]int) {
	stateMap := oracle[fromState]
	stateMap[overChar]= toState
	oracle[fromState] = stateMap
	//if(runInSilentMode==false) {
		fmt.Printf("\n    σ(%d,%c)=%d;",fromState,overChar,toState)
	//}
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
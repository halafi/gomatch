package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time")

/**
 	Implementation of Set Backward Oracle Matching algorithm (Factor based aproach).
	Searches for a set of strings in file "patterns.txt" in text file text.txt.
	
	Requires two files in the same folder as the algorithm
	@file patterns.txt containing the patterns to be searched for separated by "," 
		  !!!cannot end with "," followed by nothing
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
	patterns := strings.Split(string(patFile), ",")
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
	occurences := make(map[int][]int)
	occArray := make([]int, 0)
	var j, current int
	var word string
	n := len(t)
	lmin := computeMinLength(p)
	or, orF, f := buildOracleMultiple(reverseAll(trimToLength(p, lmin)))
	orF=orF //probably not needed
	fmt.Printf("\n\nSBOM: \n\n")
	/*for q := range orF {
		f[q] = make([]int, 0)
	}
	for i := range p {
		f[f[i]] = f[i]+p[i]
		fmt.Printf("%q has terminal state %d\n", p[i], f[i])
	}*/
	fmt.Println()
	//searching
	pos := 0
	for pos <= n - lmin {
		current = 0
		j = lmin
		fmt.Printf("Position: %d, we read: ", pos)
		for j >= 1 && stateExists(current, or) {
			fmt.Printf("%c", t[pos+j-1])
			current = getTransition(current, t[pos+j-1], or)
			if (current == -1) {
				fmt.Printf(" (FAIL) ")
			} else {
				fmt.Printf(", ")
			}
			j--
		}
		fmt.Printf("in the factor oracle. \n")
		word = getWord(pos, pos+lmin-1, t)
		if stateExists(current, or) && j == 0 && strings.HasPrefix(word, getCommonPrefix(p, f[current], lmin)) {
			for i := range p {
				//fmt.Printf("if %q = %q", p[i], word)
				if p[i] == word {
					//occurence
					fmt.Printf("- Occurence\n")
					occArray = occurences[current]
					newArray := make([]int, cap(occArray)+1)
					copy(newArray, occArray) //copy(dst, src)
					occArray = newArray
					occurences[i] = occArray
					occurences[i][len(occArray)-1] = pos
				}
			}
			j = 0
		}
		pos = pos + j + 1
	}
	fmt.Printf("\n\n")
	for key, value := range occurences {
		fmt.Printf("\nThere were %d occurences for word: %q at positions ",len(value), p[key])
		for i := range value {
			fmt.Printf("%d", value[i])
			if i != len(value)-1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf(".")
	}
	return
}

/**
	Returns a prefix of size l.
*/
func getCommonPrefix(p []string, array []int, l int) (prefix string) {
	for i := range array {
		r := []rune(p[array[i]])
		newR := make([]rune, l)
		for j := 0; j < l; j++ {
			newR[j] = r[j]
		}
		prefix = string(newR)
	}
	return prefix
}

/**
	Function that takes a set of strings, desired length and trims the set of strings to that length.
*/
func trimToLength(p []string, minLength int) (trimmed []string) {
	trimmed = make([]string, len(p))
	for i := 0; i < len(p); i++ {
		r := []rune(p[i])
		newR := make([]rune, minLength)
		for j := 0; j < minLength; j++ {
			newR[j] = r[j]
		}
		trimmed[i]=string(newR)
	}
	return trimmed
}

/**
	Function that return word in text at position begin - end.
*/
func getWord(begin, end int, t string) string {
	d := make([]uint8, end-begin+1)
	for j,i:= 0,begin; i<=end; i++ {
		d[j] = t[i]
		j++
	}
	s2 := string(d)
	return s2
}

/**
	Functions that build factor oracle.
*/
func buildOracleMultiple(p []string) (map[int]map[uint8]int, []bool, map[int][]int) {
	var parent, down int
	var o uint8
	orTrie, orTrieF, f := constructTrie(p)
	toReturn := orTrie //for getParent to have only one parent
	supply := make([]int, len(orTrieF))
	i := 0 //root of trie
	supply[i] = -1
	fmt.Printf("\n\nOracle construction: \n")
	for current := 1; current < len(orTrieF); current++ {
		o, parent = getParent(current, orTrie) //getParent might fail
		fmt.Printf("\nparent of %d is %d", current, parent)
		down = supply[parent]
		for stateExists(down, toReturn) && getTransition(down, o, toReturn) == -1 {
			createTransition(down, o, current, toReturn)
			down = supply[down]
		}
		if stateExists(down, toReturn) {
			supply[current] = getTransition(down, o, toReturn)
		} else {
			supply[current] = i
		}
		 
	}
	return orTrie, orTrieF, f
}

/**
	Function that constructs Trie as an automaton for a set of strings .
	Returns built triematon + array of terminal states
*/
func constructTrie(p []string) (map[int]map[uint8]int, []bool, map[int][]int) {
	var current, j int
	state := 1
	trie := make(map[int]map[uint8]int)
	isTerminal := make([]bool, 1)
	array := make([]int, 0)
	f := make(map[int][]int)  //0-1,2; 1-7; (terminal states for pattern i in f[i]
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
			}
			createNewState(state, trie)
			isTerminal[state]=false
			createTransition(current, p[i][j], state, trie)
			current = state
			j++
			state++
		}
		if isTerminal[current] {
			array = f[current]
			newArray := make([]int, cap(array)+1)
			copy(newArray, array) //copy(dst, src)
			array = newArray
			array[len(array)-1] = i
			f[current] = array
			fmt.Printf(" and %d", i)
		} else {
			isTerminal[current] = true //mark current as terminal
			fmt.Printf("\n%d is terminal for word number %d", current, i) 
			newArray := make([]int, 1)
			copy(newArray, array) //copy(dst, src)
			array = newArray
			array[len(array)-1] = i
			f[current] = array
		}
	}
	return trie, isTerminal, f
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
func reverseAll(s []string) (reversed []string) {
	reversed = make([]string, len(s))
	for i := 0; i < len(s); i++ {
		reversed[i] = reverse(s[i])
	}
	return reversed
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
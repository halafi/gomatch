package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time")
const debugMode bool = true //slow, if set to 'true' - prints various stuff out
/**
 	Implementation of Set Backward Oracle Matching algorithm.
	Searches for a given set of strings in file 'patterns.txt' in text that is in file 'text.txt'.
	Finds and prints occurences of each pattern.
	Requires two files in the same folder as the algorithm:
	
	@file patterns.txt containing the patterns to be searched for separated by ", " 
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
	if debugMode==true { 
		fmt.Printf("Searching for %d patterns/words:\n",len(patterns))
	}
	for i := 0; i < len(patterns); i++ {
		if (len(patterns[i]) > len(textFile)) {
			log.Fatal("There is a pattern that is longer than text! Pattern number:", i+1)
		}
		if debugMode==true { 
			fmt.Printf("%q ", patterns[i])
		}
	}
	if debugMode==true { 
		fmt.Printf("\n\nIn text (%d chars long): \n%q\n\n",len(textFile), textFile)
	}
	startTime := time.Now()
	sbom(string(textFile), patterns)
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nElapsed %f secs\n", elapsed.Seconds())
}

/**
	Function sbom performing the Set Backward Oracle Matching alghoritm.
	
	@param t string/text to be searched in
	@param p list of patterns to be serached for
*/  
func sbom(t string, p []string) {
	occurences := make(map[int][]int)
	occArray := make([]int, 0)
	lmin := computeMinLength(p)
	or, f := buildOracleMultiple(reverseAll(trimToLength(p, lmin)))
	if debugMode==true {
		fmt.Printf("\n\nSBOM:\n\n")
	}
	pos := 0	//searching
	for pos <= len(t) - lmin {
		current := 0
		j := lmin
		if debugMode==true {
			fmt.Printf("Position: %d, we read: ", pos)
		}
		for j >= 1 && stateExists(current, or) {
			if debugMode==true {
				fmt.Printf("%c", t[pos+j-1])
			}
			current = getTransition(current, t[pos+j-1], or)
			if debugMode==true {
				if (current == -1) {
					fmt.Printf(" (FAIL) ")
				} else {
					fmt.Printf(", ")
				}
			}
			j--
		}
		if debugMode==true {
			fmt.Printf("in the factor oracle. \n")
		}
		word := getWord(pos, pos+lmin-1, t)
		if stateExists(current, or) && j == 0 && strings.HasPrefix(word, getCommonPrefix(p, f[current], lmin)) {
			for i := range f[current] {
				if p[f[current][i]] == getWord(pos, pos-1+len(p[f[current][i]]), t) { //occurence
					if debugMode==true {
						fmt.Printf("- Occurence, %q = %q\n", p[f[current][i]], word)
					}
					occArray = occurences[f[current][i]]
					newArray := make([]int, cap(occArray)+1)
					copy(newArray, occArray) //copy(dst, src)
					occArray = newArray
					occurences[f[current][i]] = occArray
					occurences[f[current][i]][len(occArray)-1] = pos
				}
			}
			j = 0
		}
		pos = pos + j + 1
	}
	for key, value := range occurences {
		fmt.Printf("\nThere were %d occurences for word: %q at positions: ",len(value), p[key])
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
	Functions that build factor oracle.
*/
func buildOracleMultiple(p []string) (toReturn map[int]map[uint8]int, f map[int][]int) {
	orTrie, orTrieF, f := constructTrie(p)
	supply := make([]int, len(orTrieF))
	toReturn = orTrie
	i := 0 //root of trie
	supply[i] = -1
	if debugMode==true {
		fmt.Printf("\n\nOracle construction: \n")
	}
	for current := 1; current < len(orTrieF); current++ {
		o, parent := getParent(current, orTrie)
		down := supply[parent]
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
	return toReturn, f
}

/**
	Function that constructs Trie as an automaton for a set of strings .
	Returns built triematon + array of terminal states
*/
func constructTrie(p []string) (map[int]map[uint8]int, []bool, map[int][]int) {
	trie := make(map[int]map[uint8]int)
	isTerminal := make([]bool, 1)
	array := make([]int, 0)
	f := make(map[int][]int)  //terminal states for pattern i in f[i]
	state := 1
	if debugMode==true {
		fmt.Printf("\n\nTrie construction: \n")
	}
	createNewState(0, trie)
	for i:=0; i<len(p); i++ {
		current := 0
		j := 0
		for j < len(p[i]) && getTransition(current, p[i][j], trie) != -1 {
			current = getTransition(current, p[i][j], trie)
			j++
		}
		for j < len(p[i]) {
			if state==len(isTerminal) {
				newIsTerminal := make([]bool, cap(isTerminal)+1)
				copy(newIsTerminal, isTerminal)
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
			copy(newArray, array)
			array = newArray
			array[len(array)-1] = i
			f[current] = array
			if debugMode==true {
				fmt.Printf(" and %d", i)
			}
		} else {
			isTerminal[current] = true
			if debugMode==true {
				fmt.Printf("\n%d is terminal for word number %d", current, i) 
			}
			newArray := make([]int, 1)
			copy(newArray, array)
			array = newArray
			array[len(array)-1] = i
			f[current] = array
		}
	}
	return trie, isTerminal, f
}

/*******************          String functions          *******************/
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

/**
	Returns a prefix of size length.
*/
func getCommonPrefix(p []string, array []int, length int) (prefix string) {
	for i := range array {
		r := []rune(p[array[i]])
		newR := make([]rune, length)
		for j := 0; j < length; j++ {
			newR[j] = r[j]
		}
		prefix = string(newR)
	}
	return prefix
}

/**
	Function that takes a set of strings and their desired length, and then trims the set of strings to that length.
*/
func trimToLength(p []string, minLength int) (trimmed []string) {
	trimmed = make([]string, len(p))
	for i := range p {
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
	Function that returns word in text at position from 'begin' to 'end'.
*/
func getWord(begin, end int, t string) (s2 string) {
	d := make([]uint8, end-begin+1)
	for j,i:= 0,begin; i<=end; i++ {
		d[j] = t[i]
		j++
	}
	s2 = string(d)
	return s2
}

/**
	Function that computes minimal length or single string for a set of strings.
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

/*******************          Automaton functions          *******************/
/**
	Function that should return previous state of a state 
    (only works for trie - finds the first previous state in automaton).
*/
func getParent(state int, at map[int]map[uint8]int) (uint8, int) {
	for key, value := range at {
		for subkey, subvalue := range value {
			if subvalue == state {
				return subkey, key
			}
		}
	}
	return 'f', 0
}

/**
	Automaton function for creating a new state 'state'.
*/
func createNewState(state int, at map[int]map[uint8]int) {
	emptyMap := make(map[uint8]int)
	at[state] = emptyMap
	if debugMode==true {
		fmt.Printf("\ncreated state %d", state)
	}
}

/**
 	Automaton function for creating a transition σ(state,letter)=end.
*/
func createTransition(fromState int, overChar uint8, toState int, at map[int]map[uint8]int) {
	stateMap := at[fromState]
	stateMap[overChar]= toState
	at[fromState] = stateMap
	if debugMode==true {
		fmt.Printf("\n    σ(%d,%c)=%d;",fromState,overChar,toState)
	}
}

/**
	Returns toState from 'σ(fromState,overChar)=toState'.
	@return toState state for the desired transition function σ, -1 if there is nothing to return
*/
func getTransition(fromState int, overChar uint8, at map[int]map[uint8]int)(toState int) {
	var ok bool
	if (!stateExists(fromState, at)) {
		return -1
	}
	stateMap := at[fromState]
	toState, ok = stateMap[overChar]
	if (ok == false) {
		return -1	
	}
	return toState
}

/**
	Checks if state 'state' exists. Returns true if it does, false otherwise.
*/
func stateExists(state int, at map[int]map[uint8]int)bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	} else {
		return true
	}
}
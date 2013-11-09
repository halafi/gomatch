package main
import ("fmt"; "log"; "os"; "io/ioutil")

/** user defined CONSTANTS
	Set runInSilentMode to:
		@true to run in silent mode
		@false to print everything
	Set commandLineInput to:
		@true to take two command line arguments
		@false to take two files "pattern.txt" AND "text.txt"
*/
const runInSilentMode bool = true //very slow
const commandLineInput bool = false

/**
 	Implementation of Backward Oracle Matching algorithm (Factor based aproach).
	
	IF(commandLineInput == true) Requires two command line arguments.
	@argument string to be searched "for" (pattern, search word), no spaces allowed
	@argument one space
	@argument string to be searched "in" (text), single spaces allowed
	
	IF(commandLineInput == false) requires two files in the same folder
	@file pattern.txt containing the pattern to be searched for
	@file text.txt containing the text to be searched in
*/
func main() {
	if (commandLineInput == true) { //in case of command line input
		args := os.Args
		if (len(args) <= 2) {
			log.Fatal("Not enough arguments. Two string arguments separated by spaces are required!")
		}
		pattern := args[1]
		s := args[2]
		for i := 3; i<len(args); i++ {
			s = s +" "+ args[i]
		}
		if ( len(args[1]) > len(s) ) {
			log.Fatal("Pattern  is longer than text!")
		} 
		if(runInSilentMode==false) {
			fmt.Printf("\nRunning: Backward Oracle Matching alghoritm.\n\n")
			fmt.Printf("Search word (%d chars long): %q.\n",len(args[1]), pattern)
			fmt.Printf("Text        (%d chars long): %q.\n\n",len(s), s)
		} else {
			fmt.Printf("\nRunning: Backward Oracle Matching alghoritm in SILENT mode (see line 12 in the code).")
		}
		bom(s, pattern)
	} else if (commandLineInput == false) { //in case of file line input
		patFile, err := ioutil.ReadFile("pattern.txt")
		if err != nil {
			log.Fatal(err)
		}
		textFile, err := ioutil.ReadFile("text.txt")
		if err != nil {
			log.Fatal(err)
		}
		if (len(patFile) > len(textFile)) {
			log.Fatal("Pattern  is longer than text!")
		}
		if(runInSilentMode==false) {
			fmt.Printf("\nRunning: Backward Oracle Matching alghoritm.\n\n")
			fmt.Printf("Search word (%d chars long): %q.\n",len(patFile), patFile)
			fmt.Printf("Text        (%d chars long): %q.\n\n",len(textFile), textFile)
		} else {
			fmt.Printf("\nRunning: Backward Oracle Matching alghoritm in SILENT mode (see line 12 in the code).")
		}
		bom(string(textFile), string(patFile))
	}
}

/**
	Function bom performing the Backward Oracle Matching alghoritm.
    Prints whether the word/pattern was found + positions of possible multiple occurences
	or that the word was not found.
	
	@param t string/text to be searched in
	@param p pattern/word to be serached for
*/  
func bom(t, p string) {
	n, m := len(t), len(p)
	var current, j, pos int
	oracle := oracleOnLine(reverse(p))
	occurences := make([]int, len(t))
	currentOcc := 0
	pos = 0
	if(runInSilentMode==false) {
		fmt.Printf("\n\nWe are reading backwards in %q, searching for %q\n\nat position %d:\n",t, p, pos+m-1)
	}
	for (pos <= n - m) {
		current = 0 //initial state of the oracle
		j = m
		for j > 0 && stateExists(current, oracle) {
			if(runInSilentMode==false) {
				prettyPrint(current, j, n, pos, t, oracle)
			}
			current = getTransition(current, t[pos+j-1], oracle)
			j--
		}
		if stateExists(current, oracle){
			if(runInSilentMode==false) {
				fmt.Printf(" We got an occurence!")
			}
			occurences[currentOcc] = pos
			currentOcc++
		}
		pos = pos + j +1
		if (pos+m-1 < len(t)) {
			if(runInSilentMode==false) {
				fmt.Printf("\n\nposition %d:\n",pos+m-1)
			}
		}
	}
	fmt.Printf("\n\n")
	if (currentOcc > 0) {
		fmt.Printf("\nWord %q was found at positions: {", p)
		for k := 0; k<currentOcc-1; k++ {
			fmt.Printf("%d, ",occurences[k])
		}
		fmt.Printf("%d",occurences[currentOcc-1])
		fmt.Printf("} in %q.\n", t)
	}
	if(currentOcc == 0) {
		fmt.Printf("\nWord was not found.\n")
	}
	return
}

/**
	Construction of the factor oracle automaton for a word p.

	@param p pattern to be added
	@param supply supply map
	@return oracle built oracle
*/
func oracleOnLine(p string)(oracle map[int]map[uint8]int) {
	if(runInSilentMode==false) {
		fmt.Printf("Oracle construction: \n")
	}
	oracle = make(map[int]map[uint8]int)
	supply := make([]int, len(p)+2) //supply function
	createNewState(0, oracle)
	supply[0]=-1
	var orP string
	for j := 0; j < len(p); j++ {
		oracle, orP = oracleAddLetter(oracle, supply, orP, p[j])
	}
	return oracle
}

/**
	Adds one letter to the oracle.

	@param oracle oracle to add letter to
	@param p pattern (not whole, contained in oracle)
	@param o letter to be added
	@param supply supply map
*/
func oracleAddLetter(oracle map[int]map[uint8]int, supply []int, orP string, o uint8)(oracleToReturn map[int]map[uint8]int, orPToReturn string) { 
	m := len(orP)
	var s int
	createNewState(m + 1, oracle)
	createTransition(m, o, m + 1, oracle)
	k := supply[m]
	for k > -1 && getTransition(k,o, oracle) == -1 {
		createTransition(k, o, m + 1, oracle)
		k = supply[k]
	}
	if (k == -1) {
		s = 0
	} else {
		s = getTransition(k,o, oracle)
	}
	supply[m+1] = s
	return oracle, orP+string(o)
}

/**	
	Function that takes string and reverses it.
	
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
	if(runInSilentMode==false) {
		fmt.Printf("\ncreated state %d", state)
	}
}

/**
 	Automaton function for creating a transition σ(state,letter)=end.
	@usage createTransition(fromSate, overChar, toState, oracle)
*/
func createTransition(fromState int, overChar uint8, toState int, oracle map[int]map[uint8]int) {
	stateMap := oracle[fromState]
	stateMap[overChar]= toState
	oracle[fromState] = stateMap
	if(runInSilentMode==false) {
		fmt.Printf("\n    σ(%d,%c)=%d;",fromState,overChar,toState)
	}
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

/**
	Just some printing of what the alghoritm does.
*/
func prettyPrint(current int, j int, n int, pos int, t string, oracle map[int]map[uint8]int) {
	if (current == 0 && !(getTransition(current, t[pos+j-1], oracle) == -1)) {
		fmt.Printf("\n -->(%d)---(%c)--->(%d)", current, t[pos+j-1], getTransition(current, t[pos+j-1], oracle))
	} else if (getTransition(current, t[pos+j-1], oracle) == -1 && current !=0) {
		fmt.Printf("\n    (%d)---(%c)       ", current, t[pos+j-1])
	} else if (getTransition(current, t[pos+j-1], oracle) == -1 && current ==0) {
		fmt.Printf("\n -->(%d)---(%c)       ", current, t[pos+j-1])
	} else {
		fmt.Printf("\n    (%d)---(%c)--->(%d)", current, t[pos+j-1], getTransition(current, t[pos+j-1], oracle))
	}
	fmt.Printf(" ")
	for a := 0; a < pos+j-1; a++ {
		fmt.Printf("%c", t[a])
	}
	if (getTransition(current, t[pos+j-1], oracle) == -1) {
		fmt.Printf("[%c]", t[pos+j-1])
	} else {
		fmt.Printf("[%c]", t[pos+j-1])
	}
	for a := pos+j; a<n; a++ {
			fmt.Printf("%c", t[a])
	}
	if (getTransition(current, t[pos+j-1], oracle) == -1) {
		fmt.Printf(" FAIL on the character[%c]", t[pos+j-1])
	}	
}
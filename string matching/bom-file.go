package main
import (
	"fmt" //implements fomratted I/O.
	"log" //simple logging package
	"io/ioutil" // some I/O utility functions
)

/* 	Implementation of Backward Oracle Matching algorithm (Factor based aproach).
	Requires two files in the folder with this file:
	
	@File pattern.txt containing the pattern to be searched for
	@File text.txt containing the text to be searched in
*/
func main() {
	// Error handling & file input
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
	// Alghoritm execution
	fmt.Printf("\nRunning: Backward Oracle Matching alghoritm.\n\n")
	fmt.Printf("Search word (%d chars long): %q.\n",len(patFile), patFile)
	fmt.Printf("Text        (%d chars long): %q.\n\n",len(textFile), textFile)
	bom(string(textFile), string(patFile))
}

/*  Function bom performing the Backward Oracle Matching alghoritm.
    Prints whether the word/pattern was found and on what position in the text or not.
	
	@param t string/text to be searched in
	@param p pattern/word to be serached for
*/  
func bom(t, p string) {
	oracle := oracleOnLine(p)
	fmt.Printf("\n\nWord was not found.\n")
	oracle[0] = oracle[0]
	return
}

/*	Construction of the factor oracle automaton for a word p.

	@param p pattern to be added
	@param supply supply map
	@return oracle built oracle
*/
func oracleOnLine(p string)(oracle map[int]map[uint8]int) {
	//m := len(p)
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

/*	Adds one letter to the oracle.

	@param oracle oracle to add letter to
	@param p pattern (not whole, contained in oracle)
	@param o letter to be added
	@param supply supply map
*/
func oracleAddLetter(oracle map[int]map[uint8]int, supply []int, orP string, o uint8)(oracleToReturn map[int]map[uint8]int, orPToReturn string) { 
	m := len(orP)
	var s int
	createNewState(m + 1, oracle)
	//createNewState(m, oracle)
	createTransition(m, o, m + 1, oracle)
	k := supply[m] //0 (nil) the first time
	for k > -1 && getTransition(k,o, oracle) == 0 {
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


// 	Follows some automaton functions.
//	Automaton states are stored in map[int]map[uint8]int:
//		- for each initial state(key) there is a 'value': set of unique characters(keywords) with their destination states (values).
//		- lets assume, that state 0 is always the inital state of the automaton

/*	Automaton function for adding a new state.
	@param state state number to add state for
	@param oracle oracle to add state to
	@return oracle oracle with added state
*/
func createNewState(state int, oracle map[int]map[uint8]int) {
	emptyMap := make(map[uint8]int)
	oracle[state] = emptyMap
	fmt.Printf("\n State %d was created", state)
}

/* 	Automaton function for adding a transition from 'state' over 'letter' to 'end' state.
	@usage createTransition(fromSate, overChar, toState, oracle)
	@return oracle oracle with added transition
*/
func createTransition(fromState int, overChar uint8, toState int, oracle map[int]map[uint8]int) {
	stateMap := oracle[fromState]
	stateMap[overChar]= toState
	oracle[fromState] = stateMap
	fmt.Printf("\n σ(%d,%c)=%d;",fromState,overChar,toState)
}

/*	Returns a 'toState' from state 'fromState' over char 'overChar' :).
	@return toState state for the desired transition
*/
func getTransition(fromState int, overChar uint8, oracle map[int]map[uint8]int)(toState int) {
	stateMap := oracle[fromState]
	toState = stateMap[overChar]
	return toState
}

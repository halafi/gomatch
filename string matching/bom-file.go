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
	
	@param text string/text to be searched in
	@param word word/pattern to be serached for
*/  
func bom(t, p string) {
	oracle := oracleOnLine(p)
	fmt.Printf("\n\nWord was not found.\n")
	oracle[0] = oracle[0]
	return
}

func oracleOnLine(p string)(oracle map[int]map[uint8]int) {
	oracle = make(map[int]map[uint8]int)
	addState(-1, oracle) //used for beginning of the automaton
	addState(0, oracle)
	addTransition(-1, 'e', 0, oracle) //transition from the beginning to the initial state
	
	for i := 0; i < len(p); i++ {
		oracleAddLetter(p[i], oracle, p)
	}
	return oracle
}

func oracleAddLetter(letter uint8, oracle map[int]map[uint8]int, p string) {
	addState(len(p), oracle)
	addTransition(len(p)-1, letter, len(p), oracle)
	k := //getTransition, to be continued.
}


/*
	Automaton functions
*/
func addState(state int, oracle map[int]map[uint8]int) {
	emptyMap := make(map[uint8]int)
	oracle[state] = emptyMap
}

//Transition from begin 'state' over 'letter' to 'end' state
func addTransition(begin int, letter uint8, end int, oracle map[int]map[uint8]int) {
	stateMap := oracle[begin]
	stateMap[letter]=end
	oracle[begin] = stateMap
}
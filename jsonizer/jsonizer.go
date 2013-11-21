package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time"; /*"regexp"*/)

const debugMode bool = true

func main() {
	patternsFile, textFile := loadFiles() //reads text files and converts them into string
	matchDefs := strings.Split(patternsFile, "\n") // match number + what needs to match - definition
	regex, word := process(matchDefs) //we get all regexes and all words
	
	lineMatch := make(map[int]string) // line number + matches we found with repetitions i.g.: lineMathc[0] = "IP IP word word", which will get compared against all the matches
	//lineMatch tbDeleted if not needed
	fmt.Printf("\nJSONIZER\n--------\n")
	if debugMode==true { 
		fmt.Printf("\nThere are %d matches defined in 'patterns.txt':\n",len(matchDefs))
		for i := range matchDefs {
			fmt.Printf("Match %d: %s\n", i, matchDefs[i])
		}
		fmt.Printf("\n\nLETS DO SOME SEARCHING:\n\n")
	}
	startTime := time.Now()
	ahoCor, isFinalForWords, s := buildAc(word)
	lines := strings.Split(textFile, "\n")
	//SEARCHING
	for n := range lines {
		currentLine := strings.Split(lines[n], " ")
		for w:= range currentLine {
			currentState := 0 //state of automaton
			//hledani slov
			for wPos := range currentLine[w] {
				for getTransition(currentState, currentLine[w][wPos], ahoCor) == -1 && s[currentState] != -1 {
					currentState = s[currentState]
				}
				if getTransition(currentState,currentLine[w][wPos], ahoCor) != -1 {
					currentState = getTransition(currentState, currentLine[w][wPos], ahoCor)
				} else {
					currentState = 0
				}
				_, ok := isFinalForWords[currentState] 
				if ok {
					for i := range isFinalForWords[currentState] {
						if word[isFinalForWords[currentState][i]] == getWord(wPos-len(word[isFinalForWords[currentState][i]])+1, wPos, currentLine[w]) {
							//  word[isFinalForWords[currentState][i]] is MATCHED WORD
							/*if debugMode==true {fmt.Printf("Occurence at line %d, %s = %s\n", line, word[isFinalForWords[currentState][i]], word[isFinalForWords[current][i]])}*/
							lineMatch[n] = lineMatch[n] + " " + word[isFinalForWords[currentState][i]]
						}
					}
				}
			}
			//hledani regexu
			for k := range currentLine[w] { //pro kazdy regex check match
				k=k
				for l := range regex {
					if regex[l] == currentLine[w] { //matched,err := regexp.MatchString((string(patFile)), "50")
						fmt.Printf("pmatch")
					}
				}
			}
		}
	}
	
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nElapsed %f secs\n", elapsed.Seconds())
	for i:=0; i<len(lineMatch); i++ {
		fmt.Printf("\n Line %d: %s", i, lineMatch[i])
	}
	return
}

/**
	Takes lines of pattern's definitions and separates them into array's of words and regular expressions.
	Also return's all words given in 'lines'.
*/
func process(lines []string)(allRegexes []string, allWords []string) {
	allRegexes = make([]string, 0)
	allWords = make([]string, 0)
	for i := range lines {
		line := strings.Split(lines[i], " ")
		for j := range line {
			if line[j][0] == '<' {
				allRegexes = addWord(allRegexes, getWord(0, len(line[j])-1, line[j]))
			}
			if line[j][0] == '{' {
				allWords = addWord(allWords, getWord(1, len(line[j])-3, line[j]))
			}
		}
	}
	return allRegexes, allWords
}

/**
	Functions that builds Aho Corasick automaton.
*/
func buildAc(p []string) (acToReturn map[int]map[uint8]int, f map[int][]int, s []int) {
	acTrie, stateIsTerminal, f := constructTrie(p)
	s = make([]int, len(stateIsTerminal)) //supply function
	i := 0 //root of acTrie
	acToReturn = acTrie
	s[i] = -1

	for current := 1; current < len(stateIsTerminal); current++ {
		o, parent := getParent(current, acTrie)
		down := s[parent]
		for stateExists(down, acToReturn) && getTransition(down, o, acToReturn) == -1 {
			down = s[down]
		}
		if stateExists(down, acToReturn) {
			s[current] = getTransition(down, o, acToReturn)
			if stateIsTerminal[s[current]] == true {
				stateIsTerminal[current] = true
				f[current] = arrayUnion(f[current], f[s[current]]) //F(Current) <- F(Current) union F(S(Current))
			}
		} else {
			s[current] = i //initial state?
		}
	}
	return acToReturn, f, s
}

/**
	Function that constructs Trie as an automaton for a set of reversed & trimmed strings.
	
	@return 'trie' built prefix tree
	@return 'stateIsTerminal' array of all states and boolean values of their terminality
	@return 'f' map with keys of pattern indexes and values - arrays of p[i] terminal states
*/
func constructTrie (p []string) (trie map[int]map[uint8]int, stateIsTerminal []bool, f map[int][]int) {
	trie = make(map[int]map[uint8]int)
	stateIsTerminal = make([]bool, 1)
	f = make(map[int][]int) 
	state := 1
	createNewState(0, trie)
	for i:=0; i<len(p); i++ {
		current := 0
		j := 0
		for j < len(p[i]) && getTransition(current, p[i][j], trie) != -1 {
			current = getTransition(current, p[i][j], trie)
			j++
		}
		for j < len(p[i]) {
			stateIsTerminal = boolArrayCapUp(stateIsTerminal)
			createNewState(state, trie)
			stateIsTerminal[state] = false
			createTransition(current, p[i][j], state, trie)
			current = state
			j++
			state++
		}
		if stateIsTerminal[current] {
			newArray := intArrayCapUp(f[current])
			newArray[len(newArray)-1] = i
			f[current] = newArray
		} else {
			stateIsTerminal[current] = true
			f[current] = []int {i}
		}
	}
	return trie, stateIsTerminal, f
}

/**
	Reads IO.
*/
func loadFiles() (pFileString, tFileString string){
	pFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil {
		log.Fatal(err)
	}
	tFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	pFileString = string(pFile)
	tFileString = string(tFile)
	return pFileString, tFileString
}

/**
	Returns 'true' if array of int's 's' contains int 'e', 'false' otherwise.
	
	@author Mostafa http://stackoverflow.com/a/10485970
*/
func contains(s []int, e int) bool {
    for _, a := range s {
		if a == e {
			return true
		}
	}
    return false
}

/*******************          String functions          *******************/
/**
	Check's if word 'w 'exist in array of strings 's', if not - add's it.
	Returns 's' containing word 'w'.
*/
func addWord(s []string, w string) (output []string) {
	for i := range s {
		if s[i] == w {
			return s
		}
	}
	s = stringArrayCapUp(s)
	s[len(s)-1] = w
	return s
}

/**
	Function that returns word found in text 't' at position range 'begin' to 'end'.
*/
func getWord(begin, end int, t string) string {
	for end >= len(t) {
		return ""
	}
	d := make([]uint8, end-begin+1)
	for j, i := 0, begin; i <= end; i, j = i+1, j+1 {
		d[j] = t[i]
	}
	return string(d)
}

/*******************   Array size allocation functions  *******************/
/**
	Dynamically increases an array size of int's by 1.
*/
func intArrayCapUp (old []int)(new []int) {
	new = make([]int, cap(old)+1)
	copy(new, old)  //copy(dst,src)
	old = new
	return new
}

/**
	Dynamically increases an array size of bool's by 1.
*/
func boolArrayCapUp (old []bool)(new []bool) {
	new = make([]bool, cap(old)+1)
	copy(new, old)
	old = new
	return new
}

/**
	Dynamically increases an array size of string's by 1.
*/
func stringArrayCapUp (old []string)(new []string) {
	new = make([]string, cap(old)+1)
	copy(new, old)  //copy(dst,src)
	old = new
	return new
}

/**
	Concats two arrays of int's into one.
*/
func arrayUnion (to, from []int) (concat []int) {
	concat = to
	for i := range(from) {
		if (!contains(concat, from[i])) {
			concat = intArrayCapUp(concat)
			concat[len(concat)-1] = from[i]
		}
	}
	return concat
}

/*******************          Automaton functions          *******************/
/**
	Function that finds the first previous state of a state and returns it. 
	Used for trie where there is only one parent.
	@param 'at' automaton
*/
func getParent(state int, at map[int]map[uint8]int) (uint8, int) {
	for beginState, transitions := range at {
		for c, endState := range transitions {
			if endState == state {
				return c, beginState
			}
		}
	}
	return 0, 0 //unreachable
}

/**
	Automaton function for creating a new state 'state'.
	@param 'at' automaton
*/
func createNewState(state int, at map[int]map[uint8]int) {
	at[state] = make(map[uint8]int)
}

/**
 	Creates a transition for function σ(state,letter) = end.
	@param 'at' automaton
*/
func createTransition(fromState int, overChar uint8, toState int, at map[int]map[uint8]int) {
	at[fromState][overChar]= toState
}

/**
	Returns ending state for transition σ(fromState,overChar), '-1' if there is none.
	@param 'at' automaton
*/
func getTransition(fromState int, overChar uint8, at map[int]map[uint8]int)(toState int) {
	if (!stateExists(fromState, at)) {
		return -1
	}
	toState, ok := at[fromState][overChar]
	if (ok == false) {
		return -1	
	}
	return toState
}

/**
	Checks if state 'state' exists. Returns 'true' if it does, 'false' otherwise.
	@param 'at' automaton
*/
func stateExists(state int, at map[int]map[uint8]int)bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}
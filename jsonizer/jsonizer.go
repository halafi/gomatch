package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time")

/** 
	User defined.
	
	@true prints various extra stuff out, but slows down the execution
	@false will be quick and quiet
*/
const debugMode bool = false

/**
 	THE BIG ONE.
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
	matchLines := strings.Split(string(patFile), "\n")
	regexesPerLine, wordsPerLine, p := separate(matchLines) //declaration
	fmt.Printf("\nJSONIZER v.00000000002\n----------------------\n----------------------\n")
	if debugMode==true { 
		fmt.Printf("\nThere are %d matches defined in 'patterns.txt':\n",len(matchLines))
		for i := range matchLines {
			fmt.Printf("Match %d: %s\n", i, matchLines[i])
			fmt.Printf("Regex:")
			for j:= range regexesPerLine[i] {
				fmt.Printf("%s", regexesPerLine[i][j])
			}
			fmt.Printf(", Words:")
			for j:= range wordsPerLine[i] {
				fmt.Printf("%s", wordsPerLine[i][j])
			}
			fmt.Println()
		}
		fmt.Printf("\nI shall search for them in log file 'text.txt' that is %d chars long!\n\n",len(textFile))
	}
	t := string(textFile)
	startTime := time.Now()
	//
	// SEARCHING
	//
	occurences := make(map[int][]int)
	ac, f, s := buildAc(p)
	if debugMode==true {
		fmt.Printf("\n\nLETS DO IT:\n\n")
	}
	textLines := strings.Split(t, "\n")

	for line := range textLines { //prochazeni po radcich
		t = textLines[line]
		current := 0
		for pos := 0; pos < len(t); pos++ {
			if debugMode==true {
				fmt.Printf("Position: %d, we read: %c", pos, t[pos])
			}
			for getTransition(current, t[pos], ac) == -1 && s[current] != -1 {
				current = s[current]
			}
			if getTransition(current, t[pos], ac) != -1 {
				current = getTransition(current, t[pos], ac)
				if debugMode==true {
					fmt.Printf(" (Continue) \n")
				}
			} else {
				current = 0
				if debugMode==true {
					fmt.Printf(" (FAIL) \n")
				}
			}
			_, ok := f[current]
			if ok {
				for i := range f[current] {
					if p[f[current][i]] == getWord(pos-len(p[f[current][i]])+1, pos, t) { //check for word match
						if debugMode==true {
							fmt.Printf("Occurence at position %d, %q = %q\n", pos-len(p[f[current][i]])+1, p[f[current][i]], p[f[current][i]])
						}
						newOccurences := intArrayCapUp(occurences[f[current][i]])
						occurences[f[current][i]] = newOccurences
						occurences[f[current][i]][len(newOccurences)-1] = pos-len(p[f[current][i]])+1
					}
				}
			}
		}
	}
	
	//old output
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nElapsed %f secs\n", elapsed.Seconds())
	for key, value := range occurences { //prints all occurences of each pattern (if there was at least one)
		fmt.Printf("\nThere were %d occurences for word: %q at positions: ",len(value), p[key])
		for i := range value {
			fmt.Printf("%d", value[i])
			if i != len(value) - 1 {
				fmt.Printf(", ")
			}
		}
		fmt.Printf(".")
	}
	return
}

/**
	Takes lines of pattern's definitions and separates them into array's of words and regular expressions.
	Also return's all words given in 'lines'.
*/
func separate(lines []string)(regexesPerLine map[int][]string, wordsPerLine map[int][]string, allWords []string) {
	regexesPerLine = make(map[int][]string)
	wordsPerLine = make(map[int][]string)
	allWords = make([]string, 0)
	for i := range lines {
		line := strings.Split(lines[i], " ")
		for j := range line {
			if line[j][0] == '<' {
				currentRegexes := regexesPerLine[i]
				currentRegexes = stringArrayCapUp(currentRegexes)
				currentRegexes[len(currentRegexes)-1] = line[j]
				regexesPerLine[i] = currentRegexes
			}
			if line[j][0] == '{' {
				currentWords := wordsPerLine[i]
				currentWords = stringArrayCapUp(currentWords)
				currentWords[len(currentWords)-1] = getWord(1, len(line[j])-3, line[j])
				allWords = addWord(allWords, currentWords[len(currentWords)-1])
				wordsPerLine[i] = currentWords
			}
		}
	}
	return regexesPerLine, wordsPerLine, allWords
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
	if debugMode==true {
		fmt.Printf("\n\nAC construction: \n")
	}
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
				if debugMode==true {
					fmt.Printf(" f[%d] set to: ", current)
					for i := range f[current] {
						fmt.Printf("%d\n", f[current][i])
					}
				}
			}
		} else {
			s[current] = i //initial state?
		}
	}
	if debugMode==true {
		fmt.Printf("\nsupply function: \n")
		for i:= range s {
			fmt.Printf("\ns[%d]=%d", i, s[i])
		}
		fmt.Printf("\n\n")
		for i,j := range f {
			fmt.Printf("f[%d]=", i)
			for k := range j {
				fmt.Printf("%d\n", j[k])
			}
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
			f[current] = newArray //F(Current) <- F(Current) union {i}
			if debugMode==true {
				fmt.Printf(" and %d", i)
			}
		} else {
			stateIsTerminal[current] = true
			f[current] = []int {i}  //F(Current) <- {i}
			if debugMode==true {
				fmt.Printf("\n%d is terminal for word number %d", current, i) 
			}
		}
	}
	return trie, stateIsTerminal, f
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
	if debugMode==true {
		fmt.Printf("\ncreated state %d", state)
	}
}

/**
 	Creates a transition for function σ(state,letter) = end.
	@param 'at' automaton
*/
func createTransition(fromState int, overChar uint8, toState int, at map[int]map[uint8]int) {
	at[fromState][overChar]= toState
	if debugMode==true {
		fmt.Printf("\n    σ(%d,%c)=%d;",fromState,overChar,toState)
	}
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
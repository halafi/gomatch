package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time"; "regexp"; "os"; "strconv")

func main() {
	startTime := time.Now()
	//READ IO
	pFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil {
		log.Fatal(err)
	}
	tFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	tokFile, err := ioutil.ReadFile("tokens.txt")
	if err != nil {
		log.Fatal(err)
	}
	tokenFile, patternsFile, textFile := string(tokFile),string(pFile),string(tFile)
	//Preprocessing
	matchesPerLine := make(map[int][]string)
	patterns := make([]string, 0) //patterns
	lines := strings.Split(patternsFile, "\r\n")
	for i := range lines {
		if len(lines[i]) != 0 {
			patterns = addWord(patterns, lines[i])
		}
	}
	//Print some stuff out
	fmt.Printf("\nJSONIZER v.0.2 \n-----------------------\n")
	//Pattern matching
	ac, f, s := buildAc(patterns)
	lines = strings.Split(textFile, "\r\n")
	for n := range lines {
		words := strings.Split(lines[n], " ")
		current := 0
		for wordIndex := range words { //for each word on a single line
			for getTransition(current, words[wordIndex], ac) == -1 && len(getTransitionTokens(current, ac)) == 0 && s[current] != -1 {
				current = s[current]
			}
			passableTokens := make([]string, 0)
			if len(getTransitionTokens(current, ac)) > 0 { //we can move in 'ac' by some regex
				tokens := getTransitionTokens(current, ac)
				for r := range tokens {
					tokenToMatch := getWord(1, len(tokens[r])-2, tokens[r])
					tokenToMatchSplit := strings.Split(tokenToMatch, ":")
					if len(tokenToMatchSplit) == 2 { //CASE 1: token defined as i.e. <IP:ipAdresa>, output ipAdresa = ...
						regex := regexp.MustCompile(getToken(tokenFile, tokenToMatchSplit[0]))
						if regex.MatchString(words[wordIndex]) { //we got one match
							passableTokens = addWord(passableTokens, tokens[r])
						}
					} else if len(tokenToMatchSplit) == 1 { //CASE 2: token defined as token only, i.e.: <IP>, output IP = ...
						regex := regexp.MustCompile(getToken(tokenFile, tokenToMatch))
						if regex.MatchString(words[wordIndex]) { //we got one match
							passableTokens = addWord(passableTokens, tokens[r])
						}
					} else {
						log.Fatal("Problem in token definition: <"+tokenToMatch+"> use only <TOKEN> or <TOKEN:name>.")
					}
				}
				if len(passableTokens) > 1 {
					log.Fatal("we can match multiple tokens for one word")
				} else if len(passableTokens) == 1{
					current = getTransition(current, passableTokens[0], ac)
				}
			}
			if len(passableTokens) == 0 && getTransition(current, words[wordIndex], ac) != -1 { //move state by a specific word
				current = getTransition(current, words[wordIndex], ac)
			}
			_, ok := f[current]
			if ok { //if(one or more matches)
				for i := range f[current] { // for each pattern that ends at 'current' state
					if isMatch(lines[n], patterns[f[current][i]], tokenFile) { //if(current pattern[i] is match for current line)
						currentMatchText := matchString(lines[n], patterns[f[current][i]], tokenFile)
						if len(currentMatchText) > 1 { //CASE of regex matches (needs to print tokens)
							matchesPerLine[n] = addWord(matchesPerLine[n], strconv.Itoa(f[current][i]+1)+", {"+currentMatchText+"}") 
						} else { //CASE of only word matches
							matchesPerLine[n] = addWord(matchesPerLine[n], strconv.Itoa(f[current][i]+1)) 
						}
					}
				}
			}
		}
	}
	//writing output to a file output.txt
	path := "output.txt"
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for n := range lines { 
		if len(matchesPerLine[n]) == 0 {
			_, err := file.WriteString("NO_MATCH\r\n")
			if err != nil {
				log.Fatal(err)
			}
		} else if len(matchesPerLine[n]) == 1 { //one match found
			_, err := file.WriteString("MATCH + ["+matchesPerLine[n][0]+"]\r\n")
			if err != nil {
				log.Fatal(err)
			}
		} else { //multiple matches found
			longest := 0
			longestOutput := ""
			for j := range matchesPerLine[n] {
				currentMatch := strings.Split(matchesPerLine[n][j], ",")
				if len(currentMatch) > longest {
					longest = len(currentMatch)
					longestOutput = matchesPerLine[n][j]
				}
			}
			_, err := file.WriteString("MATCH + ["+longestOutput+"]\r\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nElapsed %f secs\n", elapsed.Seconds())
	return
}

/**
	Deep assert for pattern and a text line
*/
func isMatch(logLine string, pattern string, tokenFile string) bool {
	logLineWords := strings.Split(logLine, " ")
	patternWords := strings.Split(pattern, " ")
	if len(logLineWords) != len(patternWords) {
		return false
	} else {
		for i := range logLineWords {
			if logLineWords[i] != patternWords[i] {
				tokenToMatch := getWord(1, len(patternWords[i])-2, patternWords[i])
				tokenToMatchSplit := strings.Split(tokenToMatch, ":")
				if len(tokenToMatchSplit) == 2 { 
					regex := regexp.MustCompile(getToken(tokenFile, tokenToMatchSplit[0]))
					if !regex.MatchString(logLineWords[i]) {
						return false
					}
				} else if len(tokenToMatchSplit) == 1 {
					regex := regexp.MustCompile(getToken(tokenFile, tokenToMatch))
					if !regex.MatchString(logLineWords[i]) {
						return false
					}
				}
			}
		}
	}
	return true
}

/**
	Returns match string for output writing.
*/
func matchString(logLine string, pattern string, tokenFile string) string {
	logLineWords := strings.Split(logLine, " ")
	patternWords := strings.Split(pattern, " ")
	output := ""
	tokens := false
	for i := range logLineWords {
		if logLineWords[i] != patternWords[i] {
			tokenToMatch := getWord(1, len(patternWords[i])-2, patternWords[i])
			tokenToMatchSplit := strings.Split(tokenToMatch, ":")
			if len(tokenToMatchSplit) == 2 { 
				regex := regexp.MustCompile(getToken(tokenFile, tokenToMatchSplit[0]))
				if regex.MatchString(logLineWords[i]) {
					if tokens {
						output = output + ", "
					}
					output = output + tokenToMatch +" = "+logLineWords[i]
					tokens = true
				}
			} else if len(tokenToMatchSplit) == 1 {
				regex := regexp.MustCompile(getToken(tokenFile, tokenToMatch))
				if regex.MatchString(logLineWords[i]) {
					if tokens {
						output = output + ", "
					}
					output = output + tokenToMatch +" = "+logLineWords[i]
					tokens = true
				}
			}
		}
	}
	return output
}

/*******************            AC functions          *******************/
/**
        Functions that builds the Aho Corasick automaton.
*/
func buildAc(p []string) (acToReturn map[int]map[string]int, f map[int][]int, s []int) {
        acTrie, stateIsTerminal, f := constructTrie(p)
		fmt.Println()
        s = make([]int, len(stateIsTerminal))
        i := 0
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
                        s[current] = i
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
func constructTrie (p []string) (trie map[int]map[string]int, stateIsTerminal []bool, f map[int][]int) {
        trie = make(map[int]map[string]int)
        stateIsTerminal = make([]bool, 1)
        f = make(map[int][]int) 
        state := 1
        createNewState(0, trie)
        for i := range p {
				words := strings.Split(p[i], " ")
                current := 0
                j := 0 //word index
                for j < len(words) && getTransition(current, words[j], trie) != -1 {
                        current = getTransition(current, words[j], trie)
                        j++
                }
                for j < len(words) {
                        stateIsTerminal = boolArrayCapUp(stateIsTerminal)
                        createNewState(state, trie)
                        stateIsTerminal[state] = false
                        createTransition(current, words[j], state, trie)
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
						//fmt.Printf("%d is terminal for pattern number %d\n", current, i) 
                }
        }
        return trie, stateIsTerminal, f
}

/*******************          String functions          *******************/
/**
	Returns regex for desired token in string 'tokenFile'.
*/
func getToken(tokenFile, wanted string) string {
	tokenLines := strings.Split(tokenFile, "\r\n")
	for n := range tokenLines {
		token := strings.Split(tokenLines[n], " ")
		if len(token) == 2 && token[0] == wanted {
			return token[1]
		}
	}
	log.Fatal("NO TOKEN DEFINITION in tokens.txt FOR: ", wanted)
	return ""
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
        @author 'Walter' http://stackoverflow.com/a/10043083
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

/**
        Function that computes minimal length string in a set of strings.
*/
func computeMinLength(p []string) (lmin int){
        lmin = len(p[0])
        for i:=1; i<len(p); i++ {
                if (len(p[i])<lmin) {
                        lmin = len(p[i])
                }
        }
        return lmin
}

/*******************            Array functions            *******************/
/**
	Functions 'type'ArrayCapUp dynamically increases an 'type's array 
	maximum size by 1. (copy(dst,src))
*/
func byteArrayCapUp (old []byte)(new []byte) {
	new = make([]byte, cap(old)+1)
	copy(new, old)  
	return new
}

func intArrayCapUp (old []int)(new []int) {
	new = make([]int, cap(old)+1)
	copy(new, old) 
	return new
}


func boolArrayCapUp (old []bool)(new []bool) {
	new = make([]bool, cap(old)+1)
	copy(new, old)
	return new
}

func stringArrayCapUp (old []string)(new []string) {
	new = make([]string, cap(old)+1)
	copy(new, old)  //copy(dst,src)
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

/*******************          Automaton functions          *******************/
/**
	Function that finds the first previous state of a state and returns it. 
	Used for trie where there is only one parent, otherwise won't work.
	@param 'at' automaton
*/
func getParent(state int, at map[int]map[string]int) (string, int) {
	for beginState, transitions := range at {
		for c, endState := range transitions {
			if endState == state {
				return c, beginState
			}
		}
	}
	return " ", 0 //unreachable
}

/**
	Returns all transitioning tokens for a state 'state'.
	@param 'at' automaton
*/
func getTransitionTokens(state int, at map[int]map[string]int) ([]string) {
	toReturn := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] == '<' && s[len(s)-1] == '>' {
			toReturn = addWord(toReturn, s)
		}
	}
	return toReturn
}

/**
	Automaton function for creating a new state 'state'.
	@param 'at' automaton
*/
func createNewState(state int, at map[int]map[string]int) {
	at[state] = make(map[string]int)
	//fmt.Printf("State %d created.\n", state)
}

/**
 	Creates a transition for function σ(state,string) = end.
	@param 'at' automaton
*/
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	at[fromState][overString]= toState
	//fmt.Printf("σ(%d,%s) = %d\n", fromState, overString, toState)
}

/**
	Returns ending state for transition σ(fromState,overString), '-1' if there is none.
	@param 'at' automaton
*/
func getTransition(fromState int, overString string, at map[int]map[string]int)(toState int) {
	if (!stateExists(fromState, at)) {
		return -1
	}
	toState, ok := at[fromState][overString]
	if (ok == false) {
		return -1	
	}
	return toState
}

/**
	Checks if state 'state' exists. Returns 'true' if it does, 'false' otherwise.
	@param 'at' automaton
*/
func stateExists(state int, at map[int]map[string]int)bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}
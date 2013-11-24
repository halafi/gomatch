package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time"; "regexp"; "os"; "strconv")

func main() {
	startTime := time.Now()
	outputText := make(map[int][]string) //line number + text
	//#1 - Reads Input files
	pFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil { log.Fatal(err) }
	tFile, err := ioutil.ReadFile("text.txt")
	if err != nil { log.Fatal(err) }
	tokFile, err := ioutil.ReadFile("tokens.txt")
	if err != nil { log.Fatal(err) }
	tokenFile, patternsFile, textFile := string(tokFile),string(pFile),string(tFile)
	//#2 - Preprocessing
	p := make([]string, 0)
	matches := make(map[int][]string)
	lines := strings.Split(patternsFile, "\r\n")
	for i := range lines {
		line := strings.Split(lines[i], " ")
		for j := range line {
			if line[j][0] == '{' {
				p = addWord(p, getWord(1, len(line[j])-2, line[j]))
			}
		}
		matches[i] = strings.Split(lines[i], " ")
	}
	//#3 - Print some stuff out
	fmt.Printf("\nJSONIZER\n-----------------------\nPatterns.txt\n")
	for i,arrayOfS := range matches {
		fmt.Printf("Match %d: ", i+1)
		for j := range arrayOfS {
			fmt.Printf("%s ", arrayOfS[j])
		}
		fmt.Println()
	}
	//#4 - searching for matches
	lines = strings.Split(textFile, "\n")
	for n := range lines { 
		currentLine := strings.Split(lines[n], " ")
		for m := range matches {
			firstSubmatch := true
			wordPos := 0
			for mW := 0; mW < len(matches[m]) && mW < len(currentLine); mW++ {
				if mW == 0 { //set default output text for current line
					outputText[n] = stringArrayCapUp(outputText[n])
					outputText[n][len(outputText[n])-1] = ""
				}
				if matches[m][mW][0] == '<' { //REGEX_MATCHING
					tokenToMatch := getWord(1, len(matches[m][mW])-2, matches[m][mW])
					regex := regexp.MustCompile(getToken(tokenFile, tokenToMatch))
					if  !regex.MatchString(currentLine[mW]) { //NO_MATCH
						outputText[n][len(outputText[n])-1] = "" 
						break
					} else if mW == 0 && len(matches[m]) == 1{ //one match containing one pattern ended
						outputText[n][len(outputText[n])-1] = strconv.Itoa(m+1)+", {"+tokenToMatch+"="+currentLine[mW]+"}"
					} else if mW < len(matches[m])-1 && len(matches[m]) > 1 { //longer pattern got submatch
						if firstSubmatch {
							outputText[n][len(outputText[n])-1] = strconv.Itoa(m+1)+", {"+tokenToMatch+"="+currentLine[mW]
							firstSubmatch = false
						} else {
							outputText[n][len(outputText[n])-1] = outputText[n][len(outputText[n])-1]+", "+tokenToMatch+"="+currentLine[mW]
						}
					} else if mW == len(matches[m])-1 && len(matches[m]) > 1 { //longer pattern: submatch at the end of match_def
						if firstSubmatch {
							outputText[n][len(outputText[n])-1] = strconv.Itoa(m+1)+", {"+tokenToMatch+"="+currentLine[mW]+"}"
							firstSubmatch = false
						} else {
							outputText[n][len(outputText[n])-1] = outputText[n][len(outputText[n])-1]+", "+tokenToMatch+"="+currentLine[mW]+"}"
						}
					}
				} else if matches[m][mW][0] == '{' { //WORD_MATCHING
					if len(p) > 0 {
						wordOccurences := searchSBOM(p, lines[n])
						//wordOccurences := searchAC(p, lines[n]) //je rychlejsi pro mene kratkych patternu
						wordToMatch := getWord(1, len(matches[m][mW])-2, matches[m][mW])
						if !contains(wordOccurences[wordToMatch],wordPos) { //NO_MATCH
							outputText[n][len(outputText[n])-1] = ""
							break
						} else if mW == 0 && len(matches[m]) == 1 { //one match containing one word ended
							outputText[n][len(outputText[n])-1] = strconv.Itoa(m+1)
						} else if mW == 0 && len(matches[m]) > 1{ //match starts with word that we found
							outputText[n][len(outputText[n])-1] = strconv.Itoa(m+1)
							firstSubmatch = false
						} else if mW == len(matches[m])-1 && len(matches[m]) > 1 { //longer pattern: submatch at the end of match_def{
							if firstSubmatch == false{
								outputText[n][len(outputText[n])-1] = outputText[n][len(outputText[n])-1] + "}"
							}
						}// otherwise no need to do anything
					}
				} else if matches[m][mW][0] == '_' { //ANY_MATCHING
					if mW == len(matches[m])-1 {
						if firstSubmatch {
							outputText[n][len(outputText[n])-1] = strconv.Itoa(m+1)
						} else { //closes text if it needs to
							outputText[n][len(outputText[n])-1] = outputText[n][len(outputText[n])-1] + "}"
						}
					}
				} else {
					outputText[n][len(outputText[n])-1] = "ERROR_MATCHING: "+string(matches[m][mW])
					break
				}
				wordPos = wordPos + len(currentLine[mW]) +1
			}
		}
	}
	//#5 - writing output to a file output.txt
	path := "output.txt"
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for n := range lines {
		noMatch := true
		printedMatch := false
		for j := range outputText[n] {
			if outputText[n][j] != "" && !strings.Contains(outputText[n][j], "ERROR") {
				noMatch = false
				if printedMatch == false {
					_, err := file.WriteString("MATCH + ")
					if err != nil {
						log.Fatal(err)
					}
					printedMatch = true
				}
				_, err := file.WriteString("["+outputText[n][j]+"]")
				if err != nil {
					log.Fatal(err)
				}
			} else if strings.Contains(outputText[n][j], "ERROR") {
				_, err := file.WriteString(outputText[n][j] + " LINE " + strconv.Itoa(n))
				if err != nil {
					log.Fatal(err)
				}
				noMatch = false
				break
			}
		}
		if noMatch == true { //all strings were empty, print NO_MATCH
			_, err := file.WriteString("NO_MATCH\r\n")
			if err != nil {
			log.Fatal(err)
			}
		} else {
			_, err := file.WriteString("\r\n")
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	elapsed := time.Since(startTime)
	fmt.Printf("\n\nElapsed %f secs\n", elapsed.Seconds())
	return
}

/******************* Multiple string matching functions *******************/
func searchAC(p []string, t string) map[string][]int {
	occurences := make(map[string][]int)
	ac, f, s := buildAc(p)
	current := 0
	for pos := 0; pos < len(t); pos++ {
		for getTransition(current, t[pos], ac) == -1 && s[current] != -1 {
			current = s[current]
		}
		if getTransition(current, t[pos], ac) != -1 {
			current = getTransition(current, t[pos], ac)
		} else {
			current = 0
		}
		_, ok := f[current]
		if ok {
			for i := range f[current] {
				if p[f[current][i]] == getWord(pos-len(p[f[current][i]])+1, pos, t) {
					occurences[p[f[current][i]]] = intArrayCapUp(occurences[p[f[current][i]]])
					occurences[p[f[current][i]]][len(occurences[p[f[current][i]]])-1] = pos-len(p[f[current][i]])+1
				}
			}
		}
	}
	return occurences
}

func searchSBOM(p []string, t string) map[string][]int {
	occurences := make(map[string][]int)
	lmin := computeMinLength(p)
	or, f := buildOracleMultiple(reverseAll(trimToLength(p, lmin)))
	pos := 0
	for pos <= len(t) - lmin {
			current := 0
			j := lmin
			for j >= 1 && stateExists(current, or) {
					current = getTransition(current, t[pos+j-1], or)
					j--
			}
			word := getWord(pos, pos+lmin-1, t)
			if stateExists(current, or) && j == 0 && strings.HasPrefix(word, getCommonPrefix(p, f[current], lmin)) {
					for i := range f[current] {
							if p[f[current][i]] == getWord(pos, pos-1+len(p[f[current][i]]), t) {
									occurences[p[f[current][i]]] = intArrayCapUp(occurences[p[f[current][i]]])
									occurences[p[f[current][i]]][len(occurences[p[f[current][i]]])-1] = pos
							}
					}
					j = 0
			}
			pos = pos + j + 1
	}
	return occurences
}

/**
        Function that builds factor oracle used by sbom.
*/
func buildOracleMultiple (p []string) (orToReturn map[int]map[uint8]int, f map[int][]int) {
        orTrie, stateIsTerminal, f := constructTrie(p)
        s := make([]int, len(stateIsTerminal)) //supply function
        i := 0 //root of trie
        orToReturn = orTrie
        s[i] = -1
        for current := 1; current < len(stateIsTerminal); current++ {
                o, parent := getParent(current, orTrie)
                down := s[parent]
                for stateExists(down, orToReturn) && getTransition(down, o, orToReturn) == -1 {
                        createTransition(down, o, current, orToReturn)
                        down = s[down]
                }
                if stateExists(down, orToReturn) {
                        s[current] = getTransition(down, o, orToReturn)
                } else {
                        s[current] = i
                }
        }
        return orToReturn, f
}

/**
	Functions that builds Aho Corasick automaton.
*/
func buildAc(p []string) (acToReturn map[int]map[uint8]int, f map[int][]int, s []int) {
	acTrie, stateIsTerminal, f := constructTrie(p)
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
	Returns regex for desired token in string 'tokenFile'.
*/
func getToken(tokenFile, wanted string) string {
	tokenLines := strings.Split(tokenFile, "\r\n")
	for n := range tokenLines {
		token := strings.Split(tokenLines[n], " ")
		if token[0] == wanted {
			return token[1]
		}
	}
	log.Fatal("NO TOKEN DEFINITION in tokens.txt FOR: ", wanted)
	return ""
}

/**
        Returns a prefix size 'lmin' for one string 'p' of first index found in 'f'.
        It is not needed to compare all the strings from 'p' indexed in 'f',
        thanks to the konwledge of 'lmin'.
*/
func getCommonPrefix(p []string, f []int, lmin int) string {
        r := []rune(p[f[0]])
        newR := make([]rune, lmin)
        for j := 0; j < lmin; j++ {
                newR[j] = r[j]
        }
        return string(newR)
}

/**
        Function that takes a set of strings 'p' and their wanted 'length'
        and then trims each string in that set to have desired 'length'.
*/
func trimToLength(p []string, length int) (trimmedP []string) {
        trimmedP = make([]string, len(p))
        for i := range p {
                r := []rune(p[i])
                newR := make([]rune, length)
                for j := 0; j < length; j++ {
                        newR[j] = r[j]
                }
                trimmedP[i]=string(newR)
        }
        return trimmedP
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

/*******************   Array size allocation functions  *******************/
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
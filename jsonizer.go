package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time"; "regexp"; "os"; "strconv"; /*"encoding/json"*/)

//EXPERIMENTAL: change this, if you wish to search in a file that has words separated by something different than spaces
const wordSeparator = " "

func main() {
	fmt.Printf("\nJSONIZER v.0.4 \n-----------------------\n")
	
	fmt.Println("Step 1: Processing input.\n")
	logFilePath, outputPath, tokensFilePath, patternsFilePath := "text.txt", "output.json", "tokens.txt", "patterns.txt"
	tokenDefinitions, patternsString, logString := parseFile(tokensFilePath), parseFile(patternsFilePath), parseFile(logFilePath)
	logLines, patterns := splitFileString(logString), splitFileString(patternsString)
	
	fmt.Println("Step 2: Matching: "+logFilePath+".")
	startTime := time.Now()
	trie, f := constructTrie(tokenDefinitions, patterns)
	
	matchesPerLine := make(map[int][]string)
	for n := range logLines { //for each log line
		words, current := strings.Split(logLines[n], wordSeparator), 0
		for w := range words { //for each word
			transitionTokens := getTransitionTokens(current, trie)
			validTokens := make([]string, 0)
			if getTransition(current, words[w], trie) != -1 { //we can move by a word: 'words[w]'
				current = getTransition(current, words[w], trie)
			} else if len(transitionTokens) > 0 { //we can move by some regex
				for t := range transitionTokens { //for each token leading from 'current' state
					tokenWithoutBrackets := getWord(1, len(transitionTokens[t])-2, transitionTokens[t])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
						case 2: { //token defined as i.e. IP:ipAdresa
							regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBracketsSplit[0]))
							if regex.MatchString(words[w]) {
								validTokens = addWord(validTokens, transitionTokens[t])
							}
						}
						case 1: { //token defined as token only, i.e.: IP
							regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBrackets))
							if regex.MatchString(words[w]) {
								validTokens = addWord(validTokens, transitionTokens[t])
							}
						}
						default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
				if len(validTokens) > 1 { //we got i.e. string "user" that matches both <WORD> and i.e. <USERNAME>...
					log.Fatal("Multiple acceptable tokens for one word at log line: "+strconv.Itoa(n+1)+", position: "+strconv.Itoa(w+1)+".")	
				} else if len(validTokens) == 1 { //we can move exactly by one regex/token
					current = getTransition(current, validTokens[0], trie)
				}
			} else {
				//fmt.Printf("FAIL: %s \n", words[w])
				break
			}
			_, isCurrentFinalState := f[current]
			if isCurrentFinalState {
				for i := range f[current] { // for each pattern that ends at 'current' state
					if isMatch(logLines[n], patterns[f[current][i]], tokenDefinitions) {
						outputText := matchString(logLines[n], patterns[f[current][i]], tokenDefinitions)
						if len(outputText) > 1 { //CASE of regex matches (needs to print tokens)
							matchesPerLine[n] = addWord(matchesPerLine[n], strconv.Itoa(f[current][i]+1)+", {"+outputText+"}") 
						} else { //CASE of only word matches
							matchesPerLine[n] = addWord(matchesPerLine[n], strconv.Itoa(f[current][i]+1)) 
						}
					}
				}
			}
		}
	}
	
	elapsedMatch := time.Since(startTime)
	fmt.Printf("        Elapsed %f secs\n\n", elapsedMatch.Seconds())
	
	fmt.Println("Step 3: Writing output to file: "+outputPath)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for n := range logLines { 
		switch len(matchesPerLine[n]) {
			case 0: { //no match for current line
				_, err := file.WriteString("NO_MATCH"+"\r\n")
				if err != nil {
					log.Fatal(err)
				}
			}
			case 1: { //one match for current line
				_, err := file.WriteString("MATCH + ["+matchesPerLine[n][0]+"]"+"\r\n")
				if err != nil {
					log.Fatal(err)
				}
			}
			default: { //multiple matches for current line
				longest := 0
				longestOutput := ""
				for j := range matchesPerLine[n] {
					currentMatch := strings.Split(matchesPerLine[n][j], ",")
					if len(currentMatch) > longest {
						longest = len(currentMatch)
						longestOutput = matchesPerLine[n][j]
					}
				}
				_, err := file.WriteString("MATCH + ["+longestOutput+"]"+"\r\n")
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
	fmt.Println("\n\nAll Done.")
	return
}

/**
	Deep assert for pattern and a text line
*/
func isMatch(logLine string, pattern string, tokenFile string) bool {
	logLineWords := strings.Split(logLine, wordSeparator)
	patternWords := strings.Split(pattern, wordSeparator)
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
	logLineWords := strings.Split(logLine, wordSeparator)
	patternWords := strings.Split(pattern, wordSeparator)
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
					output = output + tokenToMatchSplit[1] +" = "+logLineWords[i]
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

/**
        Function that constructs Trie as an automaton for a set of reversed & trimmed strings.
        
        @return 'trie' built prefix tree
        @return 'stateIsTerminal' array of all states and boolean values of their terminality
        @return 'f' map with keys of pattern indexes and values - arrays of p[i] terminal states
*/
func constructTrie (tokenDefinitions string, p []string) (trie map[int]map[string]int, f map[int][]int) {
        trie = make(map[int]map[string]int)
        stateIsTerminal := make([]bool, 1)
        f = make(map[int][]int) 
        state := 1
        createNewState(0, trie)
        for i := range p {
				words := strings.Split(p[i], wordSeparator)
                current := 0
                j := 0
                for j < len(words) && getTransition(current, words[j], trie) != -1 {
                        current = getTransition(current, words[j], trie)
                        j++
                }
                for j < len(words) {
                        stateIsTerminal = boolArrayCapUp(stateIsTerminal)
                        createNewState(state, trie)
                        stateIsTerminal[state] = false
						
						if len(getTransitionWords(current, trie)) > 0 && words[j][0] == '<' && words[j][len(words[j])-1] == '>' { //check for conflict when adding regex transition
							transitions := getTransitionWords(current, trie)
							tokenWithoutBrackets := getWord(1, len(words[j])-2, words[j])
							tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
							for w := range transitions {
								switch len(tokenWithoutBracketsSplit) {
									case 2: { //token defined as i.e. IP:ipAdresa
										regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBracketsSplit[0]))
										if regex.MatchString(transitions[w]) {
											log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitions[w])	
										}
									}
									case 1: { //token defined as token only, i.e.: IP
										regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBrackets))
										if regex.MatchString(transitions[w]) {
											log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitions[w])	
										}
									}
									default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
								}
							}
						} else if len(getTransitionTokens(current, trie)) > 0 && words[j][0] != '<' && words[j][len(words[j])-1] != '>' { //check for conflict when adding word
							tokens := getTransitionTokens(current, trie)
							for t := range tokens {
								tokenWithoutBrackets := getWord(1, len(tokens[t])-2, tokens[t])
								tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
								switch len(tokenWithoutBracketsSplit) {
									case 2: { //token defined as i.e. IP:ipAdresa
										regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBracketsSplit[0]))
										if regex.MatchString(words[j]) {
											log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+tokens[t])	
										}
									}
									case 1: { //token defined as token only, i.e.: IP
										regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBrackets))
										if regex.MatchString(words[j]) {
											log.Fatal("Conflict in patterns definition, token "+tokens[t]+" matches word "+words[j])	
										}
									}
									default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
								}
							}
						}
						
                        createTransition(current, words[j], state, trie)
                        current = state
                        j++
                        state++
                }
                if stateIsTerminal[current] {
                        f[current] = intArrayCapUp(f[current])
                        f[current][len(f[current])-1] = i
                } else {
                        stateIsTerminal[current] = true
                        f[current] = []int {i}
                }
        }
        return trie, f
}

/*******************          File functions           *******************/
/**
	Simple file reader that return string content of the file.
*/
func parseFile(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileString := string(file)
	return fileString
}

/**
	Function that parses file into single lines (array of strings).
*/
func splitFileString(fileString string) []string {
	splitFile := make([]string, 1) 
	splitFile[0] = fileString //default single pattern, no line break
	if strings.Contains(fileString, "\r\n") { //CR+LF
		splitFile = strings.Split(fileString, "\r\n")
	} else if strings.Contains(fileString, "\n") { //LF
		splitFile = strings.Split(fileString, "\n")
	} else if strings.Contains(fileString, "\r") { //CR
		splitFile = strings.Split(fileString, "\r")
	}
	return splitFile
}

/*******************          String functions          *******************/
/**
	Returns regex for desired token in string 'tokenFile'.
*/
func getToken(tokenFile, wanted string) string {
	tokenLines := strings.Split(tokenFile, "\r\n")
	for n := range tokenLines {
		token := strings.Split(tokenLines[n], wordSeparator)
		if len(token) == 2 && token[0] == wanted {
			return token[1]
		}
	}
	log.Fatal("NO TOKEN DEFINITION in tokens.txt FOR: ", wanted)
	return "" //unreachable
}

/**
	Check's if word 'w 'exist in array of strings 's', if not - add's it.
	Returns 's' containing word 'w'.
*/
func addWord(s []string, w string) []string {
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

/*******************            Array functions            *******************/
/**
	Increases an array of byte's maximum size by 1.
*/
func byteArrayCapUp (old []byte)(new []byte) {
	new = make([]byte, cap(old)+1)
	copy(new, old)  
	return new
}

/**
	Increases an array of int's maximum size by 1.
*/
func intArrayCapUp (old []int)(new []int) {
	new = make([]int, cap(old)+1)
	copy(new, old) 
	return new
}

/**
	Increases an array of bool's maximum size by 1.
*/
func boolArrayCapUp (old []bool)(new []bool) {
	new = make([]bool, cap(old)+1)
	copy(new, old)
	return new
}

/**
	Increases an array of string's maximum size by 1.
*/
func stringArrayCapUp (old []string)(new []string) {
	new = make([]string, cap(old)+1)
	copy(new, old)  //copy(dst,src)
	return new
}

/**
	Concats two arrays of types int into one.
*/
func arrayUnion (to, from []int) []int {
	concat := to
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
	Function that finds the previous state of a state and returns it. 
	Used for trie where there is only one parent, otherwise won't work.
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
	Returns all transitioning tokens for a state 'state' in automaton 'at'.
*/
func getTransitionTokens(state int, at map[int]map[string]int) []string {
	transitionTokens := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] == '<' && s[len(s)-1] == '>' {
			transitionTokens = addWord(transitionTokens, s)
		}
	}
	return transitionTokens
}

/**
	Returns all transitioning words for a state 'state' in automaton 'at'.
*/
func getTransitionWords(state int, at map[int]map[string]int) []string {
	transitions := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] != '<' && s[len(s)-1] != '>' {
			transitions = addWord(transitions, s)
		}
	}
	return transitions
}

/**
	Automaton function for creating a new state 'state' in automaton 'at'.
*/
func createNewState(state int, at map[int]map[string]int) {
	at[state] = make(map[string]int)
}

/**
 	Creates a transition for function 'σ(fromState,overString) = toState' in automaton 'at'.
*/
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	at[fromState][overString]= toState
	//fmt.Printf("σ(%d,%s) = %d\n", fromState, overString, toState)
}

/**
	Returns ending state - 'toState' for transition 'σ(fromState,overString)' in automaton 'at'.
	State -1 is returned if there is none.
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
	Checks if state 'state' exists in automaton 'at'.
	Returns 'true' if it does, 'false' otherwise.
*/
func stateExists(state int, at map[int]map[string]int) bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}
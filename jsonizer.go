package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time"; "regexp"; "os"; "strconv"; "encoding/json")

//EXPERIMENTAL: change this, if you wish to search in a file that has words separated by something different than spaces
const wordSeparator = " "

type Match struct { //struct used for storing found matches
	Type string //Match name
	Body map[string]string //Body[token] = token value
}

func main() {
	/*args := os.Args
	if (len(args) < 2 || len(args) > 2) {
		log.Fatal("Wrong ammount of arguments. \nTwo arguments required:  'inputLogPath patternsFilePath'.")
	}
	pattern := args[1]
	s := args[2]*/
	logFilePath, outputPath, tokensFilePath, patternsFilePath := "text.txt", "output.json", "Tokens", "patterns.txt"
	tokenDefinitions, patternsString, logString := parseFile(tokensFilePath), parseFile(patternsFilePath), parseFile(logFilePath)
	logLines, patterns := splitFileString(logString), splitFileString(patternsString)
	
	fmt.Printf("\nMatching...")
	startTime := time.Now()
	trie, f := constructTrie(tokenDefinitions, patterns)
	matchPerLine := make([]Match, len(logLines))
	
	for n := range logLines {
		words, current := strings.Split(logLines[n], wordSeparator), 0
		for w := range words {
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
				break
			}
			_, isCurrentFinalState := f[current]
			if isCurrentFinalState {
				patternSplit := strings.Split(patterns[f[current]], "##")
				body := matchString(logLines[n], patternSplit[1], tokenDefinitions)
				if len(body) > 1 { //CASE of regex matches (needs to print tokens)
					matchPerLine[n] = Match{patternSplit[0], body}
				} else { //CASE of only word matches
					matchPerLine[n] = Match{patternSplit[0], body}
				}
			}
		}
	}
	elapsedMatch := time.Since(startTime)
	fmt.Printf(", Elapsed: %fs", elapsedMatch.Seconds())
	//JSON Output
	indent := "  " //u can use anything like "   " or "\t" ..., nothing else should need to be changed (determines how is the output printed)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	_, err = file.WriteString("[")
		if err != nil {
			log.Fatal(err)
		}
	first := true
	for n := range matchPerLine {
		if matchPerLine[n].Type	!= "" {
			if !first {
				_, err = file.WriteString(",")
				if err != nil {
					log.Fatal(err)
				}
			}
			if first {
				first = false
			}
			b, err := json.MarshalIndent(matchPerLine[n], indent+indent, indent)
			if err != nil {
				log.Fatal(err)
			}
			_, err = file.WriteString("\r\n"+indent+"{\r\n"+indent+indent+"\"Event\": " + string(b)+"\r\n"+indent+"}")
			if err != nil {
				log.Fatal(err)
			}
			
		}
	}
	_, err = file.WriteString("\r\n]")
	if err != nil {
		log.Fatal(err)
	}
	return
}

/**
	Returns match string for output writing.
*/
func matchString(logLine string, pattern string, tokenFile string) map[string]string {
	logLineWords := strings.Split(logLine, wordSeparator)
	patternWords := strings.Split(pattern, wordSeparator)
	output := make(map[string]string)
	for i := range logLineWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := getWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			if len(tokenWithoutBracketsSplit) == 2 {
				regex := regexp.MustCompile(getToken(tokenFile, tokenWithoutBracketsSplit[0]))
				if regex.MatchString(logLineWords[i]) {
					output[tokenWithoutBracketsSplit[1]] = logLineWords[i]
				}
			} else if len(tokenWithoutBracketsSplit) == 1 {
				regex := regexp.MustCompile(getToken(tokenFile, tokenWithoutBrackets))
				if regex.MatchString(logLineWords[i]) {
					output[tokenWithoutBrackets] = logLineWords[i]
				}
			} else {
				log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
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
func constructTrie (tokenDefinitions string, p []string) (trie map[int]map[string]int, f map[int]int) {
        trie = make(map[int]map[string]int)
        stateIsTerminal := make([]bool, 1)
        f = make(map[int]int) 
        state := 1
        createNewState(0, trie)
        for i := range p {
				patternsNameSplit := strings.Split(p[i], "##")
				if len(patternsNameSplit) != 2 {
					log.Fatal("Error with pattern number ",i+1," name, use [NAME##<token> word ...]. Name cannot be empty.")
				}
				if len(patternsNameSplit[0]) == 0 {
					log.Fatal("Error with pattern number ",i+1,": name cannot be empty.")
				}
				
				words := strings.Split(patternsNameSplit[1], wordSeparator)
				
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
						log.Fatal("Duplicate pattern definition.")
                } else {
                        stateIsTerminal[current] = true
                        f[current] = i
                }
        }
        return trie, f
}

/*******************          File functions           *******************/
/**
	Simple file reader that returns string content of the file.
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
	Function that parses string file into single lines (array of strings).
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
	tokenLines := splitFileString(tokenFile)
	for n := range tokenLines {
		token := strings.Split(tokenLines[n], " ")
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
	Concats two arrays of type int into one.
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
	Returns 'true' if an array of int's 's' contains int 'e', 'false' otherwise.
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
	Function that finds the previous state (parent) of a state and returns it. 
	Works only in automaton/tree with single parents.
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
	Returns all transitioning tokens (without words) for a state 'state' in automaton 'at'.
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
	Returns all transitioning words (without tokens) for a state 'state' in automaton 'at'.
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
	Creates a new state 'state' in automaton 'at'.
*/
func createNewState(state int, at map[int]map[string]int) {
	at[state] = make(map[string]int)
}

/**
 	Creates a transition 'σ(fromState,overString) = toState' in automaton 'at'.
*/
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	at[fromState][overString]= toState
}

/**
	Returns ending state - 'toState' for transition 'σ(fromState,overString)' in automaton 'at'.
	State -1 is returned if there is no transition.
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

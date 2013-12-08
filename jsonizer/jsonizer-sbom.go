package main
import ("fmt"; "log"; "strings"; "io/ioutil"; "time"; "regexp"; "os"; "strconv")

func main() {
	startTime := time.Now()
	//Preprocessing, reading input files, printing some stuff out
	tokensFile, err := ioutil.ReadFile("tokens.txt")
	if err != nil {
		log.Fatal(err)
	}
	patternsFile, err := ioutil.ReadFile("patterns.txt")
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := ioutil.ReadFile("text.txt")
	if err != nil {
		log.Fatal(err)
	}
	tokensString, patternsString, logString := string(tokensFile), string(patternsFile), string(logFile)
	matchesPerLine := make(map[int][]string)
	patterns := make([]string, 0) //patterns
	lines := strings.Split(patternsString, "\r\n")
	for i := range lines {
		if len(lines[i]) != 0 {
			patterns = addWord(patterns, lines[i])
		}
	}
	fmt.Printf("\nJSONIZER v.0.3 \n-----------------------\n")
	//searching
	lmin := computeMinLength(patterns)
	/*fmt.Println(lmin)
	fmt.Println(trimToLength(patterns, lmin))*/
	or, f := buildOracleMultiple(reverseAll(trimToLength(patterns, lmin)))
	
	lines = strings.Split(logString, "\r\n")
	for n := range lines {
		words := strings.Split(lines[n], " ")
		current := 0
		wordIndex := 0
		or, f, tokensString, current, wordIndex = or, f, tokensString, current, wordIndex
		for wordIndex <= len(words) - lmin {
			fmt.Println("once")
			current := 0
			j := lmin	
			for j >= 1 && stateExists(current, or) {
				moved := false
				passableTokens := make([]string, 0)
				if len(getTransitionTokens(current, or)) > 0 { //we can move by some regex
					tokens := getTransitionTokens(current, or)
					for r := range tokens {
						tokenToMatch := getWord(1, len(tokens[r])-2, tokens[r])
						tokenToMatchSplit := strings.Split(tokenToMatch, ":")
						if len(tokenToMatchSplit) == 2 { //CASE 1: token defined as i.e. <IP:ipAdresa>, output ipAdresa = ...
							regex := regexp.MustCompile(getToken(tokensString, tokenToMatchSplit[0]))
							if regex.MatchString(words[wordIndex]) { //we got one match
								passableTokens = addWord(passableTokens, tokens[r])
								fmt.Printf("%s matches %s\n", tokens[r], words[wordIndex])
							}
						} else if len(tokenToMatchSplit) == 1 { //CASE 2: token defined as token only, i.e.: <IP>, output IP = ...
							regex := regexp.MustCompile(getToken(tokensString, tokenToMatch))
							if regex.MatchString(words[wordIndex]) { //we got one match
								passableTokens = addWord(passableTokens, tokens[r])
								fmt.Printf("%s matches %s\n", tokens[r], words[wordIndex])
							}
						} else {
							log.Fatal("Problem in token definition: <"+tokenToMatch+"> use only <TOKEN> or <TOKEN:name>.")
						}
					}
					if len(passableTokens) > 1 { //we got i.e. string "user" that matches both <WORD> and i.e. <USERNAME>...
						log.Fatal("We can match multiple tokens for one word. This shouldn't happen.")	
					} else if len(passableTokens) == 1 {
						current = getTransition(current, passableTokens[0], or)
						fmt.Println(current)
						j--
						moved = true
					}
				} 
				if len(passableTokens) == 0 && moved == false  { // move by some word
					current = getTransition(current, words[wordIndex+j-1], or)
					fmt.Printf("moved by a word to %d, j = %d\n", current, j)
					j--
				}
			}
			
			word := buildWord(wordIndex, wordIndex+lmin, words) //buildWord(from,to)
			fmt.Println(strings.HasPrefix(word, getCommonPrefix(patterns, f[current], lmin)))
			if stateExists(current, or) && j == 0 && strings.HasPrefix(word, getCommonPrefix(patterns, f[current], lmin)) {
				fmt.Println("possible occurence\n")
				for i := range f[current] { //some occurence, verify all possible against the text
					if patterns[f[current][i]] == lines[n] {
						fmt.Println("occurence")
						currentMatchText := matchString(lines[n], patterns[f[current][i]], tokensString)
						if len(currentMatchText) > 1 { //CASE of regex matches (needs to print tokens)
							matchesPerLine[n] = addWord(matchesPerLine[n], strconv.Itoa(f[current][i]+1)+", {"+currentMatchText+"}") 
						} else { //CASE of only word matches
							matchesPerLine[n] = addWord(matchesPerLine[n], strconv.Itoa(f[current][i]+1)) 
						}
					}
				}
				j = 0
			}
			wordIndex = wordIndex + j + 1
		}
	}
	//Output printing
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

/*******************            SBOM functions          *******************/
/**
	Function that builds factor oracle used by sbom.
*/
func buildOracleMultiple (p []string) (orToReturn map[int]map[string]int, f map[int][]int) {
	orTrie, stateIsTerminal, f := constructTrie(p)
	s := make([]int, len(stateIsTerminal))
	i := 0
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

/**
	Function that takes a set of strings 'p' and their wanted number of words - 'length'
	and then trims each string in that set to have desired ammount of words.
*/
func trimToLength(p []string, length int) (trimmedP []string) {
	trimmedP = make([]string, len(p))
	for i := range p {
		psplit := strings.Split(p[i], " ")
		for j := 0; j < length; j++ {
			trimmedP[i] = trimmedP[i] + psplit[j]
			if j != length-1 {
				trimmedP[i] = trimmedP[i] + " "
			}
		}
	}
	return trimmedP
}

/**        
	Function that takes an array of patterns (words separated by " ") and reverses them.
*/
func reverseAll(s []string) (reversed []string) {
	reversed = make([]string, len(s))
	for i := 0; i < len(s); i++ {
		reversed[i] = reverse(s[i])
	}
	return reversed
}

/**        
	Function that takes a string of words separated by " " and reverses the order of words.
*/
func reverse(s string) (reversed string) {
	reversed = ""
	split := strings.Split(s, " ")
	for i, j := 0, len(split)-1; i < len(split); j, i = j -1, i+1 {
		reversed = reversed + split[j]
		if i != len(split)-1 {
			reversed = reversed + " "
		}
	}
	return reversed
}

/**
    Function that computes minimal length of a pattern in a set of patterns.
*/
func computeMinLength(p []string) (lmin int){
	pSplit := strings.Split(p[0], " ")
	lmin = len(pSplit)
	for i:=1; i < len(p); i++ {
		pSplit = strings.Split(p[i], " ")
		if (len(pSplit) < lmin) {
			lmin = len(pSplit)
		}
	}
	return lmin
}

/**
	Returns a prefix size 'lmin' for one string 'p' of first index found in 'f'.
	It is not needed to compare all the strings from 'p' indexed in 'f',
	thanks to the konwledge of 'lmin'.
*/
func getCommonPrefix(p []string, f []int, lmin int) string {
	split := strings.Split(p[f[0]], " ")
	prefix := ""
	for j := 0; j < lmin; j++ {
		prefix = prefix + split[j]
		if j != lmin - 1 {
			prefix = prefix + " "
		}
	}
	return prefix
}

func buildWord(from , to int, words []string) string{
	out := ""
	for i := from; i < to; i++ {
		out = out + words[i]
		if i != to-1 {
			out = out + " "
		}
	}
	return out
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

/*******************            Array functions            *******************/
/**
	Functions 'type'ArrayCapUp increases an array of 'type's maximum size by 1.
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
	Concats two arrays of types int into one.
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
	Returns all transitioning tokens for a state 'state' in automaton 'at'.
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
	fmt.Printf("σ(%d,%s) = %d\n", fromState, overString, toState)
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
func stateExists(state int, at map[int]map[string]int)bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}
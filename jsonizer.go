package main
import (
	"code.google.com/p/go.crypto/ssh/terminal"
	"fmt";
	"log";
	"strings";
	"io/ioutil";
	"regexp";
	"os";
	"strconv";
	"encoding/json"; //json conversion
)

const (
	indent = "   " //determines JSON output indent (formatting), you can use anything like three spaces(default) or "\t"...
	wordSeparator = " " //change this, if you wish to search in a file that has words separated by something different than spaces
) 

type Match struct {
	Type string
	Body map[string]string
}

/**
	Function main performs reading input files, matching of Log file and printing JSON output.
*/
func main() {
	var logLines []string
	if ! terminal.IsTerminal(0) { //Stdin not empty
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		logLines = lineSplit(string(bytes))
	} else {
		logFilePath := "text.txt"
		logString := fileToString(logFilePath)
		logLines = lineSplit(logString)
	}
    
	patternsFilePath := "Patterns"
	
	//outputPath := "output.json"
	tokensFilePath := "Tokens"
	tokenDefinitions, patternsString := fileToString(tokensFilePath), fileToString(patternsFilePath)
	patterns := lineSplit(patternsString)

	trie, finalFor, stateIsTerminal := constructPrefixTree(tokenDefinitions, patterns)
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
					tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
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
			if stateIsTerminal[current] {
				patternSplit := strings.Split(patterns[finalFor[current]], "##")
				body := getMatchBody(logLines[n], patternSplit[1], tokenDefinitions)
				if len(body) > 1 { //CASE of regex matches (needs to print tokens)
					matchPerLine[n] = Match{patternSplit[0], body}
				} else { //CASE of only word matches
					matchPerLine[n] = Match{patternSplit[0], body}
				}
			}
		}
	}
	//Output
	fmt.Printf("%s", getJSON(matchPerLine))
	return
}

/**
	Function that constructs prefix tree/automaton for a set strings.
	
	@param 'p' an array of patterns with words/tokens separated by single spaces each
	@param 'tokenDefinitions' string of text lines with "token_name(space)regex"
	@return 'trie' built prefix tree
	@return 'finalFor' an array of int's, for each final state it is equal to one pattern number
	@return 'stateIsTerminal' an array of boolean values for each state representing if its terminal
*/
func constructPrefixTree (tokenDefinitions string, p []string) (trie map[int]map[string]int, finalFor []int, stateIsTerminal []bool) {
	trie = make(map[int]map[string]int)
	stateIsTerminal = make([]bool, 1)
	finalFor = make([]int, 1) 
	state := 1
	for i := range p {
		if p[i]=="" { //EOF, extra line break
			break
		}
		patternsNameSplit := strings.Split(p[i], "##") //separate pattern name from its definition
		if len(patternsNameSplit) != 2 {
			log.Fatal("Error with pattern number ",i+1," name, use [NAME##<token> word ...].")
		}
		if len(patternsNameSplit[0]) == 0 {
			log.Fatal("Error with pattern number ",i+1,": name cannot be empty.")
		}
		if len(patternsNameSplit[1]) == 0 {
			log.Fatal("Error with pattern number ",i+1,": pattern cannot be empty.")
		}
		words := strings.Split(patternsNameSplit[1], wordSeparator)
		current, j := 0, 0
		for j < len(words) && getTransition(current, words[j], trie) != -1 {
			current = getTransition(current, words[j], trie)
			j++
		}
		for j < len(words) {
			stateIsTerminal = boolArrayCapUp(stateIsTerminal)
			finalFor = intArrayCapUp(finalFor)
			stateIsTerminal[state] = false
			if len(getTransitionWords(current, trie)) > 0 && words[j][0] == '<' && words[j][len(words[j])-1] == '>' { //conflict check when adding regex transition
				transitionWords := getTransitionWords(current, trie)
				for w := range transitionWords {
					tokenWithoutBrackets := cutWord(1, len(words[j])-2, words[j])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
						case 2: {
							regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBracketsSplit[0]))
							if regex.MatchString(transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitionWords[w]+".")	
							}
						}
						case 1: {
							regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBrackets))
							if regex.MatchString(transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitionWords[w]+".")	
							}
						}
						default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
					}
				}
			} else if len(getTransitionTokens(current, trie)) > 0 && words[j][0] != '<' && words[j][len(words[j])-1] != '>' { //conflict check when adding word
				transitionTokens := getTransitionTokens(current, trie)
				for t := range transitionTokens {
					tokenWithoutBrackets := cutWord(1, len(transitionTokens[t])-2, transitionTokens[t])
					tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
					switch len(tokenWithoutBracketsSplit) {
						case 2: {
							regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBracketsSplit[0]))
							if regex.MatchString(words[j]) {
								log.Fatal("Conflict in patterns definition, token "+transitionTokens[t]+" matches word "+words[j]+".")	
							}
						}
						case 1: {
							regex := regexp.MustCompile(getToken(tokenDefinitions, tokenWithoutBrackets))
							if regex.MatchString(words[j]) {
								log.Fatal("Conflict in patterns definition, token "+transitionTokens[t]+" matches word "+words[j]+".")	
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
			log.Fatal("Duplicate pattern definition detected, pattern number: ",i+1,".")
		} else {
			stateIsTerminal[current] = true
			finalFor[current] = i
		}
	}
	return trie, finalFor, stateIsTerminal
}

/*****************************************************************************/
/*********************          File functions           *********************/

/**
	Simple file reader that returns a string content of the file.
	
	@param 'filePath' path to the given file
	@return string with contents of file at 'filePath'
*/
func fileToString(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return  string(file)
}

/**
	Function that parses a mutli-line string into single lines (array of strings).
	
	@param 'input' string to be splitted
	@return 'inputSplit' an array of single 'input' lines
*/
func lineSplit(input string) []string {
	inputSplit := make([]string, 1) 
	inputSplit[0] = input //default single pattern, no line break
	if strings.Contains(input, "\r\n") { //CR+LF
		inputSplit = strings.Split(input, "\r\n")
	} else if strings.Contains(input, "\n") { //LF
		inputSplit = strings.Split(input, "\n")
	} else if strings.Contains(input, "\r") { //CR
		inputSplit = strings.Split(input, "\r")
	}
	return inputSplit
}

/*****************************************************************************/
/*********************          String functions          ********************/

/**
	For a set of matches this function returns a string containing JSON data.
	If you wish to change how the output looks you can set custom value for
	constant 'indent'.
	
	@param 'matchPerLine' an array of type Match representing for each log line found match
	@return string JSON formatted output
*/
func getJSON(matchPerLine []Match) string {
	output := "["
	first := true
	for n := range matchPerLine {
		if matchPerLine[n].Type	!= "" {
			if !first {
				output = output + ","
			} else {
				first = false
			}
			b, err := json.MarshalIndent(matchPerLine[n], indent+indent, indent)
			if err != nil {
				log.Fatal(err)
			}
			output = output + "\r\n"+indent+"{\r\n"+indent+indent+"\"Event\": " + string(b)+"\r\n"+indent+"}"
		}
	}
	return output + "\r\n]"
}

/**
	Returns a Match body - string containing found tokens and their given values.
	
	@param 'logLine' logLine is used to get token values
	@param 'pattern' pattern is used to get token names
	@param 'tokenFileString' required token definitions
	@return map of key: string, value: string, where key="token name", value="token value"
*/
func getMatchBody(logLine string, pattern string, tokenFileString string) map[string]string {
	logLineWords := strings.Split(logLine, wordSeparator)
	patternWords := strings.Split(pattern, wordSeparator)
	output := make(map[string]string)
	for i := range logLineWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := cutWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			switch len(tokenWithoutBracketsSplit) {
				case 2: {
					regex := regexp.MustCompile(getToken(tokenFileString, tokenWithoutBracketsSplit[0]))
					if regex.MatchString(logLineWords[i]) {
						output[tokenWithoutBracketsSplit[1]] = logLineWords[i]
					}
				}
				case 1: {
					regex := regexp.MustCompile(getToken(tokenFileString, tokenWithoutBrackets))
					if regex.MatchString(logLineWords[i]) {
						output[tokenWithoutBrackets] = logLineWords[i]
					}
				} 
				default: log.Fatal("Problem in token definition: <"+tokenWithoutBrackets+">, use only <TOKEN> or <TOKEN:name>.")
			}
		}
	}
	return output
}

/**
	Returns a regular expression for a single token in string that contains it.
	
	@param 'tokensFileString' string of single lines "<TOKEN>(space)regex"
	@param 'token' string of wanted regular expression name
	@return string containing the regular expression for 'token'
*/
func getToken(tokensFileString, token string) string {
	tokenFileLines := lineSplit(tokensFileString)
	for n := range tokenFileLines {
		lineSplit := strings.Split(tokenFileLines[n], " ")
		if len(lineSplit) == 2 && lineSplit[0] == token {
			return lineSplit[1]
		}
	}
	log.Fatal("No token definition for: ", token,".")
	return "" //unreachable
}

/**
	Function checks if a word 'word 'exist in an array of strings, if not - it is added.
	
	@param 's' an array of strings
	@param 'word' word to be added
	@return array of strings containing word 'word' and all of the old values
*/
func addWord(s []string, word string) []string {
	for i := range s {
		if s[i] == word {
			return s
		}
	}
	s = stringArrayCapUp(s)
	s[len(s)-1] = word
	return s
}

/**
	Function that for a given word performs a cut so that the new words starts at 'begin' position
	of the old word and ends at 'end' position of the old word.
	
	@param 'begin' begin position
	@param 'end' end position
	@param 'word' word to be cut
	@return string containing some characters from 'word'
*/
func cutWord(begin, end int, word string) string {
	if end >= len(word) {
		return ""
	}
	d := make([]uint8, end-begin+1)
	for j, i := 0, begin; i <= end; i, j = i+1, j+1 {
		d[j] = word[i]
	}
	return string(d)
}

/*****************************************************************************/
/*******************            Array functions            *******************/

/**
	Increases an array of int's size by 1.
	
	@param old old array
	@return new new, bigger array with old values
*/
func intArrayCapUp (old []int)(new []int) {
	new = make([]int, cap(old)+1)
	copy(new, old) 
	return new
}

/**
	Increases an array of bool's size by 1.
	
	@param old old array
	@return new new, bigger array with old values
*/
func boolArrayCapUp (old []bool)(new []bool) {
	new = make([]bool, cap(old)+1)
	copy(new, old)
	return new
}

/**
	Increases an array of string's size by 1.
	
	@param old old array
	@return new new, bigger array with old values
*/
func stringArrayCapUp (old []string)(new []string) {
	new = make([]string, cap(old)+1)
	copy(new, old)
	return new
}

/*****************************************************************************/
/*******************          Automaton functions          *******************/

/**
	Returns all transitioning tokens (without words): only <IP>, <DATE:date> etc.
	for a state 'state' in an automaton 'at' as an array of strings.
	
	@param 'state' state that the transition tokens will be returned for
	@param 'at' automaton where the state is
	@return an array of strings containing all the possible token transitions
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
	Returns all transitioning words (without tokens): only "user", "date:" etc.
	for a state 'state' in an automaton 'at' as an array of strings.
	
	@param 'state' state that the transition words will be returned for
	@param 'at' automaton where the state is
	@return an array of strings containing all the possible word transitions
*/
func getTransitionWords(state int, at map[int]map[string]int) []string {
	transitionWords := make([]string, 0)
	for s, _ := range at[state] {
		if s[0] != '<' && s[len(s)-1] != '>' {
			transitionWords = addWord(transitionWords, s)
		}
	}
	return transitionWords
}

/**
	If there is no state 'fromState', this function creates it, 
	after that transitionion 'σ(fromState,overString) = toState' is created in an automaton 'at'.
	
	@param 'fromState' beginning state of the transition
	@param 'overString' transitioning word
	@param 'toState' ending state of the transition
	@param 'at' automaton where the transition will be created at
*/
func createTransition(fromState int, overString string, toState int, at map[int]map[string]int) {
	if stateExists(fromState, at) {
		at[fromState][overString]= toState
	} else {
		at[fromState] = make(map[string]int)
		at[fromState][overString]= toState
	}
}

/**
	Returns an ending state for transition 'σ(fromState,overString)' in automaton 'at'.
	
	@param 'fromState' state where the transition begins
	@param 'overString' transitioning word 
	@param 'at' automaton where the transition will be returned from
	@return -1 if there is no transition
	@return ending state number if there is a transition
*/
func getTransition(fromState int, overString string, at map[int]map[string]int) int {
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
	Checks if a state 'state' exists in an automaton 'at'.
	
	@param 'state' state to check existence for
	@param 'at' automaton where the state should be checked at
	@return true if it does exist
	@return false otherwise
*/
func stateExists(state int, at map[int]map[string]int) bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}

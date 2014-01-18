package main
import ( //imported Go packages
	"fmt"; //formatted I/O
	"log"; //logging package
	"strings"; //functions to manipulate strings
	"io/ioutil"; //some I/O utility functions
	"regexp"; //regular expression search
	"os"; //platform-independent interface to operating system functionality
	"strconv"; //conversions to and from string representations of basic data types
	"encoding/json"; //encoding and decoding of JSON objects as defined in RFC 4627
	"code.google.com/p/go.crypto/ssh/terminal" //provides support functions for dealing with terminals
)

const ( //USER DEFINED constants
	indent = "   " //determines JSON output indent (formatting), you can use anything like three spaces(default) or "\t"...
	wordSeparator = " " //change this, if you wish to search in a file that has words separated by something different than spaces
	patternsFilePath = "Patterns" //location of Pattern definitions
	tokensFilePath = "Tokens" //location of Token definitions 
)

type Match struct { //structure used for storing matches
	Type string
	Body map[string]string
}

func main() {
	//initial declaration
	var logLines []string
	patterns := lineSplit(fileToString(patternsFilePath))
	tokenDefinitions := fileToString(tokensFilePath)

	//reading of input log
	if ! terminal.IsTerminal(0) { //if there is standard input, read it
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			log.Fatal(err)
		}
		logLines = lineSplit(string(bytes))
	} else { //otherwise check for single filepath argument or fail
		if len(os.Args) == 2 { 
			logFile, err := ioutil.ReadFile(os.Args[1])
			if err != nil {
				log.Fatal(err)
			}
			logLines = lineSplit(string(logFile))
		} else {
			log.Fatal("Invalid program usage: no standard input or filepath argument given. \nSample usage: \n\"cat /var/log/kern.log | jsonizer > output.json\"\n\"jsonizer /var/log/kern.log\"")
		}
	}
	
	matchPerLine := make([]Match, len(logLines))
	trie, finalFor, stateIsTerminal := constructPrefixTree(tokenDefinitions, patterns) //construction of prefix tree
	
	//reading log line by line checking for matches
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
						case 2: { //token + name, i.e. <IP:ipAddress>
							if matchToken(tokenDefinitions, tokenWithoutBracketsSplit[0], words[w]) {
								validTokens = addWord(validTokens, transitionTokens[t])
							}
						}
						case 1: { //token only, i.e.: <IP>
							if matchToken(tokenDefinitions, tokenWithoutBrackets, words[w]) {
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
			if stateIsTerminal[current] && w == len(words)-1 { //we have reached leaf node in prefix tree AND end of log line - got match
				patternSplit := strings.Split(patterns[finalFor[current]], "##")
				body := getMatchBody(logLines[n], patternSplit[1], tokenDefinitions)
				
				if len(body) > 1 { //body with some tokens
					matchPerLine[n] = Match{patternSplit[0], body}
				} else { //empty body
					matchPerLine[n] = Match{patternSplit[0], nil}
				}
			}
		}
	}
	
	//printing results to standard output
	output := getJSON(matchPerLine)
	if output != "[\r\n]" {
		fmt.Printf("%s\r\n", getJSON(matchPerLine))
	}
	return
}

/**
 *	Function that constructs prefix tree/automaton for a set strings.
 *	
 *	@param 'p' an array of patterns with words/tokens separated by single spaces each
 *	@param 'tokenDefinitions' string of text lines with "token_name(space)regex"
 *	@return 'trie' built prefix tree
 *	@return 'finalFor' an array of int's, for each final state it is equal to one pattern number
 *	@return 'stateIsTerminal' an array of boolean values for each state representing if its terminal
 */
func constructPrefixTree (tokenDefinitions string, p []string) (trie map[int]map[string]int, finalFor []int, stateIsTerminal []bool) {
	trie = make(map[int]map[string]int)
	stateIsTerminal = make([]bool, 1)
	finalFor = make([]int, 1) 
	state := 1
	for i := range p {
		if p[i]== "" { //EOF, extra line break
			if i == len(p)-1 {
				break
			}
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
							if matchToken(tokenDefinitions, tokenWithoutBracketsSplit[0], transitionWords[w]) {
								log.Fatal("Conflict in patterns definition, token "+words[j]+" matches word "+transitionWords[w]+".")	
							}
						}
						case 1: {
							if matchToken(tokenDefinitions, tokenWithoutBrackets, transitionWords[w]) {
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
							if matchToken(tokenDefinitions, tokenWithoutBracketsSplit[0], words[j]) {
								log.Fatal("Conflict in patterns definition, token "+transitionTokens[t]+" matches word "+words[j]+".")	
							}
						}
						case 1: {
							if matchToken(tokenDefinitions, tokenWithoutBrackets, words[j]) {
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

/**
 *	Simple matching function, that takes a single token and a word that should match.
 *	
 *	@param 'tokenDefinitions' string of single lines "<TOKEN>(space)regex"
 *	@param 'token' token to match
 *	@param 'word' word to match
 *	@return true if token matches, false otherwise
 */
func matchToken (tokenDefinitions, token, word string) bool {
	regex := regexp.MustCompile(getToken(tokenDefinitions, token))
	if regex.MatchString(word) {
		return true
	} else {
		return false
	}
}

/**
 *	Returns a regular expression for a single token in string that contains it.
 *	
 *	@param 'tokensFileString' string of single lines "<TOKEN>(space)regex"
 *	@param 'token' string of wanted regular expression name
 *	@return string containing the regular expression for 'token'
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

/*****************************************************************************/
/*********************          File functions           *********************/

/**
 *	Simple file reader that returns a string content of the file.
 *	
 *  @param 'filePath' path to the given file
 *	@return string with contents of file at 'filePath'
 */
func fileToString(filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return  string(file)
}

/**
 *	Function that parses a mutli-line string into single lines (array of strings).
 *	
 *	@param 'input' string to be splitted
 *	@return 'inputSplit' an array of single 'input' lines
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
 *	For a set of matches this function returns a string containing JSON data.
 *	If you wish to change how the output looks you can set custom value for
 *	constant 'indent'.
 *	
 *	@param 'matchPerLine' an array of type Match representing for each log line found match
 *	@return string JSON formatted output
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
 *	Returns a Match body - string containing found tokens and their given values.
 *	
 *	@param 'logLine' logLine is used to get token values
 *	@param 'pattern' pattern is used to get token names
 *	@param 'tokenFileString' required token definitions
 *	@return map of key: string, value: string, where key="token name", value="token value"
 */
func getMatchBody(logLine string, pattern string, tokenFileString string) map[string]string {
	logLineWords := strings.Split(logLine, wordSeparator)
	patternWords := strings.Split(pattern, wordSeparator)
	output := make(map[string]string)
	for i := range patternWords {
		if logLineWords[i] != patternWords[i] {
			tokenWithoutBrackets := cutWord(1, len(patternWords[i])-2, patternWords[i])
			tokenWithoutBracketsSplit := strings.Split(tokenWithoutBrackets, ":")
			switch len(tokenWithoutBracketsSplit) {
				case 2: {
					if matchToken(tokenFileString, tokenWithoutBracketsSplit[0], logLineWords[i]) {
						output[tokenWithoutBracketsSplit[1]] = logLineWords[i]
					}
				}
				case 1: {
					if matchToken(tokenFileString, tokenWithoutBrackets, logLineWords[i]) {
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
 *  Function checks if a word 'word 'exist in an array of strings, if not - it is added.
 *	
 *	@param 's' an array of strings
 *	@param 'word' word to be added
 *	@return array of strings containing word 'word' and all of the old values
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
 *	Function that for a given word performs a cut so that the new words starts at 'begin' position
 *	of the old word and ends at 'end' position of the old word.
 * 
 *	@param 'begin' begin position
 *	@param 'end' end position
 *	@param 'word' word to be cut
 *	@return string containing some characters from 'word'
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
 *	Increases an array of int's size by 1.
 *	
 *	@param old old array
 *	@return new new, bigger array with old values
 */
func intArrayCapUp (old []int)(new []int) {
	new = make([]int, cap(old)+1)
	copy(new, old) 
	return new
}

/**
 *	Increases an array of bool's size by 1.
 *	
 *	@param old old array
 *	@return new new, bigger array with old values
 */
func boolArrayCapUp (old []bool)(new []bool) {
	new = make([]bool, cap(old)+1)
	copy(new, old)
	return new
}

/**
 *	Increases an array of string's size by 1.
 *	
 *	@param old old array
 *	@return new new, bigger array with old values
 */
func stringArrayCapUp (old []string)(new []string) {
	new = make([]string, cap(old)+1)
	copy(new, old)
	return new
}

/*****************************************************************************/
/*******************          Automaton functions          *******************/

/**
 *	Returns all transitioning tokens (without words): only <IP>, <DATE:date> etc.
 *	for a state 'state' in an automaton 'at' as an array of strings.
 *	
 *  @param 'state' state that the transition tokens will be returned for
 *	@param 'at' map with stored states and their transitions
 *	@return an array of strings containing all the possible token transitions
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
 *	Returns all transitioning words (without tokens): only "user", "date:" etc.
 *	for a state 'state' in an automaton 'at' as an array of strings.
 *	
 *	@param 'state' state that the transition words will be returned for
 *	@param 'at' map with stored states and their transitions
 *	@return an array of strings containing all the possible word transitions
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
 *	If there is no state 'fromState', this function creates it, 
 *	after that transitionion 'σ(fromState,overString) = toState' is created in an automaton 'at'.
 *	
 *	@param 'fromState' beginning state of the transition
 *	@param 'overString' transitioning word
 *	@param 'toState' ending state of the transition
 *	@param 'at' map with stored states and their transitions
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
 *	Returns an ending state for transition 'σ(fromState,overString)' in automaton 'at'.
 *	
 *	@param 'fromState' state where the transition begins
 *	@param 'overString' transitioning word 
 *	@param 'at' map with stored states and their transitions
 *	@return -1 if there is no transition
 *	@return ending state number if there is a transition
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
 *	Checks if a state 'state' exists in an automaton 'at'.
 *	
 *	@param 'state' state to check existence for
 *	@param 'at' map with stored states and their transitions
 *	@return true if it does exist
 *	@return false otherwise
 */
func stateExists(state int, at map[int]map[string]int) bool {
	_, ok := at[state]
	if (!ok || state == -1 || at[state] == nil) {
		return false
	}
	return true
}

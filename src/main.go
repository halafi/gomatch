package main

import "flag"
import "log"
import "os"

// Command-line flags.
var input = flag.String("i", "/dev/stdin", "Data input stream.")
var output = flag.String("o", "/dev/stdout", "Data output stream.")
var outputFormat = flag.String("f", "json", "Output data format, supported: json, xml, plain.")
var patternsIn = flag.String("p", "./Patterns", "Pattern definitions input.")
var tokensIn = flag.String("t", "./Tokens", "Token definitions input.")

// Function main() performs a few steps: 
func main() {
	flag.Parse()
	// Token reading
	tokenReader := openFile(*tokensIn)
	tokens := make(map[string]string)
	for {
		line, eof := readLine(tokenReader)
		if eof {
			break
		}
		token := checkToken(line)
		if token != "" {
			addToken(token, tokens)
		}
	}
	tree, finalFor, state, i := createNewTrie()
	// Initial pattern reading
	patternReader := openFile(*patternsIn)
	patternsArr := make([]string, 0)
	for {
		line, eof := readLine(patternReader)
		if eof {
			break
		}
		pattern := checkPattern(line)
		if pattern != "" {
			newPatternsArr := make([]string, cap(patternsArr)+1) // array size +1
			copy(newPatternsArr, patternsArr)
			newPatternsArr[len(newPatternsArr)-1] = pattern // add pattern to array of all patterns
			patternsArr = newPatternsArr
			tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i) // add pattern to trie
		}
	}
	patternsFileInfo, err := os.Stat(*patternsIn)
	if err != nil {
		log.Fatal(err)
	}
	lastModified := patternsFileInfo.ModTime()
	
	// Reading of input lines, matching them and writing them to output
	inputReader := openFile(*input)
	outputFile := createFile(*output)
	for {
		// If last mod time for patterns file is different, then read
		// the first line of patterns file and check for a new pattern
		patternsFileInfo, err = os.Stat(*patternsIn)
		if err != nil {
			log.Fatal(err)
		}
		if lastModified != patternsFileInfo.ModTime() {
			patternReader = openFile(*patternsIn)
			line, eof := readLine(patternReader)
			if !eof {
				pattern := checkPattern(line)
				if pattern != "" && !contains(patternsArr, pattern) {
					log.Printf("New event: \"%s\".", pattern)
					newPatternsArr := make([]string, cap(patternsArr)+1) // array size +1
					copy(newPatternsArr, patternsArr)
					newPatternsArr[len(newPatternsArr)-1] = pattern // add pattern to array of all patterns
					patternsArr = newPatternsArr
					tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i) // add pattern to trie
				}
				lastModified = patternsFileInfo.ModTime()
			}
			
		}
		logLine, eof := readLine(inputReader)
		if eof {
			writeFile(outputFile, convertMatch(getMatch(logLine, patternsArr, tokens, tree, finalFor), *outputFormat))
			break
		} else {
			writeFile(outputFile, convertMatch(getMatch(logLine, patternsArr, tokens, tree, finalFor), *outputFormat))
		}
	}
	closeFile(outputFile)
	return
}

// Calls the desired get method and returns its output.
func convertMatch(match Match, output string) string {
	if output=="JSON" || output=="json" {
		return getJSON(match)
	}
	if output=="XML" || output=="xml" {
		return getXML(match)
	}
	log.Fatal("unknown output format: \"", output +"\"")
	return ""
}

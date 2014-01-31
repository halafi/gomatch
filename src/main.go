package main

import "flag"
import "log"
import "os"

// Command-line flags.
var input = flag.String("i", "/dev/stdin", "Log data input.")
var output = flag.String("o", "/dev/stdout", "JSON data output.")
var patternsIn = flag.String("p", "./Patterns", "Pattern definitions input.")
var tokensIn = flag.String("t", "./Tokens", "Token definitions input.")

// Function main() performs a few steps: 
func main() {
	flag.Parse()
	
	tokens := readTokens(*tokensIn)
	patternReader := openFile(*patternsIn)
	patternsArr := make([]string, 0)
	tree, finalFor, state, i := createNewTrie()
	
	// Initial pattern reading
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
	
	// Reading of input lines, matching them and writing them to output.
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
			writeFile(outputFile, getJSON(getMatch(logLine, patternsArr, tokens, tree, finalFor)))
			break
		} else {
			writeFile(outputFile, getJSON(getMatch(logLine, patternsArr, tokens, tree, finalFor)))
		}
	}
	closeFile(outputFile)
	return
}

package main

import "./lib/input/logdata"
import "./lib/input/tokens"
import "./lib/input/patterns"

import "./lib/match"
import "./lib/match/trie"

import "./lib/output/json"

import "flag"

// Command-line flags.
var input = flag.String("i", "/dev/stdin", "Log data input.")
var output = flag.String("o", "/dev/stdout", "JSON data output.")
var patternsIn = flag.String("p", "./Patterns", "Pattern definitions input.")
var tokensIn = flag.String("t", "./Tokens", "Token definitions input.")

// Function main() performs a few steps: 
func main() {
	flag.Parse()
	
	tokens := tokens.ReadTokens(*tokensIn)
	patternReader := patterns.Init(*patternsIn)
	patternsArr := make([]string, 0)
	tree, finalFor, state, i := trie.Init()
	
	// Initial pattern reading
	for {
		pattern, eof := patterns.ReadPattern(patternReader)
		if eof {
			break
		} else {
			newPatternsArr := make([]string, cap(patternsArr)+1) // array size +1
			copy(newPatternsArr, patternsArr)
			newPatternsArr[len(newPatternsArr)-1] = pattern // add pattern to array of all patterns
			patternsArr = newPatternsArr
			tree, finalFor, state, i = trie.AppendPattern(tokens, pattern, tree, finalFor, state, i) // add pattern to trie
		}
	}
	
	// Reading of input lines, matching them and writing them to output.
	inputReader := logdata.Open(*input)
	outputFile := json.CreateOutputFile(*output)
	for {
		logLine, eof := logdata.ReadLine(inputReader)
		if eof {
			json.WriteOutputFile(outputFile, match.GetMatch(logLine, patternsArr, tokens, tree, finalFor))
			break
		} else {
			json.WriteOutputFile(outputFile, match.GetMatch(logLine, patternsArr, tokens, tree, finalFor))
		}
	}
	json.CloseOutputFile(outputFile)

	return
}

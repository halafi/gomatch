package main


import "./lib/input/tokens"
import "./lib/input/patterns"

//import "./lib/input/log_file"
import "./lib/input/log_unixpipe"

import "./lib/output/json"
import "./lib/match"
import "./lib/match/trie"

// Function main() performs a few steps: reading of input, matching and
// priting of JSON output to STDOUT.
func main() {
	/*if len(os.Args) == 2 {
		logLines := file.ReadLog(os.Args[1])
	}*/
	patterns := patterns.ReadPatterns("Patterns")
	tokens := tokens.ReadTokens("Tokens")
	tree, finalFor, stateIsTerminal := trie.ConstructPrefixTree(tokens, patterns)
	
	unixPipeReader := unixpipe.Init()
	for {
		//fmt.Println("try")
		logLine, eof := unixpipe.ReadLine(unixPipeReader) 
		if eof {
			//fmt.Println("eof")
			json.PrintJSON(match.GetMatch(logLine, patterns, tokens, tree, finalFor, stateIsTerminal))
			return
		} else {
			json.PrintJSON(match.GetMatch(logLine, patterns, tokens, tree, finalFor, stateIsTerminal))
		}
	}
	return
}

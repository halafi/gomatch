package main

import "./lib/input/log_file"
import "./lib/input/log_unixpipe"
import "./lib/input/tokens"
import "./lib/input/patterns"
import "./lib/match"
import "./lib/match/trie"
import "./lib/output/json"
import "os"

// Function main() performs a few steps: reading of input, matching and
// priting of JSON output to STDOUT.
func main() {
	tokens := tokens.ReadTokens("Tokens")
	patternsArr := make([]string, 0)
	patternReader := patterns.Init("Patterns")
	tree, finalFor, state, i := trie.Init()
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

	if len(os.Args) == 2 {
		logLines := file.ReadLog(os.Args[1])
		for n := range logLines {
			json.PrintJSON(match.GetMatch(logLines[n], patternsArr, tokens, tree, finalFor))
		}
	} else {
		unixPipeReader := unixpipe.Init()
		for {
			logLine, eof := unixpipe.ReadLine(unixPipeReader)
			if eof {
				json.PrintJSON(match.GetMatch(logLine, patternsArr, tokens, tree, finalFor))
				return
			} else {
				json.PrintJSON(match.GetMatch(logLine, patternsArr, tokens, tree, finalFor))
			}
		}
	}

	return
}

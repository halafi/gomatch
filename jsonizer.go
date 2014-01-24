package main


import "./lib/input/tokens"
import "./lib/input/patterns"

import "./lib/input/log_file"
//import "./lib/input/log_unix_pipe"

import "./lib/output/json"
import "./lib/match"

import "os"

// Function main() performs a few steps: reading of input, matching and
// priting of JSON output to STDOUT.
func main() {
	//logLines := unixpipe.ReadLog()
	
	//if len(os.Args) == 2 {
	logLines := file.ReadLog(os.Args[1])
	//}
	
	
	patterns := patterns.ReadPatterns("Patterns")
	tokens := tokens.ReadTokens("Tokens")

	matches := match.GetMatches(logLines, patterns, tokens)

	json.PrintJSON(matches)
	return
}

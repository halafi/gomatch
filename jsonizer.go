package main

import "./lib/match"
import "./lib/input"
import "./lib/output"

// Function main() works in a few steps: reading of input, construction
// of prefix tree (trie), matching and priting output to STDOUT.
func main() {
	logLines := input.ReadLog()
	patterns := input.ReadPatterns()
	tokenDefinitions := input.ReadTokens()

	trie, finalFor, stateIsTerminal := match.ConstructPrefixTree(tokenDefinitions, patterns) 
	
	matches := match.GetPatternsLogMatches(logLines, tokenDefinitions, patterns, trie, finalFor, stateIsTerminal)
	
	output.PrintJSON(matches)
	
	return
}

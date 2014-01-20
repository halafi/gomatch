package main

import "github.com/halafi/String-matching-Go/input"
import "github.com/halafi/String-matching-Go/output"
import "github.com/halafi/String-matching-Go/match"

// Function main() works in a few steps: reading of input, construction
// of prefix tree (trie), matching and priting output to STDOUT.
func main() {
	logLines := input.ReadLog()
	patterns := input.ReadPatterns()
	tokenDefinitions := input.ReadTokens()

	trie, finalFor, stateIsTerminal := match.ConstructPrefixTree(tokenDefinitions, patterns) 
	
	matches := match.GetMatches(logLines, tokenDefinitions, patterns, trie, finalFor, stateIsTerminal)
	
	output.PrintJSON(matches)
	return
}

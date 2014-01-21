package main

import "./lib/input"
import "./lib/output"
import "./lib/match"

// Function main() performs a few steps: reading of input, matching and 
// priting of JSON output to STDOUT.
func main() {
	logLines := input.ReadLog()
	patterns := input.ReadPatterns("Patterns")
	tokenDefinitions := input.ReadTokens("Tokens")

	matches := match.GetMatches(logLines, patterns, tokenDefinitions)
	
	output.PrintJSON(matches)
	return
}

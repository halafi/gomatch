package main

import "./lib/input/log_file"
import "./lib/input/log_unixpipe"
import "./lib/input/tokens"
import "./lib/input/patterns"
import "./lib/match"
import "./lib/match/trie"
import "./lib/output/json"
import "os"
import "flag"
import "log"

// Command-line flags.
var input = flag.String("i", "os.Stdin", "Log data input.")
var patternsIn = flag.String("p", "./Patterns", "Pattern definitions input.")
var tokensIn = flag.String("t", "./Tokens", "Token definitions input.")
var output = flag.String("o", "os.Stdout", "JSON data output.")

// Function main() performs a few steps: 
func main() {
	flag.Parse()
	tokens := tokens.ReadTokens(*tokensIn)
	patternReader := patterns.Init(*patternsIn)
	
	patternsArr := make([]string, 0)
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

	if *input != "os.Stdin" && *output =="os.Stdout" { // Read file, write pipe.
		logLines := file.ReadLog(*input)
		
		for n := range logLines {
			json.Get(match.GetMatch(logLines[n], patternsArr, tokens, tree, finalFor))
		}
	} else if *input == "os.Stdin" && *output =="os.Stdout" { // Read pipe, write pipe.
		unixPipeReader := unixpipe.Init()
		
		for {
			logLine, eof := unixpipe.ReadLine(unixPipeReader)
			if eof {
				json.Get(match.GetMatch(logLine, patternsArr, tokens, tree, finalFor))
				return
			} else {
				json.Get(match.GetMatch(logLine, patternsArr, tokens, tree, finalFor))
			}
		}
	} else if *input != "os.Stdin" && *output != "os.Stdout" { // Read file, write file.
		logLines := file.ReadLog(*input)
		
		file, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		
		for n := range logLines {
			_, err := file.WriteString(json.Get(match.GetMatch(logLines[n], patternsArr, tokens, tree, finalFor)))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if *input == "os.Stdin" && *output != "os.Stdout" { // Read pipe, write file.
		unixPipeReader := unixpipe.Init()
		
		file, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		
		for {
			logLine, eof := unixpipe.ReadLine(unixPipeReader)
			if eof {
				_, err := file.WriteString(json.Get(match.GetMatch(logLine, patternsArr, tokens, tree, finalFor)))
				if err != nil {
					log.Fatal(err)
				}
				return
			} else {
				_, err := file.WriteString(json.Get(match.GetMatch(logLine, patternsArr, tokens, tree, finalFor)))
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	return
}

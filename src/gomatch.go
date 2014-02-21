// main.go - program core.
package main

import (
	"flag"
	"log"
	"os"
)

// Command-line flags.
var input = flag.String("i", "/dev/stdin", "Data input.")
var patternsIn = flag.String("p", "Patterns", "Patterns input.")
var tokensIn = flag.String("t", "Tokens", "Tokens input.")
var output = flag.String("o", "/dev/stdout", "Matched data output.")
var noMatchOut = flag.String("u", "no_match.log", "Unmatched data output.")
var outputFormat = flag.String("f", "json", "Matched data format. Supported: json, xml, name, none.")
var inputSocket = flag.String("s", "none", "Reading from Socket.")

// main function starts when the program is executed.
func main() {
	flag.Parse()
	
	tokens := readTokens(*tokensIn)
	
	trie, finalFor, state, i := createNewTrie()
	
	patternsArr := readPatterns(*patternsIn)
	for p := range patternsArr {
		finalFor, state, i = appendPattern(tokens, patternsArr[p], trie, finalFor, state, i)
	}

	outputFile := createFile(*output) // file for matched output
	unmatchedOutputFile := createFile(*noMatchOut) // file for nomatch
	
	// for dynamic pattern insert, stores patterns file info for later
	patternReader := openFile(*patternsIn)
	patternsFileInfo, err := os.Stat(*patternsIn)
	if err != nil {
		log.Fatal(err)
	}
	lastModified := patternsFileInfo.ModTime()

	// reading of input lines from either socket or file, matching them
	// and writing them to output until EOF
	if *inputSocket != "none" {
		connection := openSocket(*inputSocket)
		// do until eof
		for {
			// check for a new pattern in patterns file at first line
			patternsFileInfo, err = os.Stat(*patternsIn)
			if err != nil {
				log.Fatal(err)
			}
			if lastModified != patternsFileInfo.ModTime() {
				patternReader = openFile(*patternsIn)
				line, eof := readLine(patternReader)
				if !eof {
					oldLen := len(patternsArr)
					patternsArr = addPattern(line, patternsArr)
					if len(patternsArr) > oldLen {
						finalFor, state, i = appendPattern(tokens, patternsArr[len(patternsArr)-1], trie, finalFor, state, i)
						log.Printf("New event: \"%s\".", patternsArr[len(patternsArr)-1].Name)
					}
					lastModified = patternsFileInfo.ModTime()
				}
			}
			// read everything from socket
			lines, eof := readFully(connection)
			for i := range lines {
				match := getMatch(lines[i], patternsArr, tokens, trie, finalFor)
				if match.Type != "" {
					writeFile(outputFile, convertMatch(match) + "\r\n")
				} else {
					writeFile(unmatchedOutputFile, lines[i] + "\r\n")
				}
			}
			if eof {
				break
			}
		}
		connection.Close()
	} else {
		inputReader := openFile(*input)
		// do until eof
		for {
			// check for a new pattern in patterns file at first line
			patternsFileInfo, err = os.Stat(*patternsIn)
			if err != nil {
				log.Fatal(err)
			}
			if lastModified != patternsFileInfo.ModTime() {
				patternReader = openFile(*patternsIn)
				line, eof := readLine(patternReader)
				if !eof {
					oldLen := len(patternsArr)
					patternsArr = addPattern(line, patternsArr)
					if len(patternsArr) > oldLen {
						finalFor, state, i = appendPattern(tokens, patternsArr[len(patternsArr)-1], trie, finalFor, state, i)
						log.Printf("New event: \"%s\".", patternsArr[len(patternsArr)-1].Name)
					}
					lastModified = patternsFileInfo.ModTime()
				}
			}
			// read log line
			logLine, eof := readLine(inputReader)
			if eof {
				match := getMatch(logLine, patternsArr, tokens, trie, finalFor)
				if match.Type != "" {
					writeFile(outputFile, convertMatch(match) + "\r\n")
				} else {
					writeFile(unmatchedOutputFile, logLine + "\r\n")
				}
				break
			} else {
				match := getMatch(logLine, patternsArr, tokens, trie, finalFor)
				if match.Type != "" {
					writeFile(outputFile, convertMatch(match) + "\r\n")
				} else {
					writeFile(unmatchedOutputFile, logLine + "\r\n")
				}
			}
		}
	}
	closeFile(outputFile)
	return
}

// convertMatch returns the desired output for a given match.
func convertMatch(match Match) string {
	if *outputFormat == "JSON" || *outputFormat == "json" {
		return getJSON(match)
	}
	if *outputFormat == "XML" || *outputFormat == "xml" {
		return getXML(match)
	}
	if *outputFormat == "name" {
		return match.Type
	}
	if *outputFormat == "none" {
		return ""
	}
	log.Fatal("unknown output format: \"", *outputFormat+"\"")
	return ""
}

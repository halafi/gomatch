// main.go - program core.
package main

import (
	"flag"
	"log"
	"os"
)

// Command-line flags.
var input = flag.String("i", "/dev/stdin", "Data input stream.")
var output = flag.String("o", "/dev/stdout", "Data output stream.")
var patternsIn = flag.String("p", "Patterns", "Pattern definitions input.")
var tokensIn = flag.String("t", "Tokens", "Token definitions input.")
var outputFormat = flag.String("f", "json", "Output data format. Supported: json, xml, name, none.")
var inputSocket = flag.String("s", "none", "Reading from Socket.")

// main starts when the program is executed.
func main() {
	// parsing of command-line flags.
	flag.Parse()
	// token reading
	tokenReader := openFile(*tokensIn)
	tokens := make(map[string]string)
	for {
		line, eof := readLine(tokenReader)
		if eof {
			break
		}
		token := checkToken(line)
		if token != "" {
			addToken(token, tokens)
		}
	}
	// initial pattern reading and trie construction
	tree, finalFor, state, i := createNewTrie()
	patternReader := openFile(*patternsIn)
	patternsArr := make([]string, 0)
	for {
		line, eof := readLine(patternReader)
		if eof {
			break
		}
		pattern := checkPattern(line)
		if pattern != "" {
			patternsArr = stringArraySizeUp(patternsArr, 1)
			patternsArr[len(patternsArr)-1] = pattern
			tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i)
		}
	}
	// store patterns file info for later use
	patternsFileInfo, err := os.Stat(*patternsIn)
	if err != nil {
		log.Fatal(err)
	}
	lastModified := patternsFileInfo.ModTime()
	// open file for output
	outputFile := createFile(*output)

	// reading of input lines from either socket or file, matching them
	// and writing them to output until EOF
	if *inputSocket != "none" {
		l := startServer(*inputSocket)
		con := openSocket(*inputSocket)
		write(con, "2013-02-26T12:24:05.425+00:00 WARN org.ssh.ServerImpl - Failed password for xtovarn from 147.251.49.42 port 46177 #1#")
		// run server
		fd, err := l.Accept()
		if err != nil {
			log.Println("accept error", err)
		}
		go echoServer(fd)
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
					pattern := checkPattern(line)
					if pattern != "" && !contains(patternsArr, pattern) {
						log.Printf("New event: \"%s\".", pattern)
						patternsArr = stringArraySizeUp(patternsArr, 1)
						patternsArr[len(patternsArr)-1] = pattern
						tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i)
					}
					lastModified = patternsFileInfo.ModTime()
				}
			}
			// read everything from socket
			lines, eof := readFully(con)
			for i := range lines {
				writeFile(outputFile, convertMatch(getMatch(lines[i], patternsArr, tokens, tree, finalFor), *outputFormat))
			}
			if eof {
				break
			}
		}
		con.Close()
		closeServer(l)
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
					pattern := checkPattern(line)
					if pattern != "" && !contains(patternsArr, pattern) {
						log.Printf("New event: \"%s\".", pattern)
						patternsArr = stringArraySizeUp(patternsArr, 1)
						patternsArr[len(patternsArr)-1] = pattern
						tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i)
					}
					lastModified = patternsFileInfo.ModTime()
				}
			}
			// read log line
			logLine, eof := readLine(inputReader)
			if eof {
				writeFile(outputFile, convertMatch(getMatch(logLine, patternsArr, tokens, tree, finalFor), *outputFormat))
				break
			} else {
				writeFile(outputFile, convertMatch(getMatch(logLine, patternsArr, tokens, tree, finalFor), *outputFormat))
			}
		}
	}
	closeFile(outputFile)
	return
}

// convertMatch returns the desired output for a given match.
func convertMatch(match Match, output string) string {
	if output == "JSON" || output == "json" {
		return getJSON(match)
	}
	if output == "XML" || output == "xml" {
		return getXML(match)
	}
	if output == "name" {
		if match.Type == "" {
			return ""
		}
		return match.Type + "\r\n"
	}
	if output == "none" {
		return ""
	}
	log.Fatal("unknown output format: \"", output+"\"")
	return ""
}

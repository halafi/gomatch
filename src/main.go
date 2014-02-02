package main

import "flag"
import "log"
import "os"

// Command-line flags.
var input = flag.String("i", "/dev/stdin", "Data input stream.")
var inputSocket = flag.String("s", "/tmp/echo.sock", "Data input Unix domain socket (none or filePath).")
var output = flag.String("o", "/dev/stdout", "Data output stream.")
var outputFormat = flag.String("f", "json", "Output data format, supported: json, xml.")
var patternsIn = flag.String("p", "./Patterns", "Pattern definitions input.")
var tokensIn = flag.String("t", "./Tokens", "Token definitions input.")

// Function main() performs a few steps: 
func main() {
	flag.Parse()
	// Token reading
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
	tree, finalFor, state, i := createNewTrie()
	// Initial pattern reading
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
			patternsArr[len(patternsArr)-1] = pattern // add pattern to array of all patterns
			tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i) // add pattern to trie
		}
	}
	patternsFileInfo, err := os.Stat(*patternsIn)
	if err != nil {
		log.Fatal(err)
	}
	lastModified := patternsFileInfo.ModTime()
	
	outputFile := createFile(*output)
	// Reading of input lines, matching them and writing them to output
	if *inputSocket != "none" { // Reading from socket
		l := startServer(*inputSocket) // start socket server
		
		con := openSocket(*inputSocket) // client connects
		write(con, "2013-02-26T12:24:05.425+00:00 WARN org.ssh.ServerImpl - Failed password for xtovarn from 147.251.49.42 port 46177 #1#") // client sends data
		
		// run server
		fd, err := l.Accept()
		if err != nil {
			log.Println("accept error", err)
		}
		go echoServer(fd)
		
		// read server fully until eof
		for {
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
	} else { // Reading from file
		inputReader := openFile(*input)
		for {
			// If last mod time for patterns file is different, then read
			// the first line of patterns file and check for a new pattern
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
						patternsArr[len(patternsArr)-1] = pattern // add pattern to array of all patterns
						tree, finalFor, state, i = appendPattern(tokens, pattern, tree, finalFor, state, i) // add pattern to trie
					}
					lastModified = patternsFileInfo.ModTime()
				}
				
			}
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

// Calls the desired get method and returns its output.
func convertMatch(match Match, output string) string {
	if output=="JSON" || output=="json" {
		return getJSON(match)
	}
	if output=="XML" || output=="xml" {
		return getXML(match)
	}
	log.Fatal("unknown output format: \"", output +"\"")
	return ""
}

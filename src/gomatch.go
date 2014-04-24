package main

import (
	"flag"
	"log"
	"os"
	"time"
)

var (
	// Command-line flags.
	inputFilePath         = flag.String("i", "/dev/stdin", "Data input.")
	inputSocketFilePath   = flag.String("s", "none", "Reading from Socket.")
	ampqConfigFilePath    = flag.String("a", "none", "Filepath for AMQP config file.")
	patternsFilePath      = flag.String("p", "Patterns", "Patterns input.")
	tokensFilePath        = flag.String("t", "Tokens", "Tokens input.")
	outputFilePath        = flag.String("o", "/dev/stdout", "Matched data output.")
	noMatchOutputFilePath = flag.String("n", "no_match.log", "Unmatched data output.")
	// Shared variables between all goroutines.
	trie          map[int]map[Token]int
	finalFor      []int
	state         int
	patternNumber int
	patterns      []Pattern
	regexMap      map[string]Regex
)

// main starts when the program is executed.
// Performs parsing of flags, loads up patterns and regexes, constructs
// a trie and runs a separate goroutine for watching file with patterns.
// Has three modes: AMQP, Socket, File/Pipeline.
// Log lines are matched and send to some output in JSON.
func main() {
	flag.Parse()

	if *ampqConfigFilePath != "none" && *inputSocketFilePath != "none" {
		log.Fatal("cannot use both socket and amqp at the same time")
	}

	trie, finalFor, state, patternNumber = initTrie()

	regexMap, patterns = readPatterns(*patternsFilePath, *tokensFilePath)
	for p := range patterns {
		finalFor, state, patternNumber = appendPattern(patterns[p], trie, finalFor, state, patternNumber, regexMap)
	}

	go watchPatterns()

	if *ampqConfigFilePath != "none" { // amqp
		// init configuration parameters
		parseAmqpConfigFile(*ampqConfigFilePath)
		// set up connections and channels, ensure that they are closed
		cSend := openConnection(amqpMatchedSendUri)
		chSend := openChannel(cSend)
		defer cSend.Close()
		defer chSend.Close()

		cReceive := openConnection(amqpReceiveUri)
		chReceive := openChannel(cReceive)
		defer cReceive.Close()
		defer chReceive.Close()

		// declare queues
		qReceive := declareQueue(amqpReceiveQueueName, chReceive)
		qSend := declareQueue(amqpMatchedSendQueueName, chSend)

		// bind the receive exchange with the receive queue
		bindReceiveQueue(chSend, qReceive)

		// start consuimng until terminated
		msgs, err := chReceive.Consume(qReceive.Name, "", true, false, false, false, nil)
		if err != nil {
			log.Fatal(err)
		}
		switch amqpReceiveFormat {
		case "plain", "PLAIN": // incoming logs
			noMatchOutputFile := createFile(*noMatchOutputFilePath)
			defer noMatchOutputFile.Close()
			for delivery := range msgs {
				match := getMatch(string(delivery.Body), patterns, trie, finalFor, regexMap)
				if match.Type != "" {
					send([]byte(marshalJson(match)), match.Type, chSend, qSend) // routing key = pattern_name
				} else {
					writeFile(noMatchOutputFile, string(delivery.Body)+"\r\n")
				}
			}
		case "json", "JSON": // incoming json
			for delivery := range msgs {
				m := unmarshalJson(delivery.Body)
				if attExists("@gomatch", m) { // att @gomatch is present
					if str, ok := m["@gomatch"].(string); ok {
						match := getMatch(str, patterns, trie, finalFor, regexMap)
						if match.Type != "" {
							m["@type"] = match.Type
							m["@p"] = match.Body
							delete(m, "@gomatch")
							entityStr, _ := m["@entity"].(string)
							send([]byte(marshalJson(m)), entityStr+"."+match.Type, chSend, qSend) // routing key = @entity.pattern_name
						}
					} else {
						log.Println("@gomatch is not a string (skipping)")
					}
				} else { // we return the former json msg
					entityStr, _ := m["@entity"].(string)
					send(delivery.Body, entityStr, chSend, qSend) // routing key = @entity
				}
			}
		default:
			log.Fatal("Unknown RabbitMQ input format, use either plain or json.")
		}

	} else if *inputSocketFilePath != "none" { // socket
		conn := openSocket(*inputSocketFilePath)
		outputFile := createFile(*outputFilePath)
		noMatchOutputFile := createFile(*noMatchOutputFilePath)
		defer conn.Close()
		defer outputFile.Close()
		defer noMatchOutputFile.Close()

		for {
			lines, eof := readFully(conn)
			for i := range lines {
				match := getMatch(lines[i], patterns, trie, finalFor, regexMap)
				if match.Type != "" {
					writeFile(outputFile, marshalMatch(match)+"\r\n")
				} else {
					writeFile(noMatchOutputFile, lines[i]+"\r\n")
				}
			}
			if eof {
				break
			}
		}

	} else { // file, pipeline
		outputFile := createFile(*outputFilePath)
		noMatchOutputFile := createFile(*noMatchOutputFilePath)
		defer outputFile.Close()
		defer noMatchOutputFile.Close()

		inputReader := openFile(*inputFilePath)

		for {
			line, eof := readLine(inputReader)
			logLine := string(line)
			match := getMatch(logLine, patterns, trie, finalFor, regexMap)
			if match.Type != "" {
				writeFile(outputFile, marshalMatch(match)+"\r\n")
			} else {
				writeFile(noMatchOutputFile, logLine+"\r\n")
			}
			if eof {
				break
			}
		}
	}
	return
}

// watchPatterns performs re-reading of the first line in Patterns file
// (if it was recently modified).
// Then tries to add that line as a new pattern to trie.
func watchPatterns() {
	patternsFileInfo, err := os.Stat(*patternsFilePath)
	if err != nil {
		log.Fatal("watchPatterns: ", err)
	}
	patternsLastModTime := patternsFileInfo.ModTime()
	for {
		time.Sleep(1 * time.Second)

		patternsFileInfo, err := os.Stat(*patternsFilePath)
		if err != nil {
			log.Println("watchPatterns: ", err)
			break
		}

		if patternsLastModTime != patternsFileInfo.ModTime() {
			patternReader := openFile(*patternsFilePath)
			line, eof := readLine(patternReader)
			if !eof {
				oldLen := len(patterns)
				patterns = addPattern(string(line), patterns, regexMap)

				if len(patterns) > oldLen {
					finalFor, state, patternNumber = appendPattern(patterns[len(patterns)-1], trie, finalFor, state, patternNumber, regexMap)
					log.Println("new event: ", patterns[len(patterns)-1].Name)
				}
				patternsLastModTime = patternsFileInfo.ModTime()
			}
		}
	}
}

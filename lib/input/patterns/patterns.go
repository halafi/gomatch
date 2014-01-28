// Package patterns provides funcionality for reading of patterns.
package patterns

import "log"
import "strings"
import "io"
import "bufio"
import "os"

// Init does the initialization of buffered io.Reader for reading from
// given file path.
func Init(filePath string) *bufio.Reader {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)
	return reader
}

// Reads a single file line using a given reader.
func ReadPattern(reader *bufio.Reader) (pattern string, eof bool) {
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return "", true
			} else {
				log.Fatal(err)
			}
		}
		return checkPattern(string(line)), false
	}
}

// Function checkPattern validates given pattern, if it passes the given
// pattern is returned, otherwise error is logged.
func checkPattern(pattern string) string {
	patternNameSplit := strings.Split(pattern, "##") //separate pattern name from its definition
	if len(patternNameSplit[0]) == 0 {
		log.Fatal("pattern error \"", pattern, "\": name cannot be empty.")
	}
	if len(patternNameSplit[1]) == 0 {
		log.Fatal("pattern error \"", pattern, "\": cannot be empty.")
	}
	if len(patternNameSplit) != 2 {
		log.Fatal("pattern error: \"", pattern+"\"")
	}
	return pattern
}

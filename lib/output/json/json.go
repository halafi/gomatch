// Package json provides output funcionality for JSON output to STDOUT.
package json

import "../../match"
import "encoding/json"
import "log"
import "os"

// Returns JSON for given match as string.
func Get(matchPerLine match.Match) string {
	if matchPerLine.Type != "" {
		b, err := json.Marshal(matchPerLine)
		if err != nil {
			log.Fatal(err)
		}
		return string(b) + "\r\n"
	}
	return ""
}

// Creates file for output data at 'filePath'.
func CreateOutputFile(filePath string) *os.File {
	f, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// Writes JSON of 'matchPerLine' to 'file'.
func WriteOutputFile(file *os.File, matchPerLine match.Match) {
	_, err := file.WriteString(Get(matchPerLine))
	if err != nil {
		log.Fatal(err)
	}
}

// Closes the given 'file'.
func CloseOutputFile(file *os.File) {
	file.Close()
}

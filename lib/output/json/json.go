// Package json provides output funcionality for JSON output to STDOUT.
package json

import "encoding/json"
import "fmt"
import "log"
import "../../match"

// PrintJSON takes a single Match input and prints it to STDOUT in 
// JSON.
func PrintJSON(matchPerLine match.Match) {
	if matchPerLine.Type != "" {
		b, err := json.Marshal(matchPerLine)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(b) + "\r\n")
	}
}

// Package json provides output funcionality for JSON output to STDOUT.
package json

import "encoding/json"
import "log"
import "../../match"

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

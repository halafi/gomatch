package main

import "encoding/json"
import "log"

// Returns JSON for given match as string.
func getJSON(matchPerLine Match) string {
	if matchPerLine.Type != "" {
		b, err := json.Marshal(matchPerLine)
		if err != nil {
			log.Fatal(err)
		}
		return string(b) + "\r\n"
	}
	return ""
}

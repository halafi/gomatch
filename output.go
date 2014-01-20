// Package output provides output funcionality - writing JSON text data
// to STDOUT.
package output

import "encoding/json"
import "fmt"
import "log"
import "match"

// Determines JSON output indent (formatting), you can use anything like
// three spaces(default) or "\t".
const indent = "   " 

// PrintJSON prints to STDOUT a formatted JSON text data.
func PrintJSON(matchPerLine []match.Match) {
	output := getJSON(matchPerLine)
	if output != "[\r\n]" {
		fmt.Printf("%s\r\n", getJSON(matchPerLine))
	}
}

// For a set of matches this function returns a string containing 
// formatted JSON data.
func getJSON(matchPerLine []match.Match) string {
	output := "["
	first := true
	for n := range matchPerLine {
		if matchPerLine[n].Type	!= "" {
			if !first {
				output = output + ","
			} else {
				first = false
			}
			b, err := json.MarshalIndent(matchPerLine[n], indent+indent, indent)
			if err != nil {
				log.Fatal(err)
			}
			output = output + "\r\n"+indent+"{\r\n"+indent+indent+"\"Event\": " + string(b)+"\r\n"+indent+"}"
		}
	}
	return output + "\r\n]"
}

// match-xml.go provides funcionality for conversion of struct Match to
// XML string.
package main

// getXML converts given Match to XML.
func getXML(match Match) string {
	xml := "<Event type=\"" + match.Type + "\">"
	if len(match.Body) != 0 {
		xml = xml + "<Body>"
		for i := 0; i < len(match.Body)-1; i = i + 2 {
			xml = xml +
				"<" + match.Body[i] + ">" +
				match.Body[i+1] +
				"</" + match.Body[i] + ">"
		}
		xml = xml + "</Body></Event>"
	} else {
		xml = xml + "</Event>"
	}
	return xml
}

package main

// getXML converts given match to XML string.
func getXML(match Match) string {
	xml := "<Event type=\"" + match.Type + "\">"
	if len(match.Body) != 0 {
		xml = xml + "<Body>"
		for i := 0; i < len(match.Body)-1; i = i + 2 {
			xml = xml + "<" + match.Body[i] + ">" + match.Body[i+1] + "</" + match.Body[i] + ">"
		}
		xml = xml + "</Body></Event>"
	} else {
		xml = xml + "</Event>"
	}
	return xml
}

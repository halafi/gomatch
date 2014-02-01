package main

import "log"

// Converts given Match to JSON.
func getJSON(match Match) string {
	json := ""
	if match.Type != "" {
		json = "{\"Event\":{\"type\":\""+match.Type+"\""
		if len(match.Body) != 0 {
			json = json +",\"body\":[{"
			for i := 0; i < len(match.Body)-1; i=i+2 {
				json = json + "\""+match.Body[i] +"\":\"" + match.Body[i+1]+"\""
				if i != len(match.Body)-2 {
					json = json +","
				}
			}
			json = json +"}]}}"
		} else {
			json = json + "}}"
		}
		return json+"\r\n"
	}
	return json
}

// Converts given Match to XML.
func getXML(match Match) string {
	xml := ""
	if match.Type != "" {
		xml = "<Event type=\""+match.Type+"\">"
		if len(match.Body) != 0 {
			xml = xml +"<Body>"
			for i := 0; i < len(match.Body)-1; i=i+2 {
				xml = xml + "<"+match.Body[i] +">" + match.Body[i+1] + "</"+match.Body[i] +">"
			}
			xml = xml +"</Body></Event>"
		} else {
			xml = xml + "</Event>"
		}
		return xml+"\r\n"
	}
	return xml
}

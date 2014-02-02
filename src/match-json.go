package main

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

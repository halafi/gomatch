package main

// getJSON converts given match to JSON string.
func getJSON(match Match) string {
	json := "{\"Event\":{\"type\":\"" + match.Type + "\""
	if len(match.Body) != 0 {
		json = json + ",\"body\":[{"
		for i := 0; i < len(match.Body)-1; i = i + 2 {
			if i != len(match.Body)-2 {
				json = json + "\"" + match.Body[i] + "\":\"" + match.Body[i+1] + "\", "
			} else {
				json = json + "\"" + match.Body[i] + "\":\"" + match.Body[i+1] + "\""
			}
		}
		json = json + "}]}}"
	} else {
		json = json + "}}"
	}
	return json
}

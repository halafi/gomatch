package main

import (
	"encoding/json"
	"log"
)

// marshalJson converts any object into a JSON string.
func marshalJson(object interface{}) []byte {
	b, err := json.Marshal(object)
	if err != nil {
		log.Println(err)
		return make([]byte, 0)
	}
	return b
}

// marshalMatch converts a given match into a JSON string.
// Performs faster than function marshalJson.
func marshalMatch(match Match) string {
	json := "{\"Event\":{\"Type\":\"" + match.Type + "\""
	if len(match.Body) != 0 {
		json = json + ",\"Body\":[{"
		for k, v := range match.Body {
			json = json + "\"" + k + "\":\"" + v + "\","
		}
		json = json[:len(json)-1] + "}]}}" // remove extra comma
	} else {
		json = json + "}}"
	}
	return json
}

// unmarshalJson takes a JSON object and decodes it into a map.
func unmarshalJson(object []byte) map[string]interface{} {
	var msg interface{}
	err := json.Unmarshal(object, &msg)
	if err != nil {
		log.Fatal(err)
	}
	return msg.(map[string]interface{})
}

// attExists checks whether an attribute exists in a given map.
func attExists(att string, m map[string]interface{}) bool {
	if _, exists := m[att]; !exists || m[att] == "" {
		return false
	}
	return true
}

package utils

import "encoding/json"

func StringifyJson(in interface{}) string {
	b, err := json.Marshal(in)
	if err != nil {
		return "[StringifyJson] " + err.Error()
	}
	return string(b)
}

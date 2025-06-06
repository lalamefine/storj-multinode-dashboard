package main

import (
	"encoding/json"
)

// MarshalJSON convertit une structure Go en chaîne JSON
func MarshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(data)
}

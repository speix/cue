package helpers

import "encoding/json"

// InBetween checks if a given number is between two others
func InBetween(number, min, max int) bool {

	if (number >= min) && (number <= max) {
		return true
	}

	return false
}

// IsJSON checks if a given string is a proper JSON object
func IsJSON(str json.RawMessage) bool {
	var js json.RawMessage

	return json.Unmarshal([]byte(string(str)), &js) == nil
}

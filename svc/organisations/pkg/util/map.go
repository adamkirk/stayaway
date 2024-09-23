package util

import (
	"strings"
)

// checkPathExists takes a dot-notated path and checks if the path exists in the map
func KeyExistsInMap(data map[string]any, path string) (bool) {
	// Split the dot notation into keys
	keys := strings.Split(path, ".")

	// Recursively check if the path exists
	return checkKeysExist(data, keys)
}

// Helper function to recursively check if the keys exist
func checkKeysExist(data map[string]any, keys []string) (bool) {
	// If there are no keys left, return false (invalid path)
	if len(keys) == 0 {
		return false
	}

	// Get the current key
	key := keys[0]

	// Check if the key exists in the map
	value, exists := data[key]
	if !exists {
		return false
	}

	// If it's the last key, the path exists
	if len(keys) == 1 {
		return true
	}

	// If there are more keys, the current value must be a map for the path to continue
	if nestedMap, ok := value.(map[string]any); ok {
		return checkKeysExist(nestedMap, keys[1:])
	}

	// If the value is not a map but more keys remain, the path doesn't exist
	return false
}
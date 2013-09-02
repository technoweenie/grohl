package grohl

import (
	"fmt"
	"strings"
)

// Builds a log message as a single line from log data.
func BuildLine(data map[string]interface{}) string {
	index := 0
	pieces := make([]string, len(data))
	for key, value := range data {
		pieces[index] = fmt.Sprintf("%s=%s", key, value)
		index = index + 1
	}

	return strings.Join(pieces, space)
}

func buildLine(context map[string]interface{}, data map[string]interface{}) string {
	return BuildLine(dupeMaps(context, data))
}

func dupeMaps(maps ...map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for _, orig := range maps {
		for key, value := range orig {
			merged[key] = value
		}
	}
	return merged
}

const space = " "

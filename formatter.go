package grohl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Builds a log message as a single line from log data.
func BuildLine(data map[string]interface{}) string {
	index := 0
	pieces := make([]string, len(data))
	for key, value := range data {
		pieces[index] = fmt.Sprintf("%s=%s", key, formatValue(value))
		index = index + 1
	}

	return strings.Join(pieces, space)
}

func formatValue(value interface{}) string {
	if value == nil {
		return "nil"
	}

	k := reflect.TypeOf(value).Kind()
	formatter := formatters[k]
	if formatter == nil {
		return "unknown"
	}

	return formatter(value)
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

var durationFormat = []byte("f")[0]

var formatters = map[reflect.Kind]func(value interface{}) string{
	reflect.String: func(value interface{}) string {
		return value.(string)
	},

	reflect.Bool: func(value interface{}) string {
		return strconv.FormatBool(value.(bool))
	},

	reflect.Int: func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int)), 10)
	},

	reflect.Int8: func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int8)), 10)
	},

	reflect.Int16: func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int16)), 10)
	},

	reflect.Int32: func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int32)), 10)
	},

	reflect.Int64: func(value interface{}) string {
		return strconv.FormatInt(value.(int64), 10)
	},

	reflect.Float32: func(value interface{}) string {
		return strconv.FormatFloat(float64(value.(float32)), durationFormat, 3, 32)
	},

	reflect.Float64: func(value interface{}) string {
		return strconv.FormatFloat(value.(float64), durationFormat, 3, 64)
	},

	reflect.Uint: func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint)), 10)
	},

	reflect.Uint8: func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	},

	reflect.Uint16: func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	},

	reflect.Uint32: func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	},

	reflect.Uint64: func(value interface{}) string {
		return strconv.FormatUint(value.(uint64), 10)
	},
}

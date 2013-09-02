package grohl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
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

	t := reflect.TypeOf(value)
	formatter := formatters[t.Kind().String()]
	if formatter == nil {
		formatter = formatters[t.String()]
	}

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

const (
	space      = " "
	timeLayout = "2006-01-02T15:04:05-0700"
)

var durationFormat = []byte("f")[0]

var formatters = map[string]func(value interface{}) string{
	"string": func(value interface{}) string {
		return value.(string)
	},

	"bool": func(value interface{}) string {
		return strconv.FormatBool(value.(bool))
	},

	"int": func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int)), 10)
	},

	"int8": func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int8)), 10)
	},

	"int16": func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int16)), 10)
	},

	"int32": func(value interface{}) string {
		return strconv.FormatInt(int64(value.(int32)), 10)
	},

	"int64": func(value interface{}) string {
		return strconv.FormatInt(value.(int64), 10)
	},

	"float32": func(value interface{}) string {
		return strconv.FormatFloat(float64(value.(float32)), durationFormat, 3, 32)
	},

	"float64": func(value interface{}) string {
		return strconv.FormatFloat(value.(float64), durationFormat, 3, 64)
	},

	"uint": func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint)), 10)
	},

	"uint8": func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	},

	"uint16": func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	},

	"uint32": func(value interface{}) string {
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	},

	"uint64": func(value interface{}) string {
		return strconv.FormatUint(value.(uint64), 10)
	},

	"time.Time": func(value interface{}) string {
		return value.(time.Time).Format(timeLayout)
	},
}

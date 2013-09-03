package grohl

import (
	"fmt"
	"testing"
	"time"
)

var exampleTime = time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
var exampleError = fmt.Errorf("error message")

var examples = map[string]LogData{
	"fn=string test=hi": LogData{
		"fn": "string", "test": "hi",
	},
	`fn=stringspace test="a b"`: LogData{
		"fn": "stringspace", "test": "a b",
	},
	`fn=stringslasher test="slasher \\\\"`: LogData{
		"fn": "stringslasher", "test": `slasher \\`,
	},
	`fn=stringeqspace test="x=4, y=10"`: LogData{
		"fn": "stringeqspace", "test": "x=4, y=10",
	},
	`fn=stringeq test="x=4,y=10"`: LogData{
		"fn": "stringeq", "test": "x=4,y=10",
	},
	`fn=stringspace test="hello world"`: LogData{
		"fn": "stringspace", "test": "hello world",
	},
	`fn=stringbothquotes test="echo 'hello' \"world\""`: LogData{
		"fn": "stringbothquotes", "test": `echo 'hello' "world"`,
	},
	`fn=stringsinglequotes test="a 'a'"`: LogData{
		"fn": "stringsinglequotes", "test": `a 'a'`,
	},
	`fn=stringdoublequotes test='echo "hello"'`: LogData{
		"fn": "stringdoublequotes", "test": `echo "hello"`,
	},
	`fn=stringbothquotesnospace test='a"`: LogData{
		"fn": "stringbothquotesnospace", "test": `'a"`,
	},
	"fn=emptystring test=nil": LogData{
		"fn": "emptystring", "test": "",
	},
	"fn=int test=1": LogData{
		"fn": "int", "test": int(1),
	},
	"fn=int8 test=1": LogData{
		"fn": "int8", "test": int8(1),
	},
	"fn=int16 test=1": LogData{
		"fn": "int16", "test": int16(1),
	},
	"fn=int32 test=1": LogData{
		"fn": "int32", "test": int32(1),
	},
	"fn=int64 test=1": LogData{
		"fn": "int64", "test": int64(1),
	},
	"fn=uint test=1": LogData{
		"fn": "uint", "test": uint(1),
	},
	"fn=uint8 test=1": LogData{
		"fn": "uint8", "test": uint8(1),
	},
	"fn=uint16 test=1": LogData{
		"fn": "uint16", "test": uint16(1),
	},
	"fn=uint32 test=1": LogData{
		"fn": "uint32", "test": uint32(1),
	},
	"fn=uint64 test=1": LogData{
		"fn": "uint64", "test": uint64(1),
	},
	"fn=float test=1.000": LogData{
		"fn": "float", "test": float32(1.0),
	},
	"fn=bool test=true": LogData{
		"fn": "bool", "test": true,
	},
	"fn=nil test=nil": LogData{
		"fn": "nil", "test": nil,
	},
	"fn=time test=2000-01-02T03:04:05+0000": LogData{
		"fn": "time", "test": exampleTime,
	},
	`fn=error test="error message"`: LogData{
		"fn": "error", "test": exampleError,
	},
}

func TestFormat(t *testing.T) {
	for expected, data := range examples {
		if actual := BuildLine(data, false); expected != actual {
			t.Errorf("Expected %s\nGot: %s", expected, actual)
		}
	}
}

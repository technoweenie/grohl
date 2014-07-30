package grohl

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var exampleTime = time.Date(2000, 1, 2, 3, 4, 5, 6, time.UTC)
var exampleError = fmt.Errorf("error message")

type ExampleStruct struct {
	Value interface{}
}

var actuals = []Data{
	Data{"fn": "string", "test": "hi"},
	Data{"fn": "stringspace", "test": "a b"},
	Data{"fn": "stringslasher", "test": `slasher \\`},
	Data{"fn": "stringeqspace", "test": "x=4, y=10"},
	Data{"fn": "stringeq", "test": "x=4,y=10"},
	Data{"fn": "stringspace", "test": "hello world"},
	Data{"fn": "stringbothquotes", "test": `echo 'hello' "world"`},
	Data{"fn": "stringsinglequotes", "test": `a 'a'`},
	Data{"fn": "stringdoublequotes", "test": `echo "hello"`},
	Data{"fn": "stringbothquotesnospace", "test": `'a"`},
	Data{"fn": "emptystring", "test": ""},
	Data{"fn": "int", "test": int(1)},
	Data{"fn": "int8", "test": int8(1)},
	Data{"fn": "int16", "test": int16(1)},
	Data{"fn": "int32", "test": int32(1)},
	Data{"fn": "int64", "test": int64(1)},
	Data{"fn": "uint", "test": uint(1)},
	Data{"fn": "uint8", "test": uint8(1)},
	Data{"fn": "uint16", "test": uint16(1)},
	Data{"fn": "uint32", "test": uint32(1)},
	Data{"fn": "uint64", "test": uint64(1)},
	Data{"fn": "float", "test": float32(1.0)},
	Data{"fn": "bool", "test": true},
	Data{"fn": "nil", "test": nil},
	Data{"fn": "time", "test": exampleTime},
	Data{"fn": "error", "test": exampleError},
	Data{"fn": "slice", "test": []byte{86, 87, 88}},
	Data{"fn": "struct", "test": ExampleStruct{Value: "testing123"}},
}

var expectations = [][]string{
	[]string{"fn=string", "test=hi"},
	[]string{"fn=stringspace", `test="a b"`},
	[]string{`fn=stringslasher`, `test="slasher \\\\"`},
	[]string{`fn=stringeqspace`, `test="x=4, y=10"`},
	[]string{`fn=stringeq`, `test="x=4,y=10"`},
	[]string{`fn=stringspace`, `test="hello world"`},
	[]string{`fn=stringbothquotes`, `test="echo 'hello' \"world\""`},
	[]string{`fn=stringsinglequotes`, `test="a 'a'"`},
	[]string{`fn=stringdoublequotes`, `test='echo "hello"'`},
	[]string{`fn=stringbothquotesnospace`, `test='a"`},
	[]string{"fn=emptystring", "test=nil"},
	[]string{"fn=int", "test=1"},
	[]string{"fn=int8", "test=1"},
	[]string{"fn=int16", "test=1"},
	[]string{"fn=int32", "test=1"},
	[]string{"fn=int64", "test=1"},
	[]string{"fn=uint", "test=1"},
	[]string{"fn=uint8", "test=1"},
	[]string{"fn=uint16", "test=1"},
	[]string{"fn=uint32", "test=1"},
	[]string{"fn=uint64", "test=1"},
	[]string{"fn=float", "test=1.000"},
	[]string{"fn=bool", "test=true"},
	[]string{"fn=nil", "test=nil"},
	[]string{"fn=time", "test=2000-01-02T03:04:05+0000"},
	[]string{`fn=error`, `test="error message"`},
	[]string{`fn=slice`, `test="[86 87 88]"`},
	[]string{`fn=struct`, `test={Value:testing123}`},
}

func TestFormat(t *testing.T) {
	for i, actual := range actuals {
		AssertData(t, actual, expectations[i])
	}
}

func TestFormatWithTime(t *testing.T) {
	data := Data{"fn": "time", "test": 1}
	actual := BuildLog(data, true)
	if !strings.HasPrefix(actual, "now=") {
		t.Errorf("Invalid prefix: %s", actual)
	}
	if !strings.HasSuffix(actual, " fn=time test=1") {
		t.Errorf("Invalid suffix: %s", actual)
	}
}

func AssertLog(t *testing.T, ctx *Context, expected []string) {
	AssertData(t, ctx.Merge(nil), expected)
}

func AssertData(t *testing.T, data Data, expected []string) {
	pairs, line := buildLogMap(data)
	for _, pair := range expected {
		if _, ok := pairs[pair]; !ok {
			t.Errorf("Expected pair '%s' in %s", pair, line)
		}
	}
	if expectedLen := len(expected); expectedLen != len(pairs) {
		t.Errorf("Expected %d pairs in %s", expectedLen, line)
	}
}

func AssertString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected %s\nGot: %s", expected, actual)
	}
}

func buildLogMap(d Data) (map[string]bool, string) {
	m := make(map[string]bool)
	parts := BuildLogParts(d, false)
	for _, pair := range parts {
		m[pair] = true
	}
	return m, strings.Join(parts, space)
}

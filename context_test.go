package grohl

import (
	"testing"
)

func TestContextMerge(t *testing.T) {
	orig := NewContext(nil)

	orig.Add("a", 1)
	orig.Add("b", 1)

	merged := orig.Merge(Data{"b": 2, "c": 3})

	AssertData(t, merged, "a=1", "b=2", "c=3")
	AssertLog(t, orig, "a=1", "b=1")
}

func TestContextStatterPrefix(t *testing.T) {
	ctx1 := NewContext(nil)
	ctx2 := NewContext(nil)
	ctx3 := ctx1.New(nil)
	AssertString(t, "", ctx1.StatterBucket)
	AssertString(t, "", ctx2.StatterBucket)
	AssertString(t, "", ctx3.StatterBucket)

	ctx1.SetStatter(nil, 1.0, "abc")
	AssertString(t, "abc", ctx1.StatterBucket)
	AssertString(t, "", ctx2.StatterBucket)
	AssertString(t, "", ctx3.StatterBucket)
}

//go:build goexperiment.jsonv2

package jsonpointer

import (
	"bytes"
	"encoding/json/jsontext"
	"testing"
)

func TestGetDecode(t *testing.T) {
	t.Parallel()

	ptr := "/A/2/B"

	value := jsontext.Value(`{"A":[{},{},{"B":"C"}]}`)
	dec := jsontext.NewDecoder(bytes.NewBuffer(value))

	result, err := GetDecode(ptr, dec)
	if string(result) != `"C"` || err != nil {
		t.Fatalf("GetDecode() = (%v, %v), want (\"C\", <nil>)", result, err)
	}
}

func TestPointerGetDecode(t *testing.T) {
	t.Parallel()

	ptr, err := Parse("/A/2/B")
	if err != nil {
		t.Fatalf("Parse() = %v, want <nil>", err)
	}

	value := jsontext.Value(`{"A":[{},{},{"B":"C"}]}`)
	dec := jsontext.NewDecoder(bytes.NewBuffer(value))

	result, err := ptr.GetDecode(dec)
	if string(result) != `"C"` || err != nil {
		t.Fatalf("Pointer.GetDecode() = (%v, %v), want (\"C\", <nil>)", result, err)
	}
}

func BenchmarkGetDecode(b *testing.B) {
	b.ReportAllocs()

	value := jsontext.Value(`{"A":[{},{},{"B":{"C":"D"}}]}`)
	buf := new(bytes.Buffer)
	dec := jsontext.NewDecoder(buf)

	for b.Loop() {
		buf.Reset()
		buf.Write(value)
		dec.Reset(buf)

		_, err := GetDecode("/A/2/B/C", dec)
		if err != nil {
			b.Fatalf("GetDecode() = %v, want <nil>", err)
		}
	}
}

func BenchmarkPointerGetDecode(b *testing.B) {
	b.ReportAllocs()

	value := jsontext.Value(`{"A":[{},{},{"B":{"C":"D"}}]}`)
	buf := new(bytes.Buffer)
	dec := jsontext.NewDecoder(buf)

	ptr, err := Parse("/A/2/B/C")
	if err != nil {
		b.Fatalf("Parse(/A/2/B/C) = %v, want <nil>", err)
	}

	for b.Loop() {
		buf.Reset()
		buf.Write(value)
		dec.Reset(buf)

		_, err := ptr.GetDecode(dec)
		if err != nil {
			b.Fatalf("Pointer.GetDecode() = %v, want <nil>", err)
		}
	}
}

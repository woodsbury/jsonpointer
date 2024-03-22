package jsonpointer

import (
	"bytes"
	"errors"
	"testing"
)

func TestGet(t *testing.T) {
	t.Parallel()

	ptr := "/A/2/B"

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": "C",
			},
		},
	}

	result, err := Get(ptr, value)
	if err != nil {
		t.Fatalf("Get() = %v, want <nil>", err)
	}

	if result != "C" {
		t.Errorf("Get() = %v, want C", result)
	}

	value = struct {
		A []map[string]any
	}{
		A: []map[string]any{
			nil,
			nil,
			{
				"B": "C",
			},
		},
	}

	result, err = Get(ptr, value)
	if err != nil {
		t.Fatalf("Get() = %v, want <nil>", err)
	}

	if result != "C" {
		t.Errorf("Get() = %v, want C", result)
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	ptrs := []string{
		"",
		"/",
		"//",
		"/~0",
		"/~1",
		"/~01",
		"/~10",
		"/0",
		"/01",
		"/1",
		"/a/b/c",
	}

	for _, ptr := range ptrs {
		_, err := Parse(ptr)
		if err != nil {
			t.Errorf("Parse(%s) = %v, want <nil>", ptr, err)
		}
	}

	var terr *invalidTokenError
	_, err := Parse("/~")
	if !errors.As(err, &terr) {
		t.Errorf("Parse(/~) = %v, want %v", err, &invalidTokenError{"~"})
	}

	var perr *invalidPointerError
	_, err = Parse("a")
	if !errors.As(err, &perr) {
		t.Errorf("Parse(a) = %v, want %v", err, &invalidPointerError{"a"})
	}
}

func TestPointerGet(t *testing.T) {
	t.Parallel()

	ptr := "/A/2/B"

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": "C",
			},
		},
	}

	p, err := Parse(ptr)
	if err != nil {
		t.Fatalf("Parse(%s) = %v, want <nil>", ptr, err)
	}

	result, err := p.Get(value)
	if err != nil {
		t.Fatalf("Pointer.Get() = %v, want <nil>", err)
	}

	if result != "C" {
		t.Errorf("Pointer.Get() = %v, want C", result)
	}

	value = struct {
		A []map[string]any
	}{
		A: []map[string]any{
			nil,
			nil,
			{
				"B": "C",
			},
		},
	}

	result, err = p.Get(value)
	if err != nil {
		t.Fatalf("Pointer.Get() = %v, want <nil>", err)
	}

	if result != "C" {
		t.Errorf("Pointer.Get() = %v, want C", result)
	}
}

func TestPointerMarshalText(t *testing.T) {
	t.Parallel()

	ptrs := [][]byte{
		[]byte(""),
		[]byte("/"),
		[]byte("//"),
		[]byte("/~0"),
		[]byte("/~1"),
		[]byte("/~01"),
		[]byte("/~10"),
		[]byte("/0"),
		[]byte("/01"),
		[]byte("/1"),
		[]byte("/a/b/c"),
	}

	for _, ptr := range ptrs {
		var p Pointer
		if err := p.UnmarshalText(ptr); err != nil {
			t.Fatalf("Pointer.UnmarshalText(%s) = %v, want <nil>", ptr, err)
		}

		d, err := p.MarshalText()
		if err != nil {
			t.Fatalf("Pointer.MarshalText() = %v, want <nil>", err)
		}

		if !bytes.Equal(d, ptr) {
			t.Errorf("Pointer.MarshalText() = %s, want %s", d, ptr)
		}
	}
}

func TestPointerString(t *testing.T) {
	t.Parallel()

	ptrs := []string{
		"",
		"/",
		"//",
		"/~0",
		"/~1",
		"/~01",
		"/~10",
		"/0",
		"/01",
		"/1",
		"/a/b/c",
	}

	for _, ptr := range ptrs {
		p, err := Parse(ptr)
		if err != nil {
			t.Fatalf("Parse(%s) = %v, want <nil>", ptr, err)
		}

		s := p.String()
		if s != ptr {
			t.Errorf("Pointer.String() = %s, want %s", s, ptr)
		}
	}
}

func BenchmarkGet(b *testing.B) {
	b.ReportAllocs()

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": "C",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get("/A/2/B", value)
		if err != nil {
			b.Fatalf("Get() = %v, want <nil>", err)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := Parse("/A/2/B")
		if err != nil {
			b.Fatalf("Parse() = %v, want <nil>", err)
		}
	}
}

func BenchmarkPointerGet(b *testing.B) {
	b.ReportAllocs()

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": "C",
			},
		},
	}

	ptr, err := Parse("/A/2/B")
	if err != nil {
		b.Fatalf("Parse(/A/2/B) = %v, want <nil>", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = ptr.Get(value)
		if err != nil {
			b.Fatalf("Pointer.Get() = %v, want <nil>", err)
		}
	}
}

func BenchmarkPointerMarshalText(b *testing.B) {
	b.ReportAllocs()

	p, err := Parse("/A/2/B")
	if err != nil {
		b.Fatalf("Parse() = %v, want <nil>", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.MarshalText()
		if err != nil {
			b.Fatalf("Pointer.MarshalText() = %v, want <nil>", err)
		}
	}
}

func BenchmarkPointerString(b *testing.B) {
	b.ReportAllocs()

	p, err := Parse("/A/2/B")
	if err != nil {
		b.Fatalf("Parse() = %v, want <nil>", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = p.String()
	}
}

func BenchmarkPointerUnmarshalText(b *testing.B) {
	b.ReportAllocs()

	data := []byte("/A/2/B")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var p Pointer
		err := p.UnmarshalText(data)
		if err != nil {
			b.Fatalf("Pointer.UnmarshalText() = %v, want <nil>", err)
		}
	}
}

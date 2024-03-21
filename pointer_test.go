package jsonpointer

import (
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

	b.ResetTimer()
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

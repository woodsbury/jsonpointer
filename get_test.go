package jsonpointer

import "testing"

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
	if result != "C" || err != nil {
		t.Fatalf("Get() = (%v, %v), want (C, <nil>)", result, err)
	}

	type B struct {
		B string
	}

	type A struct {
		A []B
	}

	value = &A{
		A: []B{
			{},
			{},
			{
				B: "C",
			},
		},
	}

	result, err = Get(ptr, value)
	if result != "C" || err != nil {
		t.Fatalf("Get() = (%v, %v), want (C, <nil>)", result, err)
	}

	type E struct {
		E string `json:"B"`
	}

	type D struct {
		D []any `json:"A"`
	}

	value = &D{
		D: []any{
			&E{},
			&E{},
			&E{
				E: "C",
			},
		},
	}

	result, err = Get(ptr, value)
	if result != "C" || err != nil {
		t.Fatalf("Get() = (%v, %v), want (C, <nil>)", result, err)
	}

	ptr = "/A/2"

	value = map[string]any{
		"A": []any{
			nil,
			nil,
			nil,
		},
	}

	result, err = Get(ptr, value)
	if result != nil || err != nil {
		t.Fatalf("Get() = (%v, %v), want (<nil>, <nil>)", result, err)
	}

	value = &D{
		D: []any{
			nil,
			nil,
			nil,
		},
	}

	result, err = Get(ptr, value)
	if result != nil || err != nil {
		t.Fatalf("Get() = (%v, %v), want (<nil>, <nil>)", result, err)
	}
}

func TestPointerGet(t *testing.T) {
	t.Parallel()

	ptr, err := Parse("/A/2/B")
	if err != nil {
		t.Fatalf("Parse(/A/2/B) = %v, want <nil>", err)
	}

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": "C",
			},
		},
	}

	result, err := ptr.Get(value)
	if result != "C" || err != nil {
		t.Fatalf("Get() = (%v, %v), want (C, <nil>)", result, err)
	}

	type B struct {
		B string
	}

	type A struct {
		A []B
	}

	value = &A{
		A: []B{
			{},
			{},
			{
				B: "C",
			},
		},
	}

	result, err = ptr.Get(value)
	if result != "C" || err != nil {
		t.Fatalf("Get() = (%v, %v), want (C, <nil>)", result, err)
	}

	type E struct {
		E string `json:"B"`
	}

	type D struct {
		D []any `json:"A"`
	}

	value = &D{
		D: []any{
			&E{},
			&E{},
			&E{
				E: "C",
			},
		},
	}

	result, err = ptr.Get(value)
	if result != "C" || err != nil {
		t.Fatalf("Get() = (%v, %v), want (C, <nil>)", result, err)
	}

	ptr, err = Parse("/A/2")
	if err != nil {
		t.Fatalf("Parse(/A/2) = %v, want <nil>", err)
	}

	value = map[string]any{
		"A": []any{
			nil,
			nil,
			nil,
		},
	}

	result, err = ptr.Get(value)
	if result != nil || err != nil {
		t.Fatalf("Get() = (%v, %v), want (<nil>, <nil>)", result, err)
	}

	value = &D{
		D: []any{
			nil,
			nil,
			nil,
		},
	}

	result, err = ptr.Get(value)
	if result != nil || err != nil {
		t.Fatalf("Get() = (%v, %v), want (<nil>, <nil>)", result, err)
	}
}

func BenchmarkGetMap(b *testing.B) {
	b.ReportAllocs()

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": map[string]any{
					"C": "D",
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get("/A/2/B/C", value)
		if err != nil {
			b.Fatalf("Get() = %v, want <nil>", err)
		}
	}
}

func BenchmarkGetStruct(b *testing.B) {
	b.ReportAllocs()

	type C struct {
		C string
	}

	type B struct {
		B C
	}

	type A struct {
		A []B
	}

	value := &A{
		A: []B{
			{},
			{},
			{
				B: C{
					C: "D",
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Get("/A/2/B/C", value)
		if err != nil {
			b.Fatalf("Get() = %v, want <nil>", err)
		}
	}
}

func BenchmarkPointerGetMap(b *testing.B) {
	b.ReportAllocs()

	var value any = map[string]any{
		"A": []any{
			map[string]any{},
			map[string]any{},
			map[string]any{
				"B": map[string]any{
					"C": "D",
				},
			},
		},
	}

	ptr, err := Parse("/A/2/B/C")
	if err != nil {
		b.Fatalf("Parse(/A/2/B/C) = %v, want <nil>", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = ptr.Get(value)
		if err != nil {
			b.Fatalf("Pointer.Get() = %v, want <nil>", err)
		}
	}
}

func BenchmarkPointerGetStruct(b *testing.B) {
	b.ReportAllocs()

	type C struct {
		C string
	}

	type B struct {
		B C
	}

	type A struct {
		A []B
	}

	value := &A{
		A: []B{
			{},
			{},
			{
				B: C{
					C: "D",
				},
			},
		},
	}

	ptr, err := Parse("/A/2/B/C")
	if err != nil {
		b.Fatalf("Parse(/A/2/B/C) = %v, want <nil>", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err = ptr.Get(value)
		if err != nil {
			b.Fatalf("Pointer.Get() = %v, want <nil>", err)
		}
	}
}

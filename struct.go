package jsonpointer

import (
	"reflect"
	"strings"
	"sync"
	"unicode"
)

func structField(field string, value *reflect.Value) bool {
	fields := getStructFields(value.Type())
	i, ok := fields[field]
	if !ok {
		return false
	}

	*value = value.FieldByIndex(i)
	return true
}

type structFields map[string][]int

var structFieldsCache sync.Map

func getStructFields(t reflect.Type) structFields {
	if fields, ok := structFieldsCache.Load(t); ok {
		return fields.(structFields)
	}

	type field struct {
		t reflect.Type
		i []int
	}

	fields := make(structFields, t.NumField())

	current := []field{}
	next := []field{{
		t: t,
	}}

	visited := make(map[reflect.Type]struct{})

	for len(next) > 0 {
		current, next = next, current[:0]

		for _, f := range current {
			if _, ok := visited[f.t]; ok {
				continue
			}

			visited[f.t] = struct{}{}

			n := f.t.NumField()
			for i := 0; i < n; i++ {
				sf := f.t.Field(i)
				if sf.Anonymous {
					ft := sf.Type
					if ft.Kind() == reflect.Pointer {
						ft = ft.Elem()
					}

					if !sf.IsExported() && ft.Kind() != reflect.Struct {
						continue
					}
				} else if !sf.IsExported() {
					continue
				}

				name := sf.Name
				tag := sf.Tag.Get("json")
				if tag != "" {
					if tag == "-" {
						continue
					}

					tag, _, _ = strings.Cut(tag, ",")
					for _, r := range tag {
						if strings.ContainsRune("!#$%&()*+-./:;<=>?@[]^_{|}~ ", r) {
							continue
						}

						if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
							tag = ""
							break
						}
					}

					if tag != "" {
						name = tag
					}
				}

				index := make([]int, len(f.i)+1)
				copy(index, f.i)
				index[len(f.i)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					ft = ft.Elem()
				}

				if !sf.Anonymous || ft.Kind() != reflect.Struct {
					fields[name] = index
					continue
				}

				next = append(next, field{
					t: ft,
					i: index,
				})
			}
		}
	}

	fieldsVal, _ := structFieldsCache.LoadOrStore(t, fields)
	return fieldsVal.(structFields)
}

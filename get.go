package jsonpointer

import (
	"reflect"
	"strings"
)

// Get resolves the JSON pointer ptr against value and returns the result.
func Get(ptr string, value any) (any, error) {
	if ptr == "" {
		return value, nil
	}

	if ptr[0] != '/' {
		return nil, &invalidPointerError{ptr}
	}

	remaining := ptr[1:]
	count := strings.Count(remaining, "/")
	result := value
	var i int
	var tok token
	var ok bool
	for i = 0; i < count; i++ {
		next := strings.IndexByte(remaining, '/')
		if next == -1 {
			return nil, &invalidPointerError{ptr}
		}

		var err error
		tok, err = parseToken(remaining[:next])
		if err != nil {
			return nil, err
		}

		result, ok, err = get(tok, result)
		if err != nil {
			return nil, err
		}

		remaining = remaining[next+1:]

		if !ok {
			i++
			break
		}
	}

	if ok {
		tok, err := parseToken(remaining)
		if err != nil {
			return nil, err
		}

		result, ok, err = get(tok, result)
		if err != nil {
			return nil, err
		}

		if ok {
			return result, nil
		}
	}

	refResult := reflect.ValueOf(result)
	if err := getReflect(tok, &refResult); err != nil {
		return nil, err
	}

	for ; i < count; i++ {
		next := strings.IndexByte(remaining, '/')
		if next == -1 {
			return nil, &invalidPointerError{ptr}
		}

		var err error
		tok, err = parseToken(remaining[:next])
		if err != nil {
			return nil, err
		}

		if err := getReflect(tok, &refResult); err != nil {
			return nil, err
		}

		remaining = remaining[next+1:]
	}

	var err error
	tok, err = parseToken(remaining)
	if err != nil {
		return nil, err
	}

	if err := getReflect(tok, &refResult); err != nil {
		return nil, err
	}

	return refResult.Interface(), nil
}

// Get resolves the JSON pointer parsed into p against value and returns the
// result.
func (p Pointer) Get(value any) (any, error) {
	result := value

	var i int
	var tok token
	var ok bool
	for i, tok = range p.tokens {
		var err error
		result, ok, err = get(tok, result)
		if err != nil {
			return nil, err
		}

		if !ok {
			break
		}
	}

	if ok {
		return result, nil
	}

	refResult := reflect.ValueOf(result)
	for i, tok = range p.tokens[i:] {
		if err := getReflect(tok, &refResult); err != nil {
			return nil, err
		}
	}

	return refResult.Interface(), nil
}

func get(tok token, value any) (any, bool, error) {
	switch v := value.(type) {
	case map[string]any:
		field, ok := v[tok.field]
		if !ok {
			return nil, false, &valueNotFoundError{tok.field}
		}

		return field, true, nil
	case *map[string]any:
		field, ok := (*v)[tok.field]
		if !ok {
			return nil, false, &valueNotFoundError{tok.field}
		}

		return field, true, nil
	case []any:
		if tok.index == -1 {
			if tok.field == "-" {
				return nil, false, &arrayIndexOutOfBoundsError{len(v)}
			}

			return nil, false, &invalidArrayIndexError{tok.field}
		}

		if tok.index >= len(v) {
			return nil, false, &arrayIndexOutOfBoundsError{tok.index}
		}

		return v[tok.index], true, nil
	case *[]any:
		if tok.index == -1 {
			if tok.field == "-" {
				return nil, false, &arrayIndexOutOfBoundsError{len(*v)}
			}

			return nil, false, &invalidArrayIndexError{tok.field}
		}

		if tok.index >= len(*v) {
			return nil, false, &arrayIndexOutOfBoundsError{len(*v)}
		}

		return (*v)[tok.index], true, nil
	case *any:
		switch v := (*v).(type) {
		case map[string]any:
			field, ok := v[tok.field]
			if !ok {
				return nil, false, &valueNotFoundError{tok.field}
			}

			return field, true, nil
		case []any:
			if tok.index == -1 {
				if tok.field == "-" {
					return nil, false, &arrayIndexOutOfBoundsError{len(v)}
				}

				return nil, false, &invalidArrayIndexError{tok.field}
			}

			if tok.index >= len(v) {
				return nil, false, &arrayIndexOutOfBoundsError{tok.index}
			}

			return v[tok.index], true, nil
		case nil:
			return nil, false, &valueNotFoundError{tok.field}
		}
	case nil:
		return nil, false, &valueNotFoundError{tok.field}
	}

	return value, false, nil
}

func getReflect(tok token, value *reflect.Value) error {
	k := value.Kind()
	for {
		if k != reflect.Interface && k != reflect.Pointer {
			break
		}

		if value.IsNil() {
			return &valueNotFoundError{tok.field}
		}

		*value = value.Elem()
		k = value.Kind()
	}

	switch k {
	case reflect.Array, reflect.Slice:
		if tok.index == -1 {
			if tok.field == "-" {
				return &arrayIndexOutOfBoundsError{value.Len()}
			}

			return &invalidArrayIndexError{tok.field}
		}

		if tok.index >= value.Len() {
			return &arrayIndexOutOfBoundsError{tok.index}
		}

		*value = value.Index(tok.index)
		return nil
	case reflect.Map:
		*value = value.MapIndex(reflect.ValueOf(tok.field))
		if !value.IsValid() {
			return &valueNotFoundError{tok.field}
		}

		return nil
	case reflect.Struct:
		if ok := structField(tok.field, value); !ok {
			return &valueNotFoundError{tok.field}
		}

		return nil
	default:
		return &valueNotFoundError{tok.field}
	}
}

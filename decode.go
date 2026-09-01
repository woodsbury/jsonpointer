//go:build goexperiment.jsonv2

package jsonpointer

import (
	"encoding/json/jsontext"
	"strings"
)

// Get resolves the JSON pointer ptr against the value being decoded by dec and
// returns the result.
func GetDecode(ptr string, dec *jsontext.Decoder) (jsontext.Value, error) {
	if ptr == "" {
		value, err := dec.ReadValue()
		if err != nil {
			return nil, &jsonDecodeError{err}
		}

		return value, nil
	}

	if ptr[0] != '/' {
		return nil, &invalidPointerError{ptr}
	}

	remaining := ptr[1:]
	count := strings.Count(remaining, "/")
	for range count {
		next := strings.IndexByte(remaining, '/')
		if next == -1 {
			return nil, &invalidPointerError{ptr}
		}

		tok, err := parseToken(remaining[:next])
		if err != nil {
			return nil, err
		}

		if err := getDecode(tok, dec); err != nil {
			return nil, err
		}

		remaining = remaining[next+1:]
	}

	tok, err := parseToken(remaining)
	if err != nil {
		return nil, err
	}

	if err := getDecode(tok, dec); err != nil {
		return nil, err
	}

	value, err := dec.ReadValue()
	if err != nil {
		return nil, &jsonDecodeError{err}
	}

	return value, nil
}

// GetDecode resolves the JSON pointer parsed into p against the value being
// decoded by dec and returns the result.
func (p Pointer) GetDecode(dec *jsontext.Decoder) (jsontext.Value, error) {
	for _, tok := range p.tokens {
		if err := getDecode(tok, dec); err != nil {
			return nil, err
		}
	}

	value, err := dec.ReadValue()
	if err != nil {
		return nil, &jsonDecodeError{err}
	}

	return value, nil
}

func getDecode(tok token, dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindBeginArray:
		if tok.index == -1 {
			if tok.field == "-" {
				return &arrayIndexOutOfBoundsError{-1}
			}

			return &invalidArrayIndexError{tok.field}
		}

		if _, err := dec.ReadToken(); err != nil {
			return &jsonDecodeError{err}
		}

		for i := 0; i < tok.index; i++ {
			if dec.PeekKind() == jsontext.KindEndArray {
				return &arrayIndexOutOfBoundsError{tok.index}
			}

			if err := dec.SkipValue(); err != nil {
				return &jsonDecodeError{err}
			}
		}
	case jsontext.KindBeginObject:
		if _, err := dec.ReadToken(); err != nil {
			return &jsonDecodeError{err}
		}

		for {
			if dec.PeekKind() == jsontext.KindEndObject {
				return &valueNotFoundError{tok.field}
			}

			key, err := dec.ReadToken()
			if err != nil {
				return &jsonDecodeError{err}
			}

			if key.String() == tok.field {
				return nil
			}

			if err := dec.SkipValue(); err != nil {
				return &jsonDecodeError{err}
			}
		}
	default:
		return &valueNotFoundError{tok.field}
	}

	return nil
}

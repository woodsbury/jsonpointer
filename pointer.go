package jsonpointer

import (
	"bytes"
	"strconv"
	"strings"
	"unsafe"
)

// Pointer represents a parsed JSON pointer. Parsing a JSON pointer allows it
// to be evaluated multiple times more efficiently.
type Pointer struct {
	tokens []token
}

// MustParse is like [Parse] but panics if the provided JSON pointer cannot be
// parsed, instead of returning an error.
func MustParse(ptr string) Pointer {
	p, err := Parse(ptr)
	if err != nil {
		panic("jsonpointer.MustParse(" + strconv.Quote(ptr) + "): invalid pointer")
	}

	return p
}

// Parse parses the JSON pointer ptr.
func Parse(ptr string) (Pointer, error) {
	if ptr == "" {
		return Pointer{}, nil
	}

	if ptr[0] != '/' {
		return Pointer{}, &invalidPointerError{ptr}
	}

	remaining := ptr[1:]
	count := strings.Count(remaining, "/")
	tokens := make([]token, count+1)
	for i := range count {
		next := strings.IndexByte(remaining, '/')
		if next == -1 {
			return Pointer{}, &invalidPointerError{ptr}
		}

		tok, err := parseToken(remaining[:next])
		if err != nil {
			return Pointer{}, err
		}

		tokens[i] = tok
		remaining = remaining[next+1:]
	}

	tok, err := parseToken(remaining)
	if err != nil {
		return Pointer{}, err
	}

	tokens[count] = tok

	return Pointer{
		tokens: tokens,
	}, nil
}

// AppendText implements the [encoding.TextAppender] interface.
func (p Pointer) AppendText(buf []byte) ([]byte, error) {
	if buf == nil {
		n := len(p.tokens)
		for _, tok := range p.tokens {
			n += len(tok.field)
		}

		buf = make([]byte, 0, n)
	}

	for _, tok := range p.tokens {
		buf = append(buf, '/')

		i := strings.IndexByte(tok.field, '~')
		j := strings.IndexByte(tok.field, '/')
		if i == -1 && j == -1 {
			buf = append(buf, tok.field...)
			continue
		}

		var k int
		if i != -1 && (j == -1 || i < j) {
			k = i
		} else {
			k = j
		}

		buf = append(buf, tok.field[:k]...)
		remaining := tok.field[k+1:]
		for {
			if i != -1 {
				buf = append(buf, '~', '0')
			} else {
				buf = append(buf, '~', '1')
			}

			i = strings.IndexByte(remaining, '~')
			j = strings.IndexByte(remaining, '/')
			if i == -1 && j == -1 {
				buf = append(buf, remaining...)
				break
			}

			var k int
			if i != -1 && i < j {
				k = i
			} else {
				k = j
			}

			buf = append(buf, remaining[:k]...)
			remaining = remaining[k+1:]
		}
	}

	return buf, nil
}

// Equal compares two Pointers and reports whether they are equal.
func (p Pointer) Equal(o Pointer) bool {
	l := len(p.tokens)
	if l != len(o.tokens) {
		return false
	}

	for i := range l {
		if p.tokens[i] != o.tokens[i] {
			return false
		}
	}

	return true
}

// IsZero reports whether the Pointer is the zero value. A zero Pointer value
// resolves against the root of a value.
func (p Pointer) IsZero() bool {
	return len(p.tokens) == 0
}

// MarshalText implements the [encoding.TextMarshaler] interface.
func (p Pointer) MarshalText() ([]byte, error) {
	return p.AppendText(nil)
}

// NumTokens returns the number of reference tokens in the Pointer value.
func (p Pointer) NumTokens() int {
	return len(p.tokens)
}

// String returns a string representation of the Pointer value.
func (p Pointer) String() string {
	n := len(p.tokens)
	for _, tok := range p.tokens {
		n += len(tok.field)
	}

	buf := make([]byte, 0, n)

	for _, tok := range p.tokens {
		buf = append(buf, '/')

		i := strings.IndexByte(tok.field, '~')
		j := strings.IndexByte(tok.field, '/')
		if i == -1 && j == -1 {
			buf = append(buf, tok.field...)
			continue
		}

		var k int
		if i != -1 && (j == -1 || i < j) {
			k = i
		} else {
			k = j
		}

		buf = append(buf, tok.field[:k]...)
		remaining := tok.field[k+1:]
		for {
			if i != -1 {
				buf = append(buf, '~', '0')
			} else {
				buf = append(buf, '~', '1')
			}

			i = strings.IndexByte(remaining, '~')
			j = strings.IndexByte(remaining, '/')
			if i == -1 && j == -1 {
				buf = append(buf, remaining...)
				break
			}

			var k int
			if i != -1 && i < j {
				k = i
			} else {
				k = j
			}

			buf = append(buf, remaining[:k]...)
			remaining = remaining[k+1:]
		}
	}

	return unsafe.String(unsafe.SliceData(buf), len(buf))
}

// Token returns the reference token at the index i. Token returns an empty
// string if the requested index doesn't exist.
func (p Pointer) Token(i int) string {
	if i < 0 || i >= len(p.tokens) {
		return ""
	}

	return p.tokens[i].field
}

// Tokens returns the reference tokens of the Pointer value.
func (p Pointer) Tokens() []string {
	tokens := make([]string, len(p.tokens))
	for i, tok := range p.tokens {
		tokens[i] = tok.field
	}

	return tokens
}

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (p *Pointer) UnmarshalText(data []byte) error {
	if len(data) == 0 {
		*p = Pointer{}
		return nil
	}

	if data[0] != '/' {
		return &invalidPointerError{string(data)}
	}

	remaining := data[1:]
	count := bytes.Count(remaining, data[0:1])

	var tokens []token
	if cap(p.tokens) < count+1 {
		tokens = make([]token, count+1)
	} else {
		tokens = p.tokens[:count+1]
	}

	for i := range count {
		next := bytes.IndexByte(remaining, '/')
		if next == -1 {
			return &invalidPointerError{string(data)}
		}

		tok, err := parseTokenBytes(remaining[:next])
		if err != nil {
			return err
		}

		tokens[i] = tok
		remaining = remaining[next+1:]
	}

	tok, err := parseTokenBytes(remaining)
	if err != nil {
		return err
	}

	tokens[count] = tok

	p.tokens = tokens
	return nil
}

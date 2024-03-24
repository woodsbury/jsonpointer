package jsonpointer

import (
	"bytes"
	"strconv"
	"strings"
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
	for i := 0; i < count; i++ {
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

// MarshalText implements the [encoding.TextMarshaler] interface.
func (p Pointer) MarshalText() ([]byte, error) {
	n := len(p.tokens)
	for _, tok := range p.tokens {
		n += len(tok.field)
	}

	var b bytes.Buffer
	b.Grow(n)

	for _, tok := range p.tokens {
		b.WriteByte('/')

		i := strings.IndexByte(tok.field, '~')
		j := strings.IndexByte(tok.field, '/')
		if i == -1 && j == -1 {
			b.WriteString(tok.field)
			continue
		}

		var k int
		if i != -1 && (j == -1 || i < j) {
			k = i
		} else {
			k = j
		}

		b.WriteString(tok.field[:k])
		remaining := tok.field[k+1:]
		for {
			if i != -1 {
				b.WriteString("~0")
			} else {
				b.WriteString("~1")
			}

			i = strings.IndexByte(remaining, '~')
			j = strings.IndexByte(remaining, '/')
			if i == -1 && j == -1 {
				b.WriteString(remaining)
				break
			}

			var k int
			if i != -1 && i < j {
				k = i
			} else {
				k = j
			}

			b.WriteString(remaining[:k])
			remaining = remaining[k+1:]
		}
	}

	return b.Bytes(), nil
}

// String returns a string representation of the Pointer value.
func (p Pointer) String() string {
	n := len(p.tokens)
	for _, tok := range p.tokens {
		n += len(tok.field)
	}

	var b strings.Builder
	b.Grow(n)

	for _, tok := range p.tokens {
		b.WriteByte('/')

		i := strings.IndexByte(tok.field, '~')
		j := strings.IndexByte(tok.field, '/')
		if i == -1 && j == -1 {
			b.WriteString(tok.field)
			continue
		}

		var k int
		if i != -1 && (j == -1 || i < j) {
			k = i
		} else {
			k = j
		}

		b.WriteString(tok.field[:k])
		remaining := tok.field[k+1:]
		for {
			if i != -1 {
				b.WriteString("~0")
			} else {
				b.WriteString("~1")
			}

			i = strings.IndexByte(remaining, '~')
			j = strings.IndexByte(remaining, '/')
			if i == -1 && j == -1 {
				b.WriteString(remaining)
				break
			}

			var k int
			if i != -1 && i < j {
				k = i
			} else {
				k = j
			}

			b.WriteString(remaining[:k])
			remaining = remaining[k+1:]
		}
	}

	return b.String()
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

	for i := 0; i < count; i++ {
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

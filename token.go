package jsonpointer

import (
	"bytes"
	"strconv"
	"strings"
)

type token struct {
	field string
	index int
}

func parseToken(tok string) (token, error) {
	if len(tok) == 0 {
		return token{
			index: -1,
		}, nil
	}

	i := strings.IndexByte(tok, '~')
	if i == -1 {
		if tok == "0" {
			return token{
				field: "0",
			}, nil
		}

		if r := tok[0]; r >= '1' && r <= '9' {
			return token{
				field: tok,
				index: atoi(tok),
			}, nil
		}

		return token{
			field: tok,
			index: -1,
		}, nil
	}

	var b strings.Builder
	b.Grow(len(tok))
	b.WriteString(tok[:i])

	remaining := tok[i:]
	for {
		if len(remaining) < 2 {
			return token{}, &invalidTokenError{tok}
		}

		switch tok[1] {
		case '0':
			b.WriteByte('~')
		case '1':
			b.WriteByte('/')
		default:
			return token{}, &invalidTokenError{tok}
		}

		remaining = remaining[2:]
		i = strings.IndexByte(remaining, '~')
		if i == -1 {
			b.WriteString(remaining)

			return token{
				field: b.String(),
				index: -1,
			}, nil
		}

		b.WriteString(remaining[:i])
		remaining = remaining[i:]
	}
}

func parseTokenBytes(tok []byte) (token, error) {
	if len(tok) == 0 {
		return token{
			index: -1,
		}, nil
	}

	i := bytes.IndexByte(tok, '~')
	if i == -1 {
		if string(tok) == "0" {
			return token{
				field: "0",
			}, nil
		}

		tokStr := string(tok)

		if r := tok[0]; r >= '1' && r <= '9' {
			index, err := strconv.Atoi(tokStr)
			if err == nil {
				return token{
					field: tokStr,
					index: index,
				}, nil
			}
		}

		return token{
			field: tokStr,
			index: -1,
		}, nil
	}

	var b bytes.Buffer
	b.Grow(len(tok))
	b.Write(tok[:i])

	remaining := tok[i:]
	for {
		if len(remaining) < 2 {
			return token{}, &invalidTokenError{string(tok)}
		}

		switch tok[1] {
		case '0':
			b.WriteByte('~')
		case '1':
			b.WriteByte('/')
		default:
			return token{}, &invalidTokenError{string(tok)}
		}

		remaining = remaining[2:]
		i = bytes.IndexByte(remaining, '~')
		if i == -1 {
			b.Write(remaining)

			return token{
				field: b.String(),
				index: -1,
			}, nil
		}

		b.Write(remaining[:i])
		remaining = remaining[i:]
	}
}

func atoi(s string) int {
	var n int
	for _, r := range s {
		t := n*10 + int(r-'0')
		if t < n {
			return -1
		}

		n = t
	}

	return n
}

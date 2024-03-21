package jsonpointer

import (
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
			index, err := strconv.Atoi(tok)
			if err == nil {
				return token{
					field: tok,
					index: index,
				}, nil
			}
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
			b.WriteRune('~')
		case '1':
			b.WriteRune('/')
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

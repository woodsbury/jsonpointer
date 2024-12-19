package jsonpointer

import (
	"errors"
	"testing"
)

func TestParseToken(t *testing.T) {
	t.Parallel()

	type test struct {
		tok   string
		field string
		index int
	}

	tests := []test{
		{"", "", -1},
		{"a", "a", -1},
		{"0", "0", 0},
		{"1", "1", 1},
		{"1a", "1a", -1},
		{"01", "01", -1},
		{"~0", "~", -1},
		{"~1", "/", -1},
		{"~01", "~1", -1},
		{"~10", "/0", -1},
	}

	for _, test := range tests {
		tok, err := parseToken(test.tok)
		if err != nil {
			t.Errorf("parseToken(%s) = %v, want <nil>", test.tok, err)
		}

		if tok.field != test.field || tok.index != test.index {
			t.Errorf("parseToken(%s) = {%s, %d}, want {%s, %d}", test.tok, tok.field, tok.index, test.field, test.index)
		}

		tok, err = parseTokenBytes([]byte(test.tok))
		if err != nil {
			t.Errorf("parseTokenBytes(%s) = %v, want <nil>", test.tok, err)
		}

		if tok.field != test.field || tok.index != test.index {
			t.Errorf("parseTokenBytes(%s) = {%s, %d}, want {%s, %d}", test.tok, tok.field, tok.index, test.field, test.index)
		}
	}

	var terr *invalidTokenError
	_, err := parseToken("~")
	if !errors.As(err, &terr) {
		t.Errorf("parseToken(~) = %v, want %v", err, &invalidTokenError{"~"})
	}

	_, err = parseTokenBytes([]byte("~"))
	if !errors.As(err, &terr) {
		t.Errorf("parseTokenBytes(~) = %v, want %v", err, &invalidTokenError{"~"})
	}

	_, err = parseToken("~2")
	if !errors.As(err, &terr) {
		t.Errorf("parseToken(~2) = %v, want %v", err, &invalidTokenError{"~2"})
	}

	_, err = parseTokenBytes([]byte("~2"))
	if !errors.As(err, &terr) {
		t.Errorf("parseTokenBytes(~2) = %v, want %v", err, &invalidTokenError{"~2"})
	}
}

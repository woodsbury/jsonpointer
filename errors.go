package jsonpointer

import "strconv"

type arrayIndexOutOfBoundsError struct {
	index int
}

func (err *arrayIndexOutOfBoundsError) Error() string {
	return "jsonpointer: array index out of bounds " + strconv.Itoa(err.index)
}

type invalidArrayIndexError struct {
	tok string
}

func (err *invalidArrayIndexError) Error() string {
	return "jsonpointer: invalid array index " + strconv.QuoteToASCII(err.tok)
}

type invalidPointerError struct {
	ptr string
}

func (err *invalidPointerError) Error() string {
	return "jsonpointer: invalid pointer " + strconv.QuoteToASCII(err.ptr)
}

type invalidTokenError struct {
	tok string
}

func (err *invalidTokenError) Error() string {
	return "jsonpointer: invalid token in pointer " + strconv.QuoteToASCII(err.tok)
}

type valueNotFoundError struct {
	tok string
}

func (err *valueNotFoundError) Error() string {
	return "jsonpointer: value not found " + strconv.QuoteToASCII(err.tok)
}

package jsonpointer

import (
	"errors"
	"strconv"
)

var (
	ErrArrayIndexOutOfBounds = errors.New("jsonpointer: array index out of bounds")
	ErrInvalidArrayIndex     = errors.New("jsonpointer: invalid array index")
	ErrInvalidPointer        = errors.New("jsonpointer: invalid pointer")
	ErrValueNotFound         = errors.New("jsonpointer: value not found")
)

type arrayIndexOutOfBoundsError struct {
	index int
}

func (err *arrayIndexOutOfBoundsError) Error() string {
	return "jsonpointer: array index out of bounds " + strconv.Itoa(err.index)
}

func (err *arrayIndexOutOfBoundsError) Is(target error) bool {
	return target == ErrArrayIndexOutOfBounds
}

type invalidArrayIndexError struct {
	tok string
}

func (err *invalidArrayIndexError) Error() string {
	return "jsonpointer: invalid array index " + strconv.QuoteToASCII(err.tok)
}

func (err *invalidArrayIndexError) Is(target error) bool {
	return target == ErrInvalidArrayIndex
}

type invalidPointerError struct {
	ptr string
}

func (err *invalidPointerError) Error() string {
	return "jsonpointer: invalid pointer " + strconv.QuoteToASCII(err.ptr)
}

func (err *invalidPointerError) Is(target error) bool {
	return target == ErrInvalidPointer
}

type invalidTokenError struct {
	tok string
}

func (err *invalidTokenError) Error() string {
	return "jsonpointer: invalid token in pointer " + strconv.QuoteToASCII(err.tok)
}

func (err *invalidTokenError) Is(target error) bool {
	return target == ErrInvalidPointer
}

type valueNotFoundError struct {
	tok string
}

func (err *valueNotFoundError) Error() string {
	return "jsonpointer: value not found " + strconv.QuoteToASCII(err.tok)
}

func (err *valueNotFoundError) Is(target error) bool {
	return target == ErrValueNotFound
}

package main

import (
	"fmt"
	"strings"
)

// ErrorType is the type for lexical error types
type ErrorType string

const (
	// ErrorSyntax is syntax error type
	ErrorSyntax ErrorType = "syntax"
	// ErrorArg is argument error type
	ErrorArg ErrorType = "argument"
	// ErrorValue is value error type
	ErrorValue ErrorType = "value"
	// ErrorUnknown is unknown error type
	ErrorUnknown ErrorType = "unknown"
)

// NewError creates a new Error
func NewError(typ ErrorType, pos uint, msg string) *Error {
	return &Error{
		typ:  typ,
		pos:  pos,
		msgs: []string{msg},
	}
}

// Error is lexical error struct
type Error struct {
	typ  ErrorType
	msgs []string
	pos  uint
}

func (e *Error) Message() string {
	return strings.Join(e.msgs, ": ")
}

// Error returns error string
func (e *Error) Error() string {
	return fmt.Sprintf(
		"%s error near index %d: %s",
		string(e.typ),
		e.pos,
		e.Message(),
	)
}

func (e *Error) SpacedError() string {
	return fmt.Sprintf(
		"%s^ %s error: %s",
		strings.Repeat(" ", int(e.pos)),
		string(e.typ),
		e.Message(),
	)
}

// AppendMsg add a message to the begining of current messages
func (e *Error) AppendMsg(msg string) {
	e.msgs = append(e.msgs, msg)
}

// PrependMsg add a message to the end of current messages
func (e *Error) PrependMsg(msg string) {
	e.msgs = append([]string{msg}, e.msgs...)
}

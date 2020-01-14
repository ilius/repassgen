package main

import (
	"fmt"
	"strings"
)

// LexErrorType is the type for lexical error types
type LexErrorType string

const (
	// LexErrorSyntax is syntax error type
	LexErrorSyntax LexErrorType = "syntax"
	// LexErrorValue is value error type
	LexErrorValue LexErrorType = "value"
	// LexErrorUnknown is unknown error type
	LexErrorUnknown LexErrorType = "unknown"
)

// NewError creates a new LexError
func NewError(typ LexErrorType, pos uint, msg string) *LexError {
	return &LexError{
		typ:  typ,
		pos:  pos,
		msgs: []string{msg},
	}
}

// LexError is lexical error struct
type LexError struct {
	typ  LexErrorType
	pos  uint
	msgs []string
}

// Error returns error string
func (e *LexError) Error() string {
	return fmt.Sprintf(
		"%s error near index %d: %s",
		string(e.typ),
		e.pos,
		strings.Join(e.msgs, ": "),
	)
}

// AppendMsg add a message to the begining of current messages
func (e *LexError) AppendMsg(msg string) {
	e.msgs = append(e.msgs, msg)
}

// PrependMsg add a message to the end of current messages
func (e *LexError) PrependMsg(msg string) {
	e.msgs = append([]string{msg}, e.msgs...)
}

package passgen

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
		typ:     typ,
		pos:     pos,
		markLen: 1,
		msgs:    []string{msg},
	}
}

// Error is lexical error struct
type Error struct {
	typ     ErrorType
	msgs    []string
	pos     uint
	markLen uint
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

func (e *Error) WithMarkLen(markLen int) *Error {
	if markLen == 0 {
		markLen = 1
	}
	e.markLen = uint(markLen)
	return e
}

func (e *Error) SpacedError() string {
	if e.markLen <= 0 {
		return fmt.Sprintf(
			"%s^ %s error: %s",
			strings.Repeat(" ", int(e.pos)),
			string(e.typ),
			e.Message(),
		)
	}
	return fmt.Sprintf(
		"%s%s %s error: %s",
		strings.Repeat(" ", int(e.pos-e.markLen+1)),
		strings.Repeat("^", int(e.markLen)),
		string(e.typ),
		e.Message(),
	)
}

func ParseSpacedError(str string) *Error {
	trimmed := strings.TrimLeft(str, " ")
	if trimmed[0] != '^' {
		return nil
	}
	remain := strings.TrimLeft(trimmed, "^")
	pos := len(str) - len(remain) - 1
	markLen := len(trimmed) - len(remain)
	parts := strings.SplitN(strings.TrimLeft(remain, " "), " ", 3)
	if parts[1] != "error:" {
		return nil
	}
	typ := parts[0]
	msgs := strings.Split(parts[2], ": ")
	return &Error{
		typ:     ErrorType(typ),
		pos:     uint(pos),
		markLen: uint(markLen),
		msgs:    msgs,
	}
}

// AppendMsg add a message to the beginning of current messages
func (e *Error) AppendMsg(msg string) {
	e.msgs = append(e.msgs, msg)
}

// PrependMsg add a message to the end of current messages
func (e *Error) PrependMsg(msg string) {
	e.msgs = append([]string{msg}, e.msgs...)
}

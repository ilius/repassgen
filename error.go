package main

import (
	"fmt"
	"strings"
)

type LexErrorType string

const (
	LexErrorSyntax  LexErrorType = "syntax"
	LexErrorValue   LexErrorType = "value"
	LexErrorUnknown LexErrorType = "unknown"
)

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

func (e *LexError) Error() string {
	return fmt.Sprintf(
		"%s error near index %d: %s",
		string(e.typ),
		e.pos,
		strings.Join(e.msgs, ": "),
	)
}

// MovePos
func (e *LexError) MovePos(offset int) {
	e.pos = uint(int(e.pos) + offset)
}

// AppendMsg
func (e *LexError) AppendMsg(msg string) {
	e.msgs = append(e.msgs, msg)
}

// PrependMsg
func (e *LexError) PrependMsg(msg string) {
	e.msgs = append([]string{msg}, e.msgs...)
}

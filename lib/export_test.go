package passgen

import "io"

type TestOutput struct {
	output    []byte
	maxLength int
}

func (out *TestOutput) Write(data []byte) (int, error) {
	if out.maxLength > 0 && len(out.output) > out.maxLength {
		return 0, io.EOF
	}
	out.output = append(out.output, data...)
	return len(data), nil
}

func (s *State) SetOutput(out io.Writer) {
	s.output = out
}

func (s *State) Output() string {
	out, ok := s.output.(*TestOutput)
	if !ok {
		panic("not a TestOutput")
	}
	return string(out.output)
}

func (e *Error) Pos() int {
	return int(e.pos)
}

func (e *Error) MarkLen() int {
	return int(e.markLen)
}

func (e *Error) Type() ErrorType {
	return e.typ
}

func NewAlterGenerator(parts [][]rune, indexList []uint64) *alterGenerator {
	return &alterGenerator{
		parts:     parts,
		indexList: indexList,
		absPos:    0,
	}
}

func NewByteGenerator(uppercase bool) *byteGenerator {
	return &byteGenerator{
		uppercase: uppercase,
	}
}

func NewEncoderFunctionCallGenerator(funcName string, argPattern []rune) *encoderFunctionCallGenerator {
	return &encoderFunctionCallGenerator{
		funcName:   funcName,
		argPattern: argPattern,
	}
}

func NewGroupGenerator(pattern []rune) *groupGenerator {
	return &groupGenerator{
		pattern: pattern,
	}
}

func NewDateGenerator(sep string, startJd int, endJd int) *dateGenerator {
	return &dateGenerator{
		sep:     sep,
		startJd: startJd,
		endJd:   endJd,
	}
}

func NewRepeatGenerator(child generatorIface, count int64) *repeatGenerator {
	return &repeatGenerator{
		child: child,
		count: count,
	}
}

func NewStaticStringGenerator(str []rune) *staticStringGenerator {
	return &staticStringGenerator{str: str}
}

func NewShuffleGenerator(argPattern []rune) *shuffleGenerator {
	return &shuffleGenerator{
		argPattern: argPattern,
	}
}

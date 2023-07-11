package passgen

func (s *State) Output() string {
	return string(s.output)
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

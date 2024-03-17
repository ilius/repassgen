package passgen

import "strconv"

// LexType is the type for lex functions
type LexType func(*State) (LexType, error)

var (
	lexRootUnicode     LexType
	lexRootUnicodeWide LexType
)

func init() {
	lexRootUnicode = makeLexUnicode(LexRoot, 'u', 6, false)
	lexRootUnicodeWide = makeLexUnicode(LexRoot, 'U', 10, false)
}

// LexRoot is the root lex implementation
func LexRoot(s *State) (LexType, error) {
	if s.buffer != nil {
		return nil, s.errorUnknown("incomplete buffer: %s", string(s.buffer))
	}
	if s.end() {
		if s.openParenth > 0 {
			return nil, s.errorSyntax("unexpected: unclosed '('")
		}
		return nil, nil
	}
	c := s.input[s.inputPos]
	s.move(1)
	switch c {
	case '\\':
		return lexBackslash, nil
	case '[':
		if s.openBracket {
			return nil, s.errorSyntax("nested '['")
		}
		s.openBracket = true
		return lexRange, nil
	case '{':
		if s.lastGen == nil {
			return nil, s.errorSyntax("nothing to repeat")
		}
		return lexRepeat, nil
	case '(':
		s.openParenth++
		s.lastGroupId++
		return lexGroup, nil
	case '$':
		return lexIdent, nil
	}
	err := s.addOutputOne(c)
	if err != nil {
		return nil, err
	}
	return LexRoot, nil
}

func backslashEscape(c rune) rune {
	switch c {
	case 't':
		return '\t'
	case 'r':
		return '\r'
	case 'n':
		return '\n'
	case 'v':
		return '\v'
	case 'f':
		return '\f'
	}
	return c
}

func lexBackslash(s *State) (LexType, error) {
	if s.end() {
		return LexRoot, nil
	}
	c := s.input[s.inputPos]
	if c >= '1' && c <= '9' {
		return processGroupRef(s, LexRoot)
	}
	s.move(1)
	switch c {
	case 'u':
		if s.buffer != nil {
			return nil, s.errorUnknown("incomplete buffer: %s", string(s.buffer))
		}
		return lexRootUnicode, nil
	case 'U':
		if s.buffer != nil {
			return nil, s.errorUnknown("incomplete buffer: %s", string(s.buffer))
		}
		return lexRootUnicodeWide, nil
	case 'd':
		return processRange(s, []rune("0123456789"))
	case 'w':
		return processRange(s, wordChars)
	}
	err := s.addOutputOne(backslashEscape(c))
	if err != nil {
		return nil, err
	}
	return LexRoot, nil
}

func lexIdent(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax(s_func_call_expected)
	}
	c := s.input[s.inputPos]
	s.move(1)
	switch c {
	case '\\', '[', '{', '$':
		return nil, s.errorSyntax(s_func_call_expected)
	case '(':
		s.openParenth++
		return lexIdentFuncCall, nil
	}
	s.buffer = append(s.buffer, c)
	return lexIdent, nil
}

func makeLexUnicode(parentLex LexType, symbol rune, width int, toBuffer bool) LexType {
	return func(s *State) (LexType, error) {
		buffer := make([]rune, 0, width)
		buffer = append(buffer, '\\', symbol)
		for ; len(buffer) < width && !s.end(); s.move(1) {
			buffer = append(buffer, s.input[s.inputPos])
		}
		if len(buffer) != width {
			s.errorMarkLen = len(buffer)
			return nil, s.errorSyntax("invalid escape sequence")
		}
		char, _, _, err := strconv.UnquoteChar(string(buffer), '"')
		if err != nil {
			s.errorMarkLen = width
			return nil, s.errorSyntax("invalid escape sequence")
		}
		if toBuffer {
			s.buffer = append(s.buffer, char)
		} else {
			err := s.addOutputOne(char)
			if err != nil {
				return nil, err
			}
		}
		return parentLex, nil
	}
}

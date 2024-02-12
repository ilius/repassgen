package passgen

import (
	"strconv"
)

func lexGroup(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.input[s.inputPos]
	s.move(1)
	switch c {
	case '\\':
		return processGroupBackslash(s, lexGroup)
	case '(':
		s.openParenth++
	case ')':
		s.openParenth--
		if s.openParenth > 0 {
			s.buff = append(s.buff, ')')
			return lexGroup, nil
		}
		return processGroupEnd(s)
	case '|':
		if s.end() {
			s.errorOffset++
			return nil, s.errorSyntax("'|' at the end of group")
		}
		s.moveBack(uint64(len(s.buff) + 1))
		return lexGroupAlter, nil
	}
	s.buff = append(s.buff, c)
	return lexGroup, nil
}

func lexGroupAlter(s *State) (LexType, error) {
	pattern := []rune{}
	length := uint64(0)
	openParenth := 1
Loop:
	for ; !s.end(); s.move(1) {
		length++
		c := s.input[s.inputPos]
		switch c {
		case '\\':
			s.move(1)
			if s.end() {
				// I don't know how to test this
				pattern = append(pattern, '\\')
				break
			}
			pattern = append(pattern, '\\', c)
		case '(':
			openParenth++
			pattern = append(pattern, '\\', c)
		case ')':
			openParenth--
			if openParenth > 0 {
				pattern = append(pattern, '\\', c)
				break
			}
			break Loop
		default:
			pattern = append(pattern, c)
		}
	}
	parts, indexList, err := splitArgsStr(pattern, '|')
	if err != nil {
		return nil, err
	}
	gen := &alterGenerator{
		parts:     parts,
		indexList: indexList,
		length:    length,
	}
	err = gen.Generate(s)
	if err != nil {
		return nil, err
	}
	s.move(1)
	s.openParenth--
	s.lastGen = gen
	s.buff = nil
	return LexRoot, nil
}

func processGroupBackslash(s *State, parentLex LexType) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.input[s.inputPos]
	s.move(1)
	s.buff = append(s.buff, '\\', c)
	return lexGroup, nil
}

func processGroupEnd(s *State) (LexType, error) {
	groupId := s.lastGroupId
	lastOutputSize := len(s.output)
	s2 := NewState(s.SharedState.Copy(), s.input)
	s2.output = s.output
	s2.errorOffset -= int64(len(s.buff) + 1)
	gen := newGroupGenerator(s.buff)
	err := gen.Generate(s2)
	if err != nil {
		return nil, err
	}
	s.output = s2.output
	s.patternEntropy = s2.patternEntropy
	s.lastGroupId = s2.lastGroupId
	s.groupsOutput[groupId] = s.output[lastOutputSize:]
	s.lastGen = gen
	s.buff = nil
	return LexRoot, nil
}

func processGroupRef(s *State, parentLex LexType) (LexType, error) {
	gid_r := []rune{}
	for ; !s.end(); s.move(1) {
		c := s.input[s.inputPos]
		if c < '0' || c > '9' {
			break
		}
		gid_r = append(gid_r, c)
	}
	gid, err := strconv.ParseInt(string(gid_r), 10, 64)
	if err != nil {
		return nil, s.errorUnknown("unexpected group id '%v'", string(gid_r))
	}
	output, ok := s.groupsOutput[uint64(gid)]
	if !ok {
		s.errorMarkLen = len(gid_r) + 1
		return nil, s.errorValue("invalid group id '%v'", gid)
	}

	err = s.addOutput(output)
	if err != nil {
		return nil, err
	}

	return parentLex, nil
}

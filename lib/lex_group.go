package passgen

import (
	"log"
	"strconv"
)

func lexGroup(s *State) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return processGroupBackslash(s, lexGroup)
	case '(':
		s.openParenth++
	case ')':
		s.openParenth--
		if s.openParenth > 0 {
			s.patternBuff = append(s.patternBuff, ')')
			return lexGroup, nil
		}
		return processGroupEnd(s)
	case '|':
		if s.end() {
			s.errorOffset++
			return nil, s.errorSyntax("'|' at the end of group")
		}
		s.moveBack(uint64(len(s.patternBuff) + 1))
		return lexGroupAlter, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexGroup, nil
}

func lexGroupAlter(s *State) (LexType, error) {
	pattern := []rune{}
	length := uint64(0)
	openParenth := 1
Loop:
	for ; !s.end(); s.move(1) {
		length++
		c := s.pattern[s.patternPos]
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
	// s.absPos = s.absPos - uint(length) - 1
	parts, indexList, err := splitArgsStr(pattern, '|')
	if err != nil {
		return nil, err
	}
	if length > s.absPos {
		// FIXME: this happens
		log.Printf(
			"pattern=`%v`, s.pattern=`%v`, length=%v, absPos=%v",
			string(pattern), string(s.pattern), length, s.absPos,
		)
	}
	gen := &alterGenerator{
		parts:     parts,
		indexList: indexList,
		absPos:    s.absPos - length,
	}
	err = gen.Generate(s)
	if err != nil {
		return nil, err
	}
	s.move(1)
	s.openParenth--
	s.lastGen = gen
	s.patternBuff = nil
	return LexRoot, nil
}

func processGroupBackslash(s *State, parentLex LexType) (LexType, error) {
	if s.end() {
		s.errorOffset++
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	s.patternBuff = append(s.patternBuff, '\\', c)
	return lexGroup, nil
}

func processGroupEnd(s *State) (LexType, error) {
	groupId := s.lastGroupId
	// lastOutputSize := len(s.output)
	s2 := NewState(NewSharedState(), s.pattern)
	s2.output = s.output
	if len(s.patternBuff) > int(s.absPos) {
		log.Printf(
			"patternBuff=%#v, len(patternBuff)=%v, absPos=%v",
			string(s.patternBuff), len(s.patternBuff), s.absPos,
		)
	}
	s2.absPos = s.absPos - uint64(len(s.patternBuff)) - 1
	s2.patternEntropy = s.patternEntropy
	s2.lastGroupId = groupId
	s2.groupsOutput = s.groupsOutput
	gen := newGroupGenerator(s.patternBuff)
	err := gen.Generate(s2)
	if err != nil {
		return nil, err
	}
	s.output = s2.output
	s.patternEntropy = s2.patternEntropy
	s.lastGroupId = s2.lastGroupId
	// s.groupsOutput[groupId] = s.output[lastOutputSize:]
	s.lastGen = gen
	s.patternBuff = nil
	return LexRoot, nil
}

func processGroupRef(s *State, parentLex LexType) (LexType, error) {
	gid_r := []rune{}
	for ; !s.end(); s.move(1) {
		c := s.pattern[s.patternPos]
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
	s.addOutput(output)
	return parentLex, nil
}

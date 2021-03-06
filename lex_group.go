package main

import (
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
		s.moveBack(uint(len(s.patternBuff) + 1))
		return lexGroupAlter, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexGroup, nil
}

func lexGroupAlter(s *State) (LexType, error) {
	pattern := []rune{}
	openParenth := 1
Loop:
	for ; !s.end(); s.move(1) {
		c := s.pattern[s.patternPos]
		switch c {
		case '\\':
			s.move(1)
			if s.end() {
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
	//s.absPos = s.absPos - uint(len(pattern)) - 1
	parts, indexList, err := splitArgsStr(pattern, '|')
	if err != nil {
		return nil, err
	}
	gen := &alterGenerator{
		parts:     parts,
		indexList: indexList,
		absPos:    s.absPos - uint(len(pattern)) - 1,
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
	lastOutputSize := len(s.output)
	s2 := NewState(&SharedState{}, s.pattern)
	s2.output = s.output
	s2.absPos = s.absPos - uint(len(s.patternBuff)) - 1
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
	s.groupsOutput[groupId] = s.output[lastOutputSize:]
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
	output, ok := s.groupsOutput[uint(gid)]
	if !ok {
		return nil, s.errorValue("invalid group id '%v'", gid)
	}
	s.addOutput(output)
	return parentLex, nil
}

package main

import "fmt"

func lexGroup(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '\\':
		return lexGroupBackslash, nil
	case '(':
		s.openParenth++
		s.startGroup()
	case '|':
		s.alterPos = append(s.alterPos, s.patternPos-1)
	case ')':
		s.openParenth--
		if s.openParenth > 0 {
			break
		}
		s.absPos -= uint(len(s.patternBuff)) + 1
		childPattern := s.patternBuff[s.patternBuffStart:]
		fmt.Printf("childPattern = %#v\n", string(childPattern))

		var gen generatorIface
		if len(s.alterPos) > 0 {
			gen = newAlterationGenerator(childPattern, s.alterPos)
		} else {
			gen = newGroupGenerator(childPattern)
		}

		err := gen.Generate(s)
		if err != nil {
			return nil, err
		}
		s.lastGen = gen
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexGroup, nil
}

func lexGroupBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.move(1)
	switch c {
	case '(', ')', '|':
		s.patternBuff = append(s.patternBuff, c)
	default:
		s.patternBuff = append(s.patternBuff, '\\', c)
	}
	return lexGroup, nil
}

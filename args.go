package main

import "fmt"

type ArgsParser struct {
	res         [][]rune
	indexList   []uint64
	openParenth int
	openBracket bool
	openCurly   bool
	backslash   bool
	sep         rune
}

func (p *ArgsParser) add(c rune) {
	i := len(p.res) - 1
	p.res[i] = append(p.res[i], c)
}

func (p *ArgsParser) checkBackslash(c rune) bool {
	if p.backslash {
		p.backslash = false
		if !(c == p.sep && p.openParenth == 0 && !p.openBracket && !p.openCurly) {
			p.add('\\')
		}
		p.add(c)
		return true
	}
	if c == '\\' {
		p.backslash = true
		return true
	}
	return false
}

func (p *ArgsParser) checkBrackets(c rune) bool {
	if p.openBracket {
		if c == ']' {
			p.openBracket = false
		}
		p.add(c)
		return true
	}
	if c == '[' {
		p.openBracket = true
		p.add(c)
		return true
	}
	return false
}

func (p *ArgsParser) checkSep(c rune, i int) bool {
	if c == p.sep {
		if p.openParenth == 0 && !p.openBracket && !p.openCurly {
			p.res = append(p.res, nil)
			p.indexList = append(p.indexList, uint64(i))
		} else {
			p.add(c)
		}
		return true
	}
	return false
}

func (p *ArgsParser) readChar(c rune, i int) error {
	if p.checkBackslash(c) {
		return nil
	}
	if p.checkBrackets(c) {
		return nil
	}
	if p.checkSep(c, i) {
		return nil
	}

	p.add(c)
	switch c {
	case '(':
		p.openParenth++
	case '{':
		if p.openCurly {
			return fmt.Errorf("nested '{'")
		}
		p.openCurly = true
	case ')':
		p.openParenth--
	case '}':
		p.openCurly = false
	}
	return nil
}

func (p *ArgsParser) parse(input []rune) error {
	for i, c := range input {
		err := p.readChar(c, i)
		if err != nil {
			return err
		}
	}
	if p.openParenth > 0 {
		return fmt.Errorf("unclosed '('")
	}
	if p.openParenth < 0 {
		return fmt.Errorf("too many ')'")
	}
	if p.openBracket {
		return fmt.Errorf("unclosed '['")
	}
	if p.openCurly {
		return fmt.Errorf("unclosed '{'")
	}
	if len(p.res[len(p.res)-1]) == 0 {
		p.res = p.res[:len(p.res)-1]
	}
	return nil
}

func (p *ArgsParser) Parse(input []rune) ([][]rune, []uint64, error) {
	err := p.parse(input)
	if err != nil {
		return nil, nil, err
	}
	return p.res, p.indexList, nil
}

func splitArgsStr(input []rune, sep rune) ([][]rune, []uint64, error) {
	p := &ArgsParser{
		res:       [][]rune{nil},
		indexList: []uint64{0},
		sep:       sep,
	}
	res, indexList, err := p.Parse(input)
	if err != nil {
		return nil, nil, err
	}
	return res, indexList, nil
}

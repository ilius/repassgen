package main

import "fmt"

func splitArgsStr(input []rune, sep rune) ([][]rune, error) {
	res := [][]rune{nil}
	openParenth := 0
	openBracket := false
	openCurly := false
	backslash := false

	add := func(c rune) {
		i := len(res) - 1
		res[i] = append(res[i], c)
	}

	for _, c := range input {
		if backslash {
			backslash = false
			if !(c == sep && openParenth == 0 && !openBracket && !openCurly) {
				add('\\')
			}
			add(c)
			continue
		}
		if c == '\\' {
			backslash = true
			continue
		}
		if openBracket {
			if c == ']' {
				openBracket = false
			}
			add(c)
			continue
		}
		if c == '[' {
			openBracket = true
			add(c)
			continue
		}
		if c == sep {
			if openParenth == 0 && !openBracket && !openCurly {
				res = append(res, nil)
			} else {
				add(c)
			}
			continue
		}
		add(c)
		switch c {
		case '(':
			openParenth++
		case '{':
			if openCurly {
				return nil, fmt.Errorf("nested '{'")
			}
			openCurly = true
		case ')':
			openParenth--
		case '}':
			openCurly = false
		}
	}
	if openParenth > 0 {
		return nil, fmt.Errorf("unclosed '('")
	}
	if openParenth < 0 {
		return nil, fmt.Errorf("too many ')'")
	}
	if openBracket {
		return nil, fmt.Errorf("unclosed '['")
	}
	if openCurly {
		return nil, fmt.Errorf("unclosed '{'")
	}
	if len(res[len(res)-1]) == 0 {
		res = res[:len(res)-1]
	}
	return res, nil
}

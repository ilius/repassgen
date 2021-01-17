package main

import "fmt"

func splitArgsStr(argsStr string) ([]string, error) {
	res := []string{""}
	openParenth := 0
	openBracket := false
	openCurly := false
	backslash := false
	for _, c := range argsStr {
		if backslash {
			backslash = false
			if !(c == ',' && openParenth == 0 && !openBracket && !openCurly) {
				res[len(res)-1] += "\\"
			}
			res[len(res)-1] += string(c)
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
			res[len(res)-1] += string(c)
			continue
		}
		if c == '[' {
			openBracket = true
			res[len(res)-1] += string(c)
			continue
		}
		if c == ',' {
			if openParenth == 0 && !openBracket && !openCurly {
				res = append(res, "")
			} else {
				res[len(res)-1] += string(c)
			}
			continue
		}
		res[len(res)-1] += string(c)
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

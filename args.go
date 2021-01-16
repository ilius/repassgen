package main

import "fmt"

func splitArgsStr(argsStr string) ([]string, error) {
	res := []string{""}
	openParenth := 0
	openBracket := 0
	openCurly := 0
	backslash := false
	for _, c := range argsStr {
		if backslash {
			backslash = false
			res[len(res)-1] += string(c)
			continue
		}
		if c == '\\' {
			backslash = true
			res[len(res)-1] += string(c)
			continue
		}
		if c == ',' {
			if openParenth == 0 && openBracket == 0 && openCurly == 0 {
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
		case '[':
			openBracket++
		case '{':
			openCurly++
		case ')':
			openParenth--
		case ']':
			openBracket--
		case '}':
			openCurly--
		}
	}
	if openParenth > 0 {
		return nil, fmt.Errorf("unclosed '('")
	}
	if openParenth < 0 {
		return nil, fmt.Errorf("too many ')'")
	}
	if openBracket > 0 {
		return nil, fmt.Errorf("unclosed '['")
	}
	if openBracket > 0 {
		return nil, fmt.Errorf("too many ']'")
	}
	if openCurly > 0 {
		return nil, fmt.Errorf("unclosed '{'")
	}
	if openCurly < 0 {
		return nil, fmt.Errorf("too many '{'")
	}
	if len(res[len(res)-1]) == 0 {
		res = res[:len(res)-1]
	}
	return res, nil
}

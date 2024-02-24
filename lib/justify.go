package passgen

import (
	"strconv"
	"strings"
)

func rjust(in []rune, width int, fillChar rune) []rune {
	if len(in) >= width {
		return in
	}
	out := make([]rune, width)
	fcc := width - len(in)
	for i := range fcc {
		out[i] = fillChar
	}
	copy(out[fcc:], in)
	return out
}

func ljust(in []rune, width int, fillChar rune) []rune {
	if len(in) >= width {
		return in
	}
	out := make([]rune, width)
	copy(out, in)
	for i := len(in); i < width; i++ {
		out[i] = fillChar
	}
	return out
}

func center(in []rune, width int, fillChar rune) []rune {
	if len(in) >= width {
		return in
	}
	out := make([]rune, width)
	fcc := (width - len(in)) / 2
	for i := range fcc {
		out[i] = fillChar
	}
	copy(out[fcc:], in)
	for i := len(in) + fcc; i < width; i++ {
		out[i] = fillChar
	}
	return out
}

type JustifyArgs struct {
	pattern  []rune
	width    int
	fillChar rune
}

type justifyGenerator struct {
	entropy     *float64
	args        *JustifyArgs
	justifyFunc func([]rune, int, rune) []rune
}

func (g *justifyGenerator) Generate(s *State) error {
	output, err := subGenerate(s, g.args.pattern)
	if err != nil {
		return err
	}
	output = g.justifyFunc(output, g.args.width, g.args.fillChar)
	s.output = append(s.output, output...)
	g.entropy = &s.patternEntropy
	return nil
}

func (g *justifyGenerator) Entropy(s *State) (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, s.errorUnknown("entropy is not calculated")
}

func parseJustifyArgs(s *State, argsStr []rune, funcName string) (*JustifyArgs, error) {
	if len(argsStr) < 3 {
		return nil, s.errorArg("%s: too few characters as arguments", funcName)
	}
	args, _, err := splitArgsStr(argsStr, ',')
	if err != nil {
		return nil, err
	}
	if len(args) < 2 {
		s.errorOffset += int64(len(argsStr) + 1)
		return nil, s.errorArg("%s: at least 2 arguments are required", funcName)
	}
	pattern := args[0]
	width, err := strconv.Atoi(strings.TrimSpace(string(args[1])))
	if err != nil {
		s.errorOffset += int64(len(args[0]) + len(args[1]) + 1)
		s.errorMarkLen = len(args[1])
		return nil, s.errorValue("invalid width %s", string(args[1]))
	}
	if width < 1 {
		s.errorOffset += int64(len(args[0]) + len(args[1]) + 1)
		s.errorMarkLen = len(args[1])
		return nil, s.errorValue("invalid width %s", string(args[1]))
	}
	fillChar := ' '
	if len(args) > 2 {
		fillCharA := args[2]
		if len(fillCharA) != 1 {
			s.errorOffset += int64(len(args[0]) + len(args[1]) + len(args[2]) + 2)
			s.errorMarkLen = len(args[2])
			return nil, s.errorValue("invalid fillChar=%#v, must have length 1", string(args[2]))
		}
		fillChar = fillCharA[0]
	}
	return &JustifyArgs{
		pattern:  pattern,
		width:    width,
		fillChar: fillChar,
	}, nil
}

func newRjustGenerator(s *State, argsStr []rune) (*justifyGenerator, error) {
	args, err := parseJustifyArgs(s, argsStr, "rjust")
	if err != nil {
		return nil, err
	}
	return &justifyGenerator{
		args:        args,
		justifyFunc: rjust,
	}, nil
}

func newLjustGenerator(s *State, argsStr []rune) (*justifyGenerator, error) {
	args, err := parseJustifyArgs(s, argsStr, "ljust")
	if err != nil {
		return nil, err
	}
	return &justifyGenerator{
		args:        args,
		justifyFunc: ljust,
	}, nil
}

func newCenterGenerator(s *State, argsStr []rune) (*justifyGenerator, error) {
	args, err := parseJustifyArgs(s, argsStr, "center")
	if err != nil {
		return nil, err
	}
	return &justifyGenerator{
		args:        args,
		justifyFunc: center,
	}, nil
}

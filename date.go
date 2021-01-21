package main

import (
	rand "crypto/rand"
	"math"
	"math/big"
	"strconv"
	"strings"

	lib "github.com/ilius/libgostarcal"
	"github.com/ilius/libgostarcal/cal_types/gregorian"
)

type dateGenerator struct {
	sep     string
	startJd int
	endJd   int
}

func (g *dateGenerator) Generate(s *State) error {
	startJd := g.startJd
	endJd := g.endJd
	randBig, err := rand.Int(rand.Reader, big.NewInt(int64(endJd-startJd)))
	if err != nil {
		panic(err)
	}
	jd := startJd + int(randBig.Int64())
	date := gregorian.JdTo(jd)
	dateStr := date.StringWithSep(g.sep)
	s.addOutputNonRepeatable([]rune(dateStr))
	entropy, err := g.Entropy(s)
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *dateGenerator) Entropy(s *State) (float64, error) {
	return math.Log2(float64(g.endJd - g.startJd)), nil
}

func newDateGenerator(s *State, argsStr string) (*dateGenerator, error) {
	if len(argsStr) < 3 {
		s.errorOffset += len(argsStr) + 1
		return nil, s.errorArg("date: too few characters as arguments")
	}
	args, err := splitArgsStr(argsStr)
	if err != nil {
		return nil, err
	}
	if len(args) < 2 {
		s.errorOffset += len(argsStr) + 1
		return nil, s.errorArg("date: at least 2 arguments are required")
	}
	startYear, err := strconv.Atoi(strings.TrimSpace(args[0]))
	if err != nil {
		s.errorOffset += len(args[0])
		return nil, s.errorValue("invalid year %s", args[0])
	}
	endYear, err := strconv.Atoi(strings.TrimSpace(args[1]))
	if err != nil {
		s.errorOffset += len(argsStr)
		return nil, s.errorValue("invalid year %s", args[1])
	}
	sep := "-"
	if len(args) > 2 {
		sep = args[2]
	}
	startJd := gregorian.ToJd(lib.NewDate(startYear, 1, 1))
	endJd := gregorian.ToJd(lib.NewDate(endYear, 1, 1))
	return &dateGenerator{
		startJd: startJd,
		endJd:   endJd,
		sep:     sep,
	}, nil
}

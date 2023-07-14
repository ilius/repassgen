package passgen

import (
	"fmt"
	"io"
	"log"
)

// GenerateInput is struct given to Generate
type GenerateInput struct {
	Pattern []rune
	Output  io.Writer
}

// GenerateOutput is struct returned by Generate
type GenerateOutput struct {
	PatternEntropy float64
}

// Generate generates random password based on given pattern
// see README.md for examples of pattern
func Generate(in GenerateInput) (*GenerateOutput, *State, error) {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("pattern=`%v`", in.Pattern)
		}
	}()
	if len(in.Pattern) > 1000 {
		return nil, nil, fmt.Errorf("pattern is too long")
	}
	ss := NewSharedState()
	s := NewState(ss, in.Pattern)
	g := NewRootGenerator()
	{
		err := g.Generate(s)
		if err != nil {
			return nil, s, err
		}
	}
	return &GenerateOutput{
		PatternEntropy: s.patternEntropy,
	}, s, nil
}

func subGenerate(s *State, pattern []rune) error {
	childGen := NewRootGenerator()
	ss := s.SharedState
	{
		s := NewState(ss, pattern)
		err := childGen.Generate(s)
		if err != nil {
			return err
		}
	}
	s.lastGen = nil
	return nil
}

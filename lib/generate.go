package passgen

import (
	"fmt"
	"log"
)

// GenerateInput is struct given to Generate
type GenerateInput struct {
	Pattern []rune
}

// GenerateOutput is struct returned by Generate
type GenerateOutput struct {
	Password       []rune
	PatternEntropy float64
}

// Generate generates random password based on given pattern
// see README.md for examples of pattern
func Generate(in GenerateInput) (*GenerateOutput, *State, error) {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("panic in Generate: %v, pattern=`%v`", r, in.Pattern)
		}
	}()
	if len(in.Pattern) > 1000 {
		return nil, nil, fmt.Errorf("pattern is too long")
	}
	s := NewState(NewSharedState(), in.Pattern)
	g := NewRootGenerator()

	err := g.Generate(s)
	if err != nil {
		return nil, s, err
	}

	return &GenerateOutput{
		Password:       s.output,
		PatternEntropy: s.patternEntropy,
	}, s, nil
}

func subGenerate(s *State, pattern []rune) ([]rune, error) {
	childGen := NewRootGenerator()
	s2 := NewState(s.SharedState, pattern)
	err := childGen.Generate(s2)
	if err != nil {
		return nil, err
	}
	s.lastGen = nil
	return s2.output, nil
}

package main

import "fmt"

func shuffle(in []rune) ([]rune, error) {
	r := NewRandSource()
	r.Shuffle(len(in), func(i int, j int) {
		in[i], in[j] = in[j], in[i]
	})
	return in, nil
}

type shuffleGenerator struct {
	argPattern        string
	argPatternEntropy *float64
	entropy           *float64
}

func (g *shuffleGenerator) Generate(s *State) error {
	argOut, err := baseFunctionCallGenerator(s, "shuffle", shuffle, g.argPattern)
	if err != nil {
		return err
	}
	g.argPatternEntropy = &argOut.PatternEntropy
	entropy, err := g.Entropy()
	if err != nil {
		return err
	}
	s.patternEntropy += entropy
	return nil
}

func (g *shuffleGenerator) Level() int {
	return 0
}

func (g *shuffleGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		fmt.Println("Using prev calculated entropy:", *g.entropy)
		return *g.entropy, nil
	}
	if g.argPatternEntropy == nil {
		return 0, fmt.Errorf("argument entropy is not calculated")
	}
	// FIXME: how to calculate entropy?
	// if char classes don't have any char in common, the answer would be
	// to calculate the number of distinct patterns that may result
	// from shuffling the pattern
	// then multiply this number with pattern entropy and return the result
	// but in gerenal...
	// what if we calculate of maximum probability of a password being chosen
	// entropy = 1 / max(probability)
	// but first, we need a way to extract char count map which is not very easy
	// for example, imagine this pattern: $shuffle($base64([:byte:]{8}))
	// pattern := []rune(g.argPattern)
	// charCountM := MakeRuneCountMap(pattern)
	// charCount := SortRuneCountMap(charCountM)
	// fmt.Println(charCount)
	entropy := *g.argPatternEntropy
	g.entropy = &entropy
	fmt.Println("Calculated entropy:", entropy)
	return entropy, nil
}

func newShuffleGenerator(arg string) (*shuffleGenerator, error) {
	return &shuffleGenerator{
		argPattern: arg,
	}, nil
}

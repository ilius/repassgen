package passgen

type repeatGenerator struct {
	child generatorIface
	count int64
}

func (g *repeatGenerator) Generate(s *State) error {
	child := g.child
	count := g.count
	for range count {
		err := child.Generate(s)
		if s.tooLong() {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *repeatGenerator) Entropy(s *State) (float64, error) {
	childEntropy, err := g.child.Entropy(s)
	if err != nil {
		return 0, err
	}
	return childEntropy * float64(g.count), nil
}

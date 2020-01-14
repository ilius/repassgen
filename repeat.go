package main

type repeatGenerator struct {
	child generatorIface
	count int
}

func (g *repeatGenerator) Generate(s *State) error {
	child := g.child
	count := g.count
	for i := 0; i < count; i++ {
		err := child.Generate(s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *repeatGenerator) Entropy() (float64, error) {
	childEntropy, err := g.child.Entropy()
	if err != nil {
		return 0, err
	}
	return childEntropy * float64(g.count), nil
}

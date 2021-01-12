package main

type generatorIface interface {
	Generate(s *State) error
	Entropy(s *State) (float64, error)
}

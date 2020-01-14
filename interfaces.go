package main

type generatorIface interface {
	Generate(s *State) error
	Entropy() (float64, error)
}

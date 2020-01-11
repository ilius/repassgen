package main

type generatorIface interface {
	Generate(s *State) error
	Level() int
	Entropy() (float64, error)
}

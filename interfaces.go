package main

type generatorIface interface {
	Generate(s *State) error
	Level() int
	Entropy(s *State) (float64, error)
	CharProb() map[rune]float64
}

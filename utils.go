package main

import (
	"fmt"
	"sort"
)

func MakeRuneCountMap(in []rune) map[rune]uint {
	m := map[rune]uint{}
	for _, r := range in {
		m[r]++
	}
	return m
}

type RuneCount struct {
	Rune  rune
	Count uint
}

func (rc RuneCount) String() string {
	return fmt.Sprintf("%s(%d)", string(rc.Rune), rc.Count)
}

func SortRuneCountMap(m map[rune]uint) []RuneCount {
	ls := make([]RuneCount, 0, len(m))
	for r, n := range m {
		ls = append(ls, RuneCount{Rune: r, Count: n})
	}
	sort.Slice(ls, func(i int, j int) bool {
		return ls[i].Count > ls[j].Count
	})
	return ls
}

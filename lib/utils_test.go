package passgen

import (
	"math/rand"
	"testing"
)

func removeDuplicateRunes_2(list []rune) []rune {
	set := map[rune]struct{}{}
	newlist := make([]rune, 0, len(list))
	for _, c := range list {
		if _, ok := set[c]; ok {
			continue
		}
		set[c] = struct{}{}
		newlist = append(newlist, c)
	}
	return newlist
}

func excludeCharsASCII_2(exclude []rune) []rune {
	ex_set := map[rune]struct{}{}
	for _, c := range exclude {
		ex_set[c] = struct{}{}
	}
	list := make([]rune, 0, 95)
	for ci := 32; ci < 127; ci++ {
		c := rune(ci)
		if _, ok := ex_set[c]; ok {
			continue
		}
		list = append(list, c)
	}
	return list
}

func Benchmark_removeDuplicateRunes(b *testing.B) {
	listLength := 100
	count := b.N * 10000
	lists := make([][]rune, count)
	for i := 0; i < count; i++ {
		list := make([]rune, listLength)
		for j := 0; j < listLength; j++ {
			list[j] = rune(rand.Intn(256))
		}
		lists[i] = list
	}
	b.Run("map[rune]bool", func(b *testing.B) {
		for _, list := range lists {
			removeDuplicateRunes(list)
		}
	})
	b.Run("map[rune]struct{}", func(b *testing.B) {
		for _, list := range lists {
			removeDuplicateRunes_2(list)
		}
	})
}

func Benchmark_excludeCharsASCII(b *testing.B) {
	listLength := 100
	count := b.N * 10000
	lists := make([][]rune, count)
	for i := 0; i < count; i++ {
		list := make([]rune, listLength)
		for j := 0; j < listLength; j++ {
			list[j] = rune(rand.Intn(256))
		}
		lists[i] = list
	}
	b.Run("map[rune]bool", func(b *testing.B) {
		for _, list := range lists {
			excludeCharsASCII(list)
		}
	})
	b.Run("map[rune]struct{}", func(b *testing.B) {
		for _, list := range lists {
			excludeCharsASCII_2(list)
		}
	})
}

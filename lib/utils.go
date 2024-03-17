package passgen

func removeDuplicateRunes(chars []rune) []rune {
	set := map[rune]bool{}
	newChars := make([]rune, 0, len(chars))
	for _, c := range chars {
		if set[c] {
			continue
		}
		set[c] = true
		newChars = append(newChars, c)
	}
	return newChars
}

func byteRange(start uint, end uint) []rune {
	a := make([]rune, 0, end-start+1)
	for x := start; x <= end; x++ {
		a = append(a, rune(x))
	}
	return a
}

func excludeCharsASCII(exclude []rune) []rune {
	ex_set := map[rune]bool{}
	for _, c := range exclude {
		ex_set[c] = true
	}
	list := make([]rune, 0, 95)
	for ci := 32; ci < 127; ci++ {
		c := rune(ci)
		if ex_set[c] {
			continue
		}
		list = append(list, c)
	}
	return list
}

func hasRune(st []rune, c rune) bool {
	for _, c2 := range st {
		if c2 == c {
			return true
		}
	}
	return false
}

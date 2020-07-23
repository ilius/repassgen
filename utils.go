package main

func removeDuplicateRunes(list []rune) []rune {
	set := map[rune]bool{}
	newlist := make([]rune, 0, len(list))
	for _, c := range list {
		if set[c] {
			continue
		}
		set[c] = true
		newlist = append(newlist, c)
	}
	return newlist
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

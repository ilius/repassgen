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

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

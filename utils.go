package main

/*func removeDuplicateRunes(list []rune) []rune {
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
}*/

func removeDuplicateBytes(list []byte) []byte {
	set := map[byte]bool{}
	newlist := make([]byte, 0, len(list))
	for _, c := range list {
		if set[c] {
			continue
		}
		set[c] = true
		newlist = append(newlist, c)
	}
	return newlist
}

func byteRange(start byte, end byte) []byte {
	a := make([]byte, 0, end-start+1)
	for x := start; x <= end; x++ {
		a = append(a, byte(x))
	}
	return a
}

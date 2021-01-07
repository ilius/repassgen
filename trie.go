/*
Copyright 2013 Herman Schaaf and Shawn Smith

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

// Trie is a trie data structure
type Trie struct {
	children map[string]*Trie
	letter   string
	values   []string
}

// Build a trie for efficient retrieval of entries
func newTrie() *Trie {
	return &Trie{map[string]*Trie{}, "", []string{}}
}

// Insert a value into the trie
func (t *Trie) insert(letters, value string) {
	lettersRune := []rune(letters)

	// loop through letters in argument word
	for l, letter := range lettersRune {

		letterStr := string(letter)

		// if letter in children
		if t.children[letterStr] != nil {
			t = t.children[letterStr]
		} else {
			// not found, so add letter to children
			t.children[letterStr] = &Trie{map[string]*Trie{}, "", []string{}}
			t = t.children[letterStr]
		}

		if l == len(lettersRune)-1 {
			// last letter, save value and exit
			t.values = append(t.values, value)
			break
		}
	}
}

// Convert a given string to the corresponding values
// in the trie. This performed in a greedy fashion,
// replacing the longest valid string it can find at any
// given point.
func (t *Trie) convert(origin string) (result string) {
	root := t
	originRune := []rune(origin)
	result = ""

	for l := 0; l < len(originRune); l++ {
		t = root
		foundVal := ""
		depth := 0
		for i := 0; i+l < len(originRune); i++ {
			letter := string(originRune[l+i])
			if t.children[letter] == nil {
				// not found
				break
			}
			if len(t.children[letter].values) > 0 {
				foundVal = t.children[letter].values[0]
				depth = i
			}
			t = t.children[letter]
		}
		if foundVal != "" {
			result += foundVal
			l += depth
		} else {
			result += string(originRune[l : l+1])
		}
	}
	return result
}

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

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var (
	consonants = []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "p", "r", "s", "t", "w", "z"}

	hiraganaRe = regexp.MustCompile(`ん([あいうえおなにぬねの])`)
	katakanaRe = regexp.MustCompile(`ン([アイウエオナニヌネノ])`)

	kanaToRomajiTrie *Trie
)

// Initialize builds the Hiragana + Katakana trie.
// Because there is no overlap between the hiragana and katakana sets,
// they both use the same trie without conflict. Nice bonus!
func InitRomaji() {
	kanaToRomajiTrie = newTrie()

	tables := []string{HiraganaTable, KatakanaTable}
	for _, table := range tables {
		rows := strings.Split(table, "\n")
		colNames := strings.Split(string(rows[0]), "\t")[1:]
		for _, row := range rows[1:] {
			cols := strings.Split(string(row), "\t")
			rowName := cols[0]
			for i, kana := range cols[1:] {
				value := rowName + colNames[i]
				kanas := strings.Split(kana, "/")
				for _, singleKana := range kanas {
					if singleKana != "" {
						// add to tries
						kanaToRomajiTrie.insert(singleKana, value)
					}
				}
			}
		}
	}
}

// KanaToRomaji converts a kana string to its romaji form
func KanaToRomaji(kana string) (romaji string) {
	// unfortunate hack to deal with double n's
	romaji = hiraganaRe.ReplaceAllString(kana, "nn$1")
	romaji = katakanaRe.ReplaceAllString(romaji, "nn$1")

	romaji = kanaToRomajiTrie.convert(romaji)

	// do some post-processing for the tsu and stripe characters
	// maybe a bit of a hacky solution - how can we improve?
	// (they act more like punctuation)
	tsus := []string{"っ", "ッ"}
	for _, tsu := range tsus {
		if strings.Index(romaji, tsu) > -1 {
			for _, c := range romaji {
				ch := string(c)
				if ch == tsu {
					i := strings.Index(romaji, ch)
					runeSize := len(ch)
					followingLetter, _ := utf8.DecodeRuneInString(romaji[i+runeSize:])
					followingLetterStr := string(followingLetter)
					if followingLetterStr != tsu {
						romaji = strings.Replace(romaji, tsu, followingLetterStr, 1)
					} else {
						romaji = strings.Replace(romaji, tsu, "", 1)
					}
				}
			}
		}
	}

	line := "ー"
	for i := strings.Index(romaji, line); i > -1; i = strings.Index(romaji, line) {
		if i > 0 {
			romaji = strings.Replace(romaji, line, "-", 1)
		} else {
			romaji = strings.Replace(romaji, line, "", 1)
		}
	}
	return romaji
}

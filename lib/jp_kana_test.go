// From https://github.com/gojp/kana

package passgen

import "testing"

type kanaTest struct {
	orig, want string
}

var hiraganaToRomajiTests = []kanaTest{
	{"ああいうえお", "aaiueo"},
	{"かんじ", "kanji"},
	{"ちゃう", "chau"},
	{"きょうじゅ", "kyouju"},
	{"な\nに	ぬ	ね	の", "na\nni	nu	ne	no"},
	{"ばか dog", "baka dog"},
	{"きった", "kitta"},
	{"はんのう", "hannnou"},
	{"ぜんいん", "zennin"},
	{"んい", "nni"},
	{"はんのう", "hannnou"},
	{"はんおう", "hannou"},
	{"あうでぃ", "audexi"},
	{"だぢづを", "dajiduwo"},
	{"ぢゃぢょぢゅ", "jazoju"},
}

func TestHiraganaToRomaji(t *testing.T) {
	for _, tt := range hiraganaToRomajiTests {
		if got := KanaToRomaji(tt.orig); got != tt.want {
			t.Errorf("KanaToRomaji(%q) = %q, want %q", tt.orig, got, tt.want)
		}
	}
}

var katakanaToRomajiTests = []kanaTest{
	{"バナナ", "banana"},
	{"カンジ", "kanji"},
	{"テレビ", "terebi"},
	{"baking バナナ pancakes", "baking banana pancakes"},
	{"ベッド", "beddo"},
	{"モーター", "mo-ta-"},
	{"ＣＤプレーヤー", "ＣＤpure-ya-"},
	{"オーバーヘッドキック", "o-ba-heddokikku"},
	{"ハンノウ", "hannnou"},
	{"アウディ", "audexi"},
	{"ダヂヅヲ", "dajiduwo"},
	{"ヂャヂョヂュ", "jazoju"},
}

func TestKatakanaToRomaji(t *testing.T) {
	for _, tt := range katakanaToRomajiTests {
		if got := KanaToRomaji(tt.orig); got != tt.want {
			t.Errorf("KanaToRomaji(%q) = %q, want %q", tt.orig, got, tt.want)
		}
	}
}

package main

var charsets = map[string][]rune{
	// POSIX character classes, https://www.regular-expressions.info/posixbrackets.html
	"alpha":  []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	"alnum":  []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
	"lower":  []rune("abcdefghijklmnopqrstuvwxyz"),
	"upper":  []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	"digit":  []rune("0123456789"),
	"xdigit": []rune("0123456789abcdefABCDEF"),
	"punct":  []rune("[!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]"),

	// Crockford's Base32
	"b32": []rune("0123456789abcdefghjkmnpqrstvwxyz"),
	"B32": []rune("0123456789ABCDEFGHJKMNPQRSTVWXYZ"),

	// standard Base64
	"b64": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"),

	// URL-compatible Base64
	"b64url": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"),

	// Visible characters and spaces (anything except control characters)
	"print": byteRange(0x20, 0x7E),

	// All whitespace characters, including line breaks
	"space": []rune(" \t\r\n\v\f"),

	// Space and tab
	"blank": []rune(" \t"),

	// Control characters
	"cntrl": append(byteRange(0x00, 0x1F), '\x7F'),

	"graph": byteRange(0x21, 0x7E),
}

func byteRange(start byte, end byte) []rune {
	a := make([]rune, 0, end-start+1)
	for x := start; x <= end; x++ {
		a = append(a, rune(x))
	}
	return a
}

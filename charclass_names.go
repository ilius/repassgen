package main

var charClasses = map[string][]rune{
	// POSIX character classes, https://www.regular-expressions.info/posixbrackets.html
	"alpha": []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"),
	"alnum": []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),

	// Word characters (letters, numbers and underscores), [A-Za-z0-9_]
	"word": []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"),

	"lower": []rune("abcdefghijklmnopqrstuvwxyz"),
	"upper": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"),

	"digit":  []rune("0123456789"),
	"xdigit": []rune("0123456789abcdefABCDEF"),

	"punct": []rune("[!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]"),

	// Crockford's Base32
	"b32": []rune("0123456789abcdefghjkmnpqrstvwxyz"),
	"B32": []rune("0123456789ABCDEFGHJKMNPQRSTVWXYZ"),

	// Standard Base32
	"B32STD": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"),

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

	// Visible characters (anything except spaces and control characters)
	"graph": byteRange(0x21, 0x7E),

	// ASCII characters
	"ascii": byteRange(0x00, 0x7F),

	// Any byte
	"byte": byteRange(0x00, 0xFF),
}

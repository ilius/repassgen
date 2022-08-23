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

	"punct": []rune("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"),

	// Crockford's Base32
	"b32": []rune("0123456789abcdefghjkmnpqrstvwxyz"),
	"B32": []rune("0123456789ABCDEFGHJKMNPQRSTVWXYZ"),

	// Standard Base32
	"B32STD": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"),

	// standard Base64
	"b64": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"),

	// URL-compatible Base64
	"b64url": []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"),

	// Printable/visible characters (anything except control characters)
	// ‘[:alnum:]’, ‘[:punct:]’, and space.
	"print": byteRange(0x20, 0x7E),

	// Graphical characters: ‘[:alnum:]’ and ‘[:punct:]’.
	"graph": byteRange(0x21, 0x7E),

	// All whitespace characters, including line breaks
	"space": []rune(" \t\r\n\v\f"),

	// Space and tab
	"blank": []rune(" \t"),

	// Control characters
	// GNU: Control characters. In ASCII, these characters have octal codes 000
	// through 037, and 177 (DEL). In other character sets, these are the
	// equivalent characters, if any.
	"cntrl": append(byteRange(0x00, 0x1F), '\x7F'),

	// ASCII characters
	"ascii": byteRange(0x00, 0x7F),

	// Any byte
	"byte": byteRange(0x01, 0xFF),
}

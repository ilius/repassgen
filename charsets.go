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
}

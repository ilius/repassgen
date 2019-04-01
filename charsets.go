package main

var charsets = map[string]string{
	// POSIX character classes, https://www.regular-expressions.info/posixbrackets.html
	"alpha":  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"alnum":  "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789",
	"lower":  "abcdefghijklmnopqrstuvwxyz",
	"upper":  "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
	"digit":  "0123456789",
	"xdigit": "0123456789abcdefABCDEF",
	"punct":  "[!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~]",

	// Crockford's Base32
	"b32": "0123456789abcdefghjkmnpqrstvwxyz",
	"B32": "0123456789ABCDEFGHJKMNPQRSTVWXYZ",

	// standard Base64
	"b64": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/",

	// URL-compatible Base64
	"b64url": "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_",
}

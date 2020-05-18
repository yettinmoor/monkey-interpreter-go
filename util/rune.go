package util

import "unicode/utf8"

func IsAlpha(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}

func IsNum(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsAlphaNum(r rune) bool {
	return IsAlpha(r) || IsNum(r)
}

func StringIsAlphaNum(s string) bool {
	first, w := utf8.DecodeRuneInString(s)
	if !IsAlpha(first) {
		return false
	}
	for _, r := range s[w:] {
		if !(IsAlpha(r) || IsNum(r)) {
			return false
		}
	}
	return true
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

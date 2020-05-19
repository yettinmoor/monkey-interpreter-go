package lexer

import (
	"unicode"
	"unicode/utf8"
)

func IsValidIdentifierHead(r rune) bool {
	return IsAlpha(r) || unicode.IsLetter(r)
}

func IsValidIdentifierRune(r rune) bool {
	return IsNum(r) || IsValidIdentifierHead(r)
}

func IsValidIdentifier(s string) bool {
	first, w := utf8.DecodeRuneInString(s)
	if !IsValidIdentifierHead(first) {
		return false
	}
	for _, r := range s[w:] {
		println(r)
		if !IsValidIdentifierRune(r) {
			return false
		}
	}
	return true
}

func IsAlpha(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_'
}

func IsNum(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsAlphaNum(r rune) bool {
	return IsAlpha(r) || IsNum(r)
}

func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r' || unicode.IsSpace(r)
}

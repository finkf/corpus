package corpus

import (
	"unicode"
)

// Tokenizer defines the interface for tokenizers.
// The given callback function should be called for every
// non-empty token. Tokens for which the given callback is called
// must not contain any whitespace.
type Tokenizer interface {
	Tokenize(func(string)) error
}

// TokenType represents the type of a token
type TokenType int

// The different types of tokens.
const (
	ucletter TokenType = 1 << iota
	lcletter
	digit
	punctuation
)

// Letter returns true if the token type denotes letters.
func (t TokenType) Letter() bool {
	return t.UpperCaseLetter() || t.LowerCaseLetter()
}

// UpperCaseLetter returns true if the token type denotes uppercase letters.
func (t TokenType) UpperCaseLetter() bool {
	return (t & ucletter) != 0
}

// LowerCaseLetter returns true if the token type denotes uppercase letters.
func (t TokenType) LowerCaseLetter() bool {
	return (t & lcletter) != 0
}

// Digit returns true if the token type denotes digits.
func (t TokenType) Digit() bool {
	return (t & digit) != 0
}

// Punctuation returns true if the token type denotes
// punctuation characters.
func (t TokenType) Punctuation() bool {
	return (t & punctuation) != 0
}

// TokenTypeOf returns the token type of a given string.
func TokenTypeOf(str string) TokenType {
	var t TokenType
	for _, r := range str {
		t |= runeFlagType(r)
	}
	return t
}

func runeFlagType(r rune) TokenType {
	if unicode.In(r, unicode.Digit) {
		return digit
	}
	if unicode.In(r, unicode.Lower) {
		return lcletter
	}
	if unicode.In(r, unicode.Title, unicode.Upper) {
		return ucletter
	}
	// Marks, nonspacing are considered
	// neither upper nor lower case letter but also not punctuation
	if unicode.In(r, unicode.Mn) {
		return 0
	}
	return punctuation
}

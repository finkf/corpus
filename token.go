package corpus

import "unicode"

// Tokener defines the interface for things that read
// a stream of tokens. If an error occurs, Err returns a non-nil value.
type Tokener interface {
	Tokens(func(Token)) error
}

// TokenType represents the type of a token
type TokenType int

// Different token types
const (
	Empty TokenType = 1 << iota
	Word
	Number
	Punctuation
	Mixed
)

// Token represents a token.
type Token string

// Type returns the TokenType of the token.
func (t Token) Type() TokenType {
	if len(t) == 0 {
		return Empty
	}
	typ := Empty
	for _, r := range t {
		tmp := runeType(r)
		if typ != Empty && typ != tmp {
			return Mixed
		}
		typ = tmp
	}
	return typ
}

// runeType returns the type of a rune.
// Everything that is not a word, or number is considered
// to be punctuation.
func runeType(r rune) TokenType {
	if IsLetter(r) {
		return Word
	}
	if unicode.IsNumber(r) {
		return Number
	}
	return Punctuation
}

// IsLetter returns true if the given
// rune is to be considered a letter.
// Any mark (unicode.IsMark(r) == true) and any letter
// (unicode.IsLetter(r))return true.
func IsLetter(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsMark(r)
}

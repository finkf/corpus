package corpus_test

import (
	"testing"

	"github.com/finkf/corpus"
)

func TestTokenTypeOf(t *testing.T) {
	tests := []struct {
		token                              string
		upper, lower, letter, digit, punct bool
	}{
		{"simple", false, true, true, false, false},
		{"Simple", true, true, true, false, false},
		{"0815", false, false, false, true, false},
		{"Simple-token", true, true, true, false, true},
		{"SIMPLE", true, false, true, false, false},
		{"ᾙallo", true, true, true, false, false},
		{"Waſſe̅r", true, true, true, false, false},
	}
	for _, tc := range tests {
		t.Run(tc.token, func(t *testing.T) {
			typ := corpus.TokenTypeOf(tc.token)
			checkType(t, "upper", tc.upper, typ.UpperCaseLetter())
			checkType(t, "lower", tc.lower, typ.LowerCaseLetter())
			checkType(t, "letter", tc.letter, typ.Letter())
			checkType(t, "digit", tc.digit, typ.Digit())
			checkType(t, "punctuation", tc.punct, typ.Punctuation())
		})
	}
}

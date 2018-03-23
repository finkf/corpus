package corpus

import "testing"

func TestToken(t *testing.T) {
	tests := []struct {
		token string
		typ   TokenType
	}{
		{"", Empty},
		{"word", Word},
		{"123", Number},
		{",", Punctuation},
		{"mixed-word", Mixed},
		{"Caſparis", Word},
	}
	for _, tc := range tests {
		t.Run(tc.token, func(t *testing.T) {
			tt := Token{tc.token}
			if got := tt.Type(); got != tc.typ {
				t.Fatalf("expected %d; got %d", got, tc.typ)
			}
			if got := tt.String(); got != tc.token {
				t.Fatalf("expeceted %q; got %q", tc.token, got)
			}
		})
	}
}

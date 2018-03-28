package corpus

import (
	"unicode"
)

type state rune

const (
	initial = 0
	rest    = 'r'
	digit   = 'd'
	alpha   = 'a'
)

type splitter struct {
	state  state
	tokens []string
}

func (s *splitter) split(str string) []string {
	var b int
	for i, r := range str {
		s.update(&b, i, r, str)
	}
	s.stop(b, str)
	return s.tokens
}

func (s *splitter) update(b *int, i int, r rune, str string) {
	next := getState(r)
	if next == s.state {
		return
	}
	if s.state == initial {
		s.state = next
		*b = 0
		return
	}
	// next != s.state
	s.tokens = append(s.tokens, str[*b:i])
	s.state = next
	*b = i
}

func (s *splitter) stop(b int, str string) {
	if s.state == initial {
		return
	}
	s.tokens = append(s.tokens, str[b:])
}

func getState(r rune) state {
	if unicode.IsLetter(r) {
		return alpha
	}
	if unicode.IsNumber(r) {
		return digit
	}
	return rest
}

package corpus

type state rune

const (
	stateInitial = 0
	stateOther   = 'r'
	stateDigit   = 'd'
	stateAlpha   = 'a'
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
	next := getState(s.state, r)
	if next == s.state {
		return
	}
	if s.state == stateInitial {
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
	if s.state == stateInitial {
		return
	}
	s.tokens = append(s.tokens, str[b:])
}

func getState(s state, r rune) state {
	switch runeFlagType(r) {
	case digit:
		return stateDigit
	case lcletter, ucletter:
		return stateAlpha
	case punctuation:
		return stateOther
	default:
		return s
	}
}

// Split splits a given string into an array of tokens
func Split(str string) []string {
	return new(splitter).split(str)
}

// tokenize splits the given string and calls the
// given callback for each token.
func tokenize(str string, f func(string)) {
	for _, token := range Split(str) {
		f(token)
	}
}

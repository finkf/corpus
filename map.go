package corpus

// Char3Grams is a map character triple counts
type Char3Grams struct {
	n uint64
	m map[string]uint64
}

// NewChar3Grams create a new Char3Grams instance.
func NewChar3Grams() *Char3Grams {
	return &Char3Grams{
		m: make(map[string]uint64),
	}
}

// AddAll adds all character 3-grams of the
// supplied string into the map.
func (m *Char3Grams) AddAll(str string) {
	EachChar3Gram(str, func(str string) {
		m.m[str]++
		m.n++
	})
}

// Add3Grams adds the 3-grams of anohter map to this.
func (m *Char3Grams) Add3Grams(o *Char3Grams) {
	for k, v := range o.m {
		m.m[k] += v
	}
	m.n += o.n
}

// Get returns the number of the supplied 3-gram.
func (m *Char3Grams) Get(str string) uint64 {
	c, ok := m.m[str]
	if !ok {
		return 0
	}
	return c
}

// Total returns the total number of 3-grams in the map.
func (m *Char3Grams) Total() uint64 {
	return m.n
}

// Len return the number of different 3-grams in the map.
func (m *Char3Grams) Len() uint64 {
	return uint64(len(m.m))
}

// // Add adds the triples of another map into a map.
// func (m CharTripleMap) Add(o CharTripleMap) CharTripleMap {
// 	for k, v := range o {
// 		m[k] += v
// 	}
// 	return m
// }

// EachChar3Gram iterates of all character 3-grams in the given string.
// It calls the supplied callback function for each such 3-gram.
func EachChar3Gram(str string, f func(string)) {
	pos := make([]int, 0, len(str)+1)
	for i := range str {
		pos = append(pos, i)
		add3Gram(pos, str, f)
	}
	pos = append(pos, len(str))
	add3Gram(pos, str, f)
}

func add3Gram(pos []int, str string, f func(string)) {
	if len(pos) < 4 {
		return
	}
	s := pos[len(pos)-4]
	e := pos[len(pos)-1]
	f(str[s:e])
}

// ReadCharTripleMap reads a new CharTripleMap from a token stream.
// The second argument is the type of tokens that should be used.
// To allow multiple types of tokens use | to combine them:
// ReadCharTripleMap(t, Word | Number)
// func ReadCharTripleMap(t Tokener, tts TokenType) (CharTripleMap, error) {
// 	m := make(CharTripleMap)
// 	for token := range t.Tokens() {
// 		if (tts & token.Type()) == 0 {
// 			continue
// 		}
// 		var rs []rune
// 		for _, r := range token.Str {
// 			rs = shiftRune(rs, r)
// 			if len(rs) == 3 {
// 				m[CharTriple{rs[0], rs[1], rs[2]}]++
// 			}
// 		}
// 	}
// 	return m, t.Err()
// }

// func shiftRune(rs []rune, r rune) []rune {
// 	if len(rs) < 3 {
// 		return append(rs, r)
// 	}
// 	rs[0] = rs[1]
// 	rs[1] = rs[2]
// 	rs[2] = r
// 	return rs
// }

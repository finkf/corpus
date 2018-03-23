package ocrddta

// CharTriple represents a sequece of
// three characters.
type CharTriple [3]rune

// CharTripleMap is a map character triple counts
type CharTripleMap map[CharTriple]uint64

// Add adds the triples of another map into a map.
func (m CharTripleMap) Add(o CharTripleMap) CharTripleMap {
	for k, v := range o {
		m[k] += v
	}
	return m
}

// ReadCharTripleMap reads a new CharTripleMap from a token stream.
// The second argument is the type of tokens that should be used.
// To allow multiple types of tokens use | to combine them:
// ReadCharTripleMap(t, Word | Number)
func ReadCharTripleMap(t Tokener, tts TokenType) (CharTripleMap, error) {
	m := make(CharTripleMap)
	for token := range t.Tokens() {
		if (tts & token.Type()) == 0 {
			continue
		}
		var rs []rune
		for _, r := range token.Str {
			rs = shiftRune(rs, r)
			if len(rs) == 3 {
				m[CharTriple{rs[0], rs[1], rs[2]}]++
			}
		}
	}
	return m, t.Err()
}

func shiftRune(rs []rune, r rune) []rune {
	if len(rs) < 3 {
		return append(rs, r)
	}
	rs[0] = rs[1]
	rs[1] = rs[2]
	rs[2] = r
	return rs
}

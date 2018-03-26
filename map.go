package corpus

import "encoding/json"

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
func (m *Char3Grams) AddAll(str string) *Char3Grams {
	EachChar3Gram(str, func(str string) {
		m.m[str]++
		m.n++
	})
	return m
}

// Add3Grams adds the 3-grams of anohter map to this.
func (m *Char3Grams) Add3Grams(o *Char3Grams) *Char3Grams {
	for k, v := range o.m {
		m.m[k] += v
	}
	m.n += o.n
	return m
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

// Each iterates over all character 3-grams in this map.
func (m *Char3Grams) Each(f func(string, uint64)) *Char3Grams {
	for k, v := range m.m {
		f(k, v)
	}
	return m
}

type jsonMap struct {
	Total, Len uint64
	NGrams     map[string]uint64
}

// MarshalJSON implements JSON marshaling.
func (m *Char3Grams) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		jsonMap{
			Total:  m.Total(),
			Len:    m.Len(),
			NGrams: m.m,
		})
}

// UnmarshalJSON implements JSON unmarshaling.
func (m *Char3Grams) UnmarshalJSON(bs []byte) error {
	var tmp jsonMap
	if err := json.Unmarshal(bs, &tmp); err != nil {
		return err
	}
	*m = Char3Grams{
		n: tmp.Total,
		m: tmp.NGrams,
	}
	return nil
}

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

// Unigrams represents the absolute unigram frequencies.
type Unigrams struct {
	unigrams map[string]uint64
	total    uint64
}

// Add adds a unigram to the map.
func (u *Unigrams) Add(unigram string) *Unigrams {
	if u.unigrams == nil {
		u.unigrams = make(map[string]uint64)
	}
	u.unigrams[unigram]++
	u.total++
	return u
}

// Total returns the total number of unigrams in the map.
func (u *Unigrams) Total() uint64 {
	return u.total
}

// Get returns the count for the given unigram.
func (u *Unigrams) Get(unigram string) uint64 {
	count, ok := u.unigrams[unigram]
	if !ok {
		return 0
	}
	return count
}

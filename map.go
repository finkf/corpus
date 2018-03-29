package corpus

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

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
func (m *Char3Grams) Each(f func(string, uint64)) {
	if m == nil {
		return
	}
	for k, v := range m.m {
		f(k, v)
	}
}

type jsonMap struct {
	Total, Len uint64
	NGrams     map[string]uint64
}

// MarshalJSON implements JSON marshaling.
func (m *Char3Grams) MarshalJSON() ([]byte, error) {
	return m.marshal(json.Marshal)
}

// UnmarshalJSON implements JSON unmarshaling.
func (m *Char3Grams) UnmarshalJSON(bs []byte) error {
	return m.unmarshal(bs, json.Unmarshal)
}

// GobEncode implement gob marhsaling.
func (m *Char3Grams) GobEncode() ([]byte, error) {
	return m.marshal(marshalGob)
}

// GobDecode implements gob unmarshaling.
func (m *Char3Grams) GobDecode(bs []byte) error {
	return m.unmarshal(bs, unmarshalGob)
}

func (m *Char3Grams) marshal(f marshalFunc) ([]byte, error) {
	return f(jsonMap{
		Total:  m.Total(),
		Len:    m.Len(),
		NGrams: m.m,
	})
}

func (m *Char3Grams) unmarshal(bs []byte, f unmarshalFunc) error {
	var tmp jsonMap
	if err := f(bs, &tmp); err != nil {
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

// Add adds a range of unigrams to the map.
func (u *Unigrams) Add(us ...string) *Unigrams {
	// No check for nil since it is an error to call Add on nil.
	if u.unigrams == nil && len(us) > 0 {
		u.unigrams = make(map[string]uint64)
	}
	for _, unigram := range us {
		u.unigrams[unigram]++
		u.total++
	}
	return u
}

// AddUnigrams adds all unigrams to the map
func (u *Unigrams) AddUnigrams(o *Unigrams) *Unigrams {
	if o == nil {
		return u
	}
	if u.unigrams == nil && len(o.unigrams) > 0 {
		u.unigrams = make(map[string]uint64)
	}
	for k, v := range o.unigrams {
		u.unigrams[k] += v
	}
	u.total += o.total
	return u
}

// Total returns the total number of unigrams in the map.
func (u *Unigrams) Total() uint64 {
	if u == nil {
		return 0
	}
	return u.total
}

// Len returns the total number different unigrams in the map.
func (u *Unigrams) Len() uint64 {
	if u == nil {
		return 0
	}
	return uint64(len(u.unigrams))
}

// Get returns the count for the given unigram.
func (u *Unigrams) Get(unigram string) uint64 {
	if u == nil {
		return 0
	}
	count, ok := u.unigrams[unigram]
	if !ok {
		return 0
	}
	return count
}

// Each calls the supplied callback function for each
// entry in the map.
func (u *Unigrams) Each(f func(string, uint64)) {
	if u == nil {
		return
	}
	for k, v := range u.unigrams {
		f(k, v)
	}
}

type jsonUnigrams struct {
	Total, Len uint64
	Unigrams   map[string]uint64
}

// MarshalJSON implements JSON marshaling.
func (u *Unigrams) MarshalJSON() ([]byte, error) {
	return u.marshal(json.Marshal)
}

// UnmarshalJSON implements JSON unmarshaling.
func (u *Unigrams) UnmarshalJSON(bs []byte) error {
	return u.unmarshal(bs, json.Unmarshal)
}

// GobEncode implement gob marhsaling.
func (u *Unigrams) GobEncode() ([]byte, error) {
	return u.marshal(marshalGob)
}

// GobDecode implements gob unmarshaling.
func (u *Unigrams) GobDecode(bs []byte) error {
	return u.unmarshal(bs, unmarshalGob)
}

func (u *Unigrams) marshal(f marshalFunc) ([]byte, error) {
	return f(
		jsonUnigrams{
			Total:    u.total,
			Len:      u.Len(),
			Unigrams: u.unigrams,
		})
}

func (u *Unigrams) unmarshal(bs []byte, f unmarshalFunc) error {
	var tmp jsonUnigrams
	if err := f(bs, &tmp); err != nil {
		return err
	}
	*u = Unigrams{
		total:    tmp.Total,
		unigrams: tmp.Unigrams,
	}
	return nil
}

// Bigrams represents a map of token 2-grams.
type Bigrams struct {
	bigrams map[string]*Unigrams
	total   uint64
}

// Add adds a range of bigrams to the map.
func (b *Bigrams) Add(bs ...string) *Bigrams {
	// No check for nil since it is an error to call Add on nil.
	if b.bigrams == nil && len(bs) > 1 {
		b.bigrams = make(map[string]*Unigrams)
	}
	for i := 1; i < len(bs); i++ {
		first := bs[i-1]
		second := bs[i]
		if _, ok := b.bigrams[first]; !ok {
			b.bigrams[first] = new(Unigrams)
		}
		b.bigrams[first].Add(second)
		b.total++
	}
	return b
}

// AddBigrams adds all bigrams to the map.
func (b *Bigrams) AddBigrams(o *Bigrams) *Bigrams {
	if o == nil {
		return b
	}
	for k, v := range o.bigrams {
		b.AddUnigrams(k, v)
	}
	return b
}

// AddUnigrams adds the unigrams for the given key into the map
func (b *Bigrams) AddUnigrams(k string, u *Unigrams) *Bigrams {
	if b.bigrams == nil {
		b.bigrams = make(map[string]*Unigrams)
	}
	if _, ok := b.bigrams[k]; !ok {
		b.bigrams[k] = new(Unigrams)
	}
	b.bigrams[k].AddUnigrams(u)
	b.total += u.total
	return b
}

// Total returns the total number of bigrams in the map.
func (b *Bigrams) Total() uint64 {
	if b == nil {
		return 0
	}
	return b.total
}

// Len returns the total number of different unigrams in the map.
func (b *Bigrams) Len() uint64 {
	if b == nil {
		return 0
	}
	return uint64(len(b.bigrams))
}

// Get returns the unigrams for the given head of a bigram.
func (b *Bigrams) Get(first string) *Unigrams {
	if b == nil {
		return nil
	}
	unigrams, ok := b.bigrams[first]
	if !ok {
		return nil
	}
	return unigrams
}

// Each calls the supplied callback function for each
// entry in the map.
func (b *Bigrams) Each(f func(string, *Unigrams)) {
	if b == nil {
		return
	}
	for k, v := range b.bigrams {
		f(k, v)
	}
}

type jsonBigrams struct {
	Total, Len uint64
	Bigrams    map[string]*Unigrams
}

// MarshalJSON implements JSON marshaling.
func (b *Bigrams) MarshalJSON() ([]byte, error) {
	return b.marshal(json.Marshal)
}

// UnmarshalJSON implements JSON unmarshaling.
func (b *Bigrams) UnmarshalJSON(bs []byte) error {
	return b.unmarshal(bs, json.Unmarshal)
}

// GobEncode implement gob marhsaling.
func (b *Bigrams) GobEncode() ([]byte, error) {
	return b.marshal(marshalGob)
}

// GobDecode implements gob unmarshaling.
func (b *Bigrams) GobDecode(bs []byte) error {
	return b.unmarshal(bs, unmarshalGob)
}

func (b *Bigrams) marshal(f marshalFunc) ([]byte, error) {
	return f(
		jsonBigrams{
			Total:   b.total,
			Len:     b.Len(),
			Bigrams: b.bigrams,
		})
}

func (b *Bigrams) unmarshal(bs []byte, f unmarshalFunc) error {
	var tmp jsonBigrams
	if err := f(bs, &tmp); err != nil {
		return err
	}
	*b = Bigrams{
		total:   tmp.Total,
		bigrams: tmp.Bigrams,
	}
	return nil
}

// Trigrams represents a map of token 3-grams.
type Trigrams struct {
	trigrams map[string]*Bigrams
	total    uint64
}

// Add adds a range of trigrams to the map.
func (t *Trigrams) Add(bs ...string) *Trigrams {
	// No check for nil since it is an error to call Add on nil.
	if t.trigrams == nil && len(bs) > 2 {
		t.trigrams = make(map[string]*Bigrams)
	}
	for i := 2; i < len(bs); i++ {
		first := bs[i-2]
		second := bs[i-1]
		third := bs[i]
		if _, ok := t.trigrams[first]; !ok {
			t.trigrams[first] = new(Bigrams)
		}
		t.trigrams[first].Add(second, third)
		t.total++
	}
	return t
}

// AddTrigrams adds the trigrams to the map.
func (t *Trigrams) AddTrigrams(o *Trigrams) *Trigrams {
	if o == nil {
		return t
	}
	for k, v := range o.trigrams {
		t.AddBigrams(k, v)
	}
	return t
}

// AddBigrams adds a key with its bigrams into the map.
func (t *Trigrams) AddBigrams(k string, b *Bigrams) *Trigrams {
	if t.trigrams == nil {
		t.trigrams = make(map[string]*Bigrams)
	}
	if _, ok := t.trigrams[k]; !ok {
		t.trigrams[k] = new(Bigrams)
	}
	t.trigrams[k].AddBigrams(b)
	t.total += b.total
	return t
}

// Total returns the total number of trigrams in the map.
func (t *Trigrams) Total() uint64 {
	if t == nil {
		return 0
	}
	return t.total
}

// Len returns the total number of different bigrams in the map.
func (t *Trigrams) Len() uint64 {
	if t == nil {
		return 0
	}
	return uint64(len(t.trigrams))
}

// Get returns the Bigrams for the given head of a trigram.
func (t *Trigrams) Get(first string) *Bigrams {
	if t == nil {
		return nil
	}
	bigrams, ok := t.trigrams[first]
	if !ok {
		return nil
	}
	return bigrams
}

// Each calls the supplied callback function for each
// entry in the map.
func (t *Trigrams) Each(f func(string, *Bigrams)) {
	if t == nil {
		return
	}
	for k, v := range t.trigrams {
		f(k, v)
	}
}

type jsonTrigrams struct {
	Total, Len uint64
	Trigrams   map[string]*Bigrams
}

// MarshalJSON implements JSON marshaling.
func (t *Trigrams) MarshalJSON() ([]byte, error) {
	return t.marshal(json.Marshal)
}

// UnmarshalJSON implements JSON unmarshaling.
func (t *Trigrams) UnmarshalJSON(bs []byte) error {
	return t.unmarshal(bs, json.Unmarshal)
}

// GobEncode implement gob marhsaling.
func (t *Trigrams) GobEncode() ([]byte, error) {
	return t.marshal(marshalGob)
}

// GobDecode implements gob unmarshaling.
func (t *Trigrams) GobDecode(bs []byte) error {
	return t.unmarshal(bs, unmarshalGob)
}

func (t *Trigrams) marshal(f marshalFunc) ([]byte, error) {
	return f(
		jsonTrigrams{
			Total:    t.total,
			Len:      t.Len(),
			Trigrams: t.trigrams,
		})
}

func (t *Trigrams) unmarshal(bs []byte, f unmarshalFunc) error {
	var tmp jsonTrigrams
	if err := f(bs, &tmp); err != nil {
		return err
	}
	*t = Trigrams{
		total:    tmp.Total,
		trigrams: tmp.Trigrams,
	}
	return nil
}

type marshalFunc func(interface{}) ([]byte, error)
type unmarshalFunc func([]byte, interface{}) error

func marshalGob(data interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func unmarshalGob(bs []byte, data interface{}) error {
	return gob.NewDecoder(bytes.NewBuffer(bs)).Decode(data)
}

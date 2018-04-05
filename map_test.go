package corpus

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"
)

func TestEachChar3Gram(t *testing.T) {
	tests := []struct {
		test string
		want []string
	}{
		{",", nil},
		{"ab", nil},
		{"abc", []string{"abc"}},
		{"abcd", []string{"abc", "bcd"}},
		{"für", []string{"für"}},
		{"Bäume", []string{"Bäu", "äum", "ume"}},
		{"Größe", []string{"Grö", "röß", "öße"}},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			var got []string
			EachChar3Gram(tc.test, func(str string) {
				got = append(got, str)
			})
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

func TestChar3GramMapSizes(t *testing.T) {
	tests := []struct {
		test       string
		total, len uint64
	}{
		{"abababa", 5, 2},
		{"abcde", 3, 3},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			m := new(Char3Grams)
			m.Add(tc.test)
			if got := m.Len(); got != tc.len {
				t.Fatalf("expected %d; got %d", tc.len, got)
			}
			if got := m.Total(); got != tc.total {
				t.Fatalf("expected %d; got %d", tc.total, got)
			}
		})
	}
}

func TestChar3GramMapAdd(t *testing.T) {
	tests := []struct {
		first, second, test string
		want, total         uint64
	}{
		{"abc", "abc", "abc", 2, 2},
		{"abc", "abc", "ab", 0, 2},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s/%s", tc.first, tc.second), func(t *testing.T) {
			a := new(Char3Grams)
			a.Add(tc.first)
			b := new(Char3Grams)
			b.Add(tc.second)
			a.Append(b)
			if got := a.Get(tc.test); got != tc.want {
				t.Fatalf("expected %d; got %d", tc.want, got)
			}
			if got := a.Total(); got != tc.total {
				t.Fatalf("expected %d; got %d", tc.want, got)
			}
		})
	}
}

func TestChar3GramsJSONMarshal(t *testing.T) {
	tests := []struct{ test string }{
		{"abcde"},
		{"Waſſer"},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			a := new(Char3Grams)
			buf := &bytes.Buffer{}
			if err := json.NewEncoder(buf).Encode(a); err != nil {
				t.Fatalf("got error: %v", err)
			}
			var b Char3Grams
			if err := json.NewDecoder(buf).Decode(&b); err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(a, &b) {
				t.Fatalf("expected %v; got %v", a, b)
			}
		})
	}
}

func TestChar3GramsJSONUnarshalError(t *testing.T) {
	tests := []struct{ test string }{
		{`{"Total":"1","NGrams":{"abc":1}}`},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			var m Char3Grams
			err := json.NewDecoder(strings.NewReader(tc.test)).Decode(&m)
			if err == nil {
				t.Fatalf("expected an error; got nil")
			}
		})
	}
}

func TestChar3GramsGobMarshal(t *testing.T) {
	tests := []struct{ test string }{
		{"abcde"},
		{"Waſſer"},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			a := new(Char3Grams)
			buf := &bytes.Buffer{}
			if err := gob.NewEncoder(buf).Encode(a); err != nil {
				t.Fatalf("got error: %v", err)
			}
			var b Char3Grams
			if err := gob.NewDecoder(buf).Decode(&b); err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(a, &b) {
				t.Fatalf("expected %v; got %v", a, b)
			}
		})
	}
}

func TestChar3GramsEach(t *testing.T) {
	tests := []struct {
		test, search string
		want         int
	}{
		{"abca", "a", 2},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			var got int
			new(Char3Grams).Add(tc.test).Each(func(k string, v uint64) {
				if strings.Contains(k, tc.search) {
					got += int(v)
				}
			})
			if got != tc.want {
				t.Fatalf("expected %d; got %d", tc.want, got)
			}
		})
	}
}

func TestUnigrams(t *testing.T) {
	tests := []struct {
		unigrams          *Unigrams
		search            string
		count, total, len uint64
	}{
		{nil, "ab", 0, 0, 0},
		{new(Unigrams).Add("ab", "cd", "ab"), "ab", 2, 3, 2},
		{new(Unigrams).Add("ab", "cd", "ab"), "cd", 1, 3, 2},
		{new(Unigrams).Add("ab", "cd", "ab"), "xy", 0, 3, 2},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s", tc.search), func(t *testing.T) {
			if got := tc.unigrams.Get(tc.search); got != tc.count {
				t.Fatalf("expected %d; got %d", tc.count, got)
			}
			if got := tc.unigrams.Total(); got != tc.total {
				t.Fatalf("expcted %d; got %d", tc.total, got)
			}
			if got := tc.unigrams.Len(); got != tc.len {
				t.Fatalf("expcted %d; got %d", tc.len, got)
			}
		})
	}
}

func TestBigrams(t *testing.T) {
	tests := []struct {
		bigrams           *Bigrams
		first, second     string
		count, total, len uint64
	}{
		{nil, "ab", "cd", 0, 0, 0},
		{new(Bigrams).Add("ab", "cd", "ab"), "ab", "cd", 1, 2, 2},
		{new(Bigrams).Add("ab", "cd", "ab"), "cd", "ab", 1, 2, 2},
		{new(Bigrams).Add("ab", "cd", "ab"), "ab", "xy", 0, 2, 2},
		{new(Bigrams).Add("ab", "cd", "ab"), "xy", "ab", 0, 2, 2},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s %s", tc.first, tc.second), func(t *testing.T) {
			if got := tc.bigrams.Get(tc.first).Get(tc.second); got != tc.count {
				t.Fatalf("expected %d; got %d", tc.count, got)
			}
			if got := tc.bigrams.Total(); got != tc.total {
				t.Fatalf("expcted %d; got %d", tc.total, got)
			}
			if got := tc.bigrams.Len(); got != tc.len {
				t.Fatalf("expcted %d; got %d", tc.len, got)
			}
		})
	}
}

func TestTrigrams(t *testing.T) {
	tests := []struct {
		trigrams             *Trigrams
		first, second, third string
		count, total, len    uint64
	}{
		{nil, "ab", "cd", "ef", 0, 0, 0},
		{new(Trigrams).Add("ab", "cd", "ab"), "ab", "cd", "ab", 1, 1, 1},
		{new(Trigrams).Add("ab", "cd", "ab"), "ab", "xy", "ab", 0, 1, 1},
		{new(Trigrams).Add("ab", "cd", "ab"), "xy", "xy", "ab", 0, 1, 1},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s %s %s", tc.first, tc.second, tc.third), func(t *testing.T) {
			if got := tc.trigrams.Get(tc.first).Get(tc.second).Get(tc.third); got != tc.count {
				t.Fatalf("expected %d; got %d", tc.count, got)
			}
			if got := tc.trigrams.Total(); got != tc.total {
				t.Fatalf("expcted %d; got %d", tc.total, got)
			}
			if got := tc.trigrams.Len(); got != tc.len {
				t.Fatalf("expcted %d; got %d", tc.len, got)
			}
		})
	}
}

func TestAddTrigrams(t *testing.T) {
	tests := []struct {
		test, other  *Trigrams
		f, s, t      string
		count, total uint64
	}{
		{new(Trigrams), nil, "ab", "cd", "ef", 0, 0},
		{new(Trigrams), new(Trigrams).Add("ab", "cd", "ef"), "ab", "cd", "ef", 1, 1},
		{new(Trigrams).Add("ab", "cd", "ef"), nil, "ab", "cd", "ef", 1, 1},
		{new(Trigrams).Add("ab", "cd", "ef"), new(Trigrams).Add("gh", "ij", "kl"), "gh", "ij", "kl", 1, 2},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s %s %s", tc.f, tc.s, tc.t), func(t *testing.T) {
			m := tc.test.Append(tc.other)
			if got := m.Get(tc.f).Get(tc.s).Get(tc.t); got != tc.count {
				log.Printf("m: %v", *m)
				t.Fatalf("expected %d; got %d", tc.count, got)
			}
			if got := m.Total(); got != tc.total {
				t.Fatalf("expceted %d; got %d", tc.total, got)
			}
		})
	}
}

func TestTrigramsJSON(t *testing.T) {
	tests := []struct{ test []string }{
		{},
		{[]string{"ab"}},
		{[]string{"ab", "cd"}},
		{[]string{"ab", "cd", "ef"}},
		{[]string{"ab", "cd", "ef", "ab", "cd", "ef"}},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.test), func(t *testing.T) {
			buf := &bytes.Buffer{}
			u := new(Trigrams).Add(tc.test...)
			if err := json.NewEncoder(buf).Encode(u); err != nil {
				t.Fatalf("got error: %v", err)
			}
			v := new(Trigrams)
			if err := json.NewDecoder(buf).Decode(v); err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(u, v) {
				t.Fatalf("expected %v; got %v", *u, *v)
			}
		})
	}
}

func TestTrigramsGob(t *testing.T) {
	tests := []struct{ test []string }{
		{},
		{[]string{"ab"}},
		{[]string{"ab", "cd"}},
		{[]string{"ab", "cd", "ef"}},
		{[]string{"ab", "cd", "ef", "ab", "cd", "ef"}},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.test), func(t *testing.T) {
			buf := &bytes.Buffer{}
			u := new(Trigrams).Add(tc.test...)
			if err := gob.NewEncoder(buf).Encode(u); err != nil {
				t.Fatalf("got error: %v", err)
			}
			v := new(Trigrams)
			if err := gob.NewDecoder(buf).Decode(v); err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(u, v) {
				t.Fatalf("expected %v; got %v", *u, *v)
			}
		})
	}
}

func TestTrigramsEach(t *testing.T) {
	tests := []struct{ test []string }{
		{},
		{[]string{"ab"}},
		{[]string{"ab", "cd"}},
		{[]string{"ab", "cd", "ef"}},
		{[]string{"ab", "cd", "ef", "ab", "cd", "xy"}},
		{[]string{"ab", "cd", "ef", "ab", "cd", "ef"}},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.test), func(t *testing.T) {
			u := new(Trigrams).Add(tc.test...)
			var es []struct{ triples []string }
			u.Each(func(k string, b *Bigrams) {
				b.Each(func(l string, u *Unigrams) {
					u.Each(func(m string, n uint64) {
						es = append(es, struct{ triples []string }{[]string{k, l, m}})
					})
				})
			})
			for _, e := range es {
				got := u.Get(e.triples[0]).Get(e.triples[1]).Get(e.triples[2])
				if got == 0 {
					t.Fatalf("could not find: %v", e.triples)
				}
			}
		})
	}
}

package corpus

import (
	"bytes"
	"encoding/json"
	"fmt"
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
			m := NewChar3Grams()
			m.AddAll(tc.test)
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
			a := NewChar3Grams()
			a.AddAll(tc.first)
			b := NewChar3Grams()
			b.AddAll(tc.second)
			a.Add3Grams(b)
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
			a := NewChar3Grams().AddAll(tc.test)
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
			NewChar3Grams().AddAll(tc.test).Each(func(k string, v uint64) {
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
		t.Run(fmt.Sprintf("%s %s", tc.first, tc.second), func(t *testing.T) {
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

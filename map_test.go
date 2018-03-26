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

func TestAddUnigrams(t *testing.T) {
	tests := []struct {
		unigrams     []string
		search       string
		count, total uint64
	}{
		{nil, "ab", 0, 0},
		{[]string{"ab", "cd", "ab"}, "ab", 2, 3},
		{[]string{"ab", "cd", "ab"}, "cd", 1, 3},
		{[]string{"ab", "cd", "ab"}, "xy", 0, 3},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.unigrams), func(t *testing.T) {
			u := &Unigrams{}
			for _, unigram := range tc.unigrams {
				u.Add(unigram)
			}
			if got := u.Get(tc.search); got != tc.count {
				t.Fatalf("expected %d; got %d", tc.count, got)
			}
			if got := u.Total(); got != tc.total {
				t.Fatalf("expcted %d; got %d", tc.total, got)
			}
		})
	}
}

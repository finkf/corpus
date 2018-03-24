package corpus

import (
	"fmt"
	"reflect"
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

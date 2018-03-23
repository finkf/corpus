package corpus

import (
	"fmt"
	"strings"
	"testing"
)

func TestReadCharTripleMap(t *testing.T) {
	tests := []struct {
		tokens string
		ts     CharTriple
		want   uint64
	}{
		{"abc", CharTriple{'a', 'b', 'c'}, 1},
		{"abc abc", CharTriple{'a', 'b', 'c'}, 2},
		{"abcde", CharTriple{'a', 'b', 'c'}, 1},
		{"abcde", CharTriple{'b', 'c', 'd'}, 1},
		{"abcde", CharTriple{'c', 'd', 'e'}, 1},
		{"abcde", CharTriple{'c', 'd', 'e'}, 1},
		{"abcde 12 abcde", CharTriple{'c', 'd', 'e'}, 2},
	}
	for _, tc := range tests {
		t.Run(tc.tokens, func(t *testing.T) {
			m, _ := ReadCharTripleMap(sst{tc.tokens}, Word)
			if got := m[tc.ts]; got != tc.want {
				t.Fatalf("expected %d; got %d", tc.want, got)
			}
		})
	}
}

func TestCharTripleMapAdd(t *testing.T) {
	tests := []struct {
		first, second string
		ts            CharTriple
		want          uint64
	}{
		{"abc", "abc", CharTriple{'a', 'b', 'c'}, 2},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%s/%s", tc.first, tc.second), func(t *testing.T) {
			a, _ := ReadCharTripleMap(sst{tc.first}, Word)
			b, _ := ReadCharTripleMap(sst{tc.second}, Word)
			m := a.Add(b)
			if got := m[tc.ts]; got != tc.want {
				t.Fatalf("expected %d; got %d", tc.want, got)
			}
		})
	}
}

type sst struct {
	tokens string
}

func (s sst) Err() error { return nil }

func (s sst) Tokens() <-chan Token {
	strs := strings.Split(s.tokens, " ")
	chn := make(chan Token, len(strs))
	go func() {
		defer close(chn)
		for _, str := range strs {
			chn <- Token{Str: str}
		}
	}()
	return chn
}

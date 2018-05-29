// +build ignore

package corpus

import (
	"reflect"
	"testing"
)

func TestSplitter(t *testing.T) {
	tests := []struct {
		test string
		want []string
	}{
		{"", nil},
		{"abc", []string{"abc"}},
		{"a,b,c", []string{"a", ",", "b", ",", "c"}},
		{"abc-def", []string{"abc", "-", "def"}},
		{"abc---def()", []string{"abc", "---", "def", "()"}},
		{"(abc)", []string{"(", "abc", ")"}},
		{"03.09.1983", []string{"03", ".", "09", ".", "1983"}},
		{"ochſen-fleiſch,", []string{"ochſen", "-", "fleiſch", ","}},
		{"fuͤr,abc", []string{"fuͤr", ",", "abc"}},
	}
	for _, tc := range tests {
		t.Run(tc.test, func(t *testing.T) {
			got := new(splitter).split(tc.test)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

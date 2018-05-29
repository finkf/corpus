// +build go1.9 go1.10

package corpus_test

import "testing"

func checkType(t *testing.T, name string, want, got bool) {
	t.Helper()
	if want != got {
		t.Errorf("expected %s = %t; got %t", name, want, got)
	}
}

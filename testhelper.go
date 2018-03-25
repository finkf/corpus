// +build !go1.9 !go1.10

package corpus

import (
	"io"
	"os"
	"testing"
)

func openDTATestFile(t *testing.T) io.ReadCloser {
	is, err := os.Open("testdata/dta.xml")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	return is
}

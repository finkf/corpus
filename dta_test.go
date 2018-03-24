package corpus

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"testing/iotest"
)

// we do not care about closing the file.
func openDTATestFile(t *testing.T) io.ReadCloser {
	t.Helper()
	is, err := os.Open("testdata/dta.xml")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	return is
}

func TestDTAErrorRead(t *testing.T) {
	r := iotest.TimeoutReader(openDTATestFile(t))
	err := dtaReadTokens(r, func(Token) {})
	if err == nil {
		t.Fatalf("expected an error; got nil")
	}
}

func TestDTARead(t *testing.T) {
	tests := []struct {
		want       []string
		skip, take int
	}{
		{[]string{"D.", "Henrici"}, 0, 2},
		{[]string{"Leib-Medicus", "Der", "Studenten", ","}, 6, 4},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.want), func(t *testing.T) {
			r := openDTATestFile(t)
			var got []string
			var skipped, taken int
			err := dtaReadTokensAndClose(r, func(t Token) {
				if skipped < tc.skip {
					skipped++
					return
				}
				if taken < tc.take {
					got = append(got, string(t))
					taken++

				}
			})
			if err != nil {
				t.Fatalf("got error: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected %v; got %v", tc.want, got)
			}
		})
	}
}

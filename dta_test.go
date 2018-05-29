// +build ignore

package corpus

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
	"testing/iotest"

	"github.com/pkg/errors"
)

func closeError(t *testing.T) io.ReadCloser {
	return errCloser(openDTATestFile(t))
}

func readError(t *testing.T) io.ReadCloser {
	return ioutil.NopCloser(iotest.TimeoutReader(openDTATestFile(t)))
}

func readAndCloseError(t *testing.T) io.ReadCloser {
	return errCloser(readError(t))
}

func TestDTAErrors(t *testing.T) {
	tests := []struct {
		reader func(*testing.T) io.ReadCloser
		want   error
	}{
		{closeError, errClose},
		{readError, iotest.ErrTimeout},
		{readAndCloseError, iotest.ErrTimeout},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%T", tc.reader), func(t *testing.T) {
			r := tc.reader(t)
			err := DTAReadTokensAndClose(r, func(Token) {})
			if errors.Cause(err) != tc.want {
				t.Fatalf("expected %s; got %v", tc.want, err)
			}
		})
	}
}

func TestDTARead(t *testing.T) {
	tests := []struct {
		want       []string
		skip, take int
	}{
		{[]string{"D", ".", "Henrici"}, 0, 3},
		{[]string{"Leib", "-", "Medicus", "Der", "Studenten", ","}, 7, 6},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc.want), func(t *testing.T) {
			r := openDTATestFile(t)
			var got []string
			var skipped, taken int
			err := DTAReadTokensAndClose(r, func(t Token) {
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

func errCloser(r io.Reader) io.ReadCloser {
	return errCloserS{r}
}

type errCloserS struct {
	r io.Reader
}

var errClose = errors.New("close")

func (r errCloserS) Read(bs []byte) (int, error) { return r.r.Read(bs) }
func (r errCloserS) Close() error                { return errClose }

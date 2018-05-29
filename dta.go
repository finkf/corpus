package corpus

import (
	"encoding/xml"
	"io"
	"os"

	"github.com/pkg/errors"
)

// DTA denotes a handle for DTA files.
type DTA struct {
	r io.Reader
}

// NewDTAFile returns an new DTA handle that reads from
// the given file name.
// Do not forget to close the returned DTA handle after
// usage.
func NewDTAFile(path string) (DTA, error) {
	in, err := os.Open(path)
	if err != nil {
		return DTA{}, err
	}
	return NewDTA(in), nil
}

// NewDTA returns a new DTA handle.
func NewDTA(r io.Reader) DTA {
	return DTA{r}
}

// EachXMLToken calls the given callback function for
// each XML token in the dta document.
func (dta DTA) EachXMLToken(f func(xml.Token)) error {
	d := xml.NewDecoder(dta.r)
	var err error
	var t xml.Token
	for t, err = d.Token(); err == nil; t, err = d.Token() {
		f(t)
	}
	if err == io.EOF {
		return nil
	}
	// returns nil if err == nil
	return errors.Wrapf(err, "invalid dta corpus file")
}

// Tokenize implements the tokenize interface for DTA handles.
func (dta DTA) Tokenize(f func(string)) error {
	var inToken bool
	return dta.EachXMLToken(func(t xml.Token) {
		switch tt := t.(type) {
		case xml.CharData:
			if inToken {
				tokenize(string(tt), f)
			}
		case xml.StartElement:
			inToken = tt.Name.Local == "token"
		case xml.EndElement:
			inToken = false
		}
	})
}

// Close the underlying reader of the DTA handle.
func (dta DTA) Close() error {
	switch c := dta.r.(type) {
	case io.Closer:
		return c.Close()
	}
	return nil
}

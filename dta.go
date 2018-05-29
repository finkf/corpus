// +build ignore

package corpus

import (
	"encoding/xml"
	"io"

	"github.com/pkg/errors"
)

// DTAReadTokensAndClose is a conveniece function that reads
// all tokens in a DTA file and closes the reader.
func DTAReadTokensAndClose(r io.ReadCloser, f func(Token)) (err error) {
	defer func() {
		e2 := r.Close()
		if e2 != nil && err == nil {
			err = e2
		}
	}()
	err = DTAReadTokens(r, f)
	return
}

// DTAReadTokens reads all tokens form a DTA corpus file.
func DTAReadTokens(r io.Reader, f func(Token)) error {
	d := xml.NewDecoder(r)
	var err error
	var t xml.Token
	var inToken bool
	for t, err = d.Token(); err == nil; t, err = d.Token() {
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
	}
	if err == io.EOF {
		return nil
	}
	return errors.Wrapf(err, "invalid dta corpus file")
}

func tokenize(t string, f func(Token)) {
	for _, s := range new(splitter).split(t) {
		f(Token(s))
	}
}

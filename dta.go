package corpus

import (
	"encoding/xml"
	"io"

	"github.com/pkg/errors"
)

// DTAReadTokensAndClose is a conveniece function that reads
// all tokens in a DTA file and closes the reader.
// The error of the call to close is ignored.
func DTAReadTokensAndClose(r io.ReadCloser, f func(Token)) error {
	defer func() { _ = r.Close() }()
	return DTAReadTokens(r, f)
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
	f(Token(t))
}

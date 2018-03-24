package corpus

import (
	"encoding/xml"
	"io"

	"github.com/pkg/errors"
)

func dtaReadTokensAndClose(r io.ReadCloser, f func(Token)) error {
	defer func() { _ = r.Close() }()
	return dtaReadTokens(r, f)
}

func dtaReadTokens(r io.Reader, f func(Token)) error {
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

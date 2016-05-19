// Package href extracts HTML <a href> attributes.
package href

import (
	"io"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Extractor extracts <a href> links from HTML.
type Extractor struct {
	z *html.Tokenizer
}

// New returns an Extractor that reads from r.
func New(r io.Reader) *Extractor {
	z := html.NewTokenizer(r)
	return &Extractor{z: z}
}

// Next returns the next link, or io.EOF.
func (e *Extractor) Next() (string, error) {
	for {
		tt := e.z.Next()
		switch tt {
		case html.ErrorToken:
			err := e.z.Err()
			return "", err
		case html.StartTagToken:
			t := e.z.Token()
			if t.DataAtom != atom.A {
				continue
			}
			for _, attr := range t.Attr {
				if attr.Namespace != "" {
					continue
				}
				if attr.Key != "href" {
					continue
				}
				return attr.Val, nil
			}
		}
	}
}

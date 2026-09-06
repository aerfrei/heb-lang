// Package tanach parses tanach.us book XML files into per-word rows.
package tanach

import (
	"encoding/xml"
	"io"
	"strconv"

	"github.com/aerfrei/heb-lang/internal/hebrew"
)

type Word struct {
	Book     string
	Chapter  int
	Verse    int
	Position int
	Letters  string
	Text     string // original pointed text, before stripping to bare letters
}

// ParseBook streams a book's XML and returns one Word per word, in order.
// Where a verse has a kethiv/qere pair (<k> then <q>), the qere is used
// and the kethiv is discarded.
func ParseBook(r io.Reader, book string) ([]Word, error) {
	dec := xml.NewDecoder(r)
	var words []Word
	var chapter, verse, pos int

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		start, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}

		switch start.Name.Local {
		case "teiHeader":
			if err := dec.Skip(); err != nil {
				return nil, err
			}
		case "c":
			chapter = attrInt(start, "n")
		case "v":
			verse = attrInt(start, "n")
			pos = 0
		case "w", "q":
			var text string
			if err := dec.DecodeElement(&text, &start); err != nil {
				return nil, err
			}
			pos++
			words = append(words, Word{
				Book:     book,
				Chapter:  chapter,
				Verse:    verse,
				Position: pos,
				Letters:  hebrew.StripToLetters(text),
				Text:     text,
			})
		case "k":
			// Kethiv paired with the following <q>; the qere wins, so discard this.
			var discard string
			if err := dec.DecodeElement(&discard, &start); err != nil {
				return nil, err
			}
		}
	}

	return words, nil
}

func attrInt(start xml.StartElement, name string) int {
	for _, a := range start.Attr {
		if a.Name.Local == name {
			n, _ := strconv.Atoi(a.Value)
			return n
		}
	}
	return 0
}

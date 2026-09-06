// Command gencitations builds wordle/citations.json: for every 5-letter word
// found in the Tanach, its first occurrence (book, chapter, verse, full
// pointed verse text), for the wordle game to fetch on demand.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aerfrei/heb-lang/internal/tanach"
)

const (
	booksDir = "tanach"
	outPath  = "wordle/citations.json"
)

// canonicalOrder is the traditional Tanach (Torah/Nevi'im/Ketuvim) book
// order, so "first occurrence" means first in the Bible, not alphabetical.
var canonicalOrder = []string{
	"Genesis", "Exodus", "Leviticus", "Numbers", "Deuteronomy",
	"Joshua", "Judges", "Samuel_1", "Samuel_2", "Kings_1", "Kings_2",
	"Isaiah", "Jeremiah", "Ezekiel",
	"Hosea", "Joel", "Amos", "Obadiah", "Jonah", "Micah", "Nahum",
	"Habakkuk", "Zephaniah", "Haggai", "Zechariah", "Malachi",
	"Psalms", "Proverbs", "Job", "Song_of_Songs", "Ruth", "Lamentations",
	"Ecclesiastes", "Esther", "Daniel", "Ezra", "Nehemiah",
	"Chronicles_1", "Chronicles_2",
}

type verseKey struct {
	book    string
	chapter int
	verse   int
}

// verse is one entry of the "verses" array in the output JSON: [book, chapter, verse, text].
type verse [4]any

func main() {
	var (
		verseOrder []verseKey // every verse seen, in canonical order (first-seen)
		verseSeen  = map[verseKey]bool{}
		verseText  = map[verseKey]string{}
		firstVerse = map[string]verseKey{} // 5-letter word -> its first verse
		count      = map[string]int{}      // 5-letter word -> total occurrences, any pointing
	)

	for _, book := range canonicalOrder {
		path := filepath.Join(booksDir, book+".xml")
		f, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		ws, err := tanach.ParseBook(f, book)
		f.Close()
		if err != nil {
			log.Fatalf("%s: %v", book, err)
		}

		for _, w := range ws {
			key := verseKey{w.Book, w.Chapter, w.Verse}
			if !verseSeen[key] {
				verseSeen[key] = true
				verseOrder = append(verseOrder, key)
			}
			if verseText[key] == "" {
				verseText[key] = w.Text
			} else {
				verseText[key] += " " + w.Text
			}

			if tileCount(w.Letters) == 5 {
				count[w.Letters]++
				if _, ok := firstVerse[w.Letters]; !ok {
					firstVerse[w.Letters] = key
				}
			}
		}
	}

	// Only verses actually cited as some word's first occurrence make it
	// into the output — most of the Tanach's text is never a citation
	// target, so there's no reason to ship it.
	cited := map[verseKey]bool{}
	for _, key := range firstVerse {
		cited[key] = true
	}

	var verses []verse
	verseIndex := map[verseKey]int{}
	for _, key := range verseOrder {
		if !cited[key] {
			continue
		}
		verseIndex[key] = len(verses)
		verses = append(verses, verse{key.book, key.chapter, key.verse, verseText[key]})
	}

	// Each word maps to [firstVerseIndex, totalOccurrences] — the count
	// covers every pointing/spelling that reduces to the same 5 letters.
	words := make(map[string][2]int, len(firstVerse))
	for w, key := range firstVerse {
		words[w] = [2]int{verseIndex[key], count[w]}
	}

	out := struct {
		Verses []verse           `json:"verses"`
		Words  map[string][2]int `json:"words"`
	}{verses, words}

	f, err := os.Create(outPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(out); err != nil {
		log.Fatal(err)
	}

	info, _ := os.Stat(outPath)
	fmt.Printf("%s: %d verses, %d words, %d bytes\n", outPath, len(verses), len(words), info.Size())
}

// tileCount counts game tiles, not runes: a shin/sin dot (U+05C1/U+05C2)
// rides with the preceding consonant rather than forming its own tile.
func tileCount(letters string) int {
	n := 0
	for _, r := range letters {
		if r != 0x05C1 && r != 0x05C2 {
			n++
		}
	}
	return n
}

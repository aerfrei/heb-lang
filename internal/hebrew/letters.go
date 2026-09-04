// Package hebrew provides helpers for normalizing pointed Hebrew text.
package hebrew

import "strings"

// StripToLetters reduces a pointed Hebrew word to its bare consonants.
// It keeps the shin/sin dot (U+05C1/U+05C2), since that distinguishes two
// different letters, but drops vowel points, dagesh, cantillation marks,
// and any other punctuation.
func StripToLetters(word string) string {
	var b strings.Builder
	for _, r := range word {
		switch {
		case r >= 0x05D0 && r <= 0x05EA: // Hebrew consonants
			b.WriteRune(r)
		case r == 0x05C1 || r == 0x05C2: // shin dot / sin dot
			b.WriteRune(r)
		}
	}
	return b.String()
}

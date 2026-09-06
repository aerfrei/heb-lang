// Package hebrew provides helpers for normalizing pointed Hebrew text.
package hebrew

import "strings"

// finalToRegular maps each final letter form to its regular form. Hebrew
// orthography always uses the final form at the true end of a word, so this
// loses no distinguishing information — it's as cosmetic as the vowel
// points StripToLetters already drops.
var finalToRegular = map[rune]rune{
	'ך': 'כ',
	'ם': 'מ',
	'ן': 'נ',
	'ף': 'פ',
	'ץ': 'צ',
}

// StripToLetters reduces a pointed Hebrew word to its bare consonants.
// It keeps the shin/sin dot (U+05C1/U+05C2), since that distinguishes two
// different letters, but drops vowel points, dagesh, cantillation marks,
// and any other punctuation. Final letter forms are normalized to their
// regular form.
func StripToLetters(word string) string {
	var b strings.Builder
	for _, r := range word {
		switch {
		case r >= 0x05D0 && r <= 0x05EA: // Hebrew consonants
			if reg, ok := finalToRegular[r]; ok {
				r = reg
			}
			b.WriteRune(r)
		case r == 0x05C1 || r == 0x05C2: // shin dot / sin dot
			b.WriteRune(r)
		}
	}
	return b.String()
}

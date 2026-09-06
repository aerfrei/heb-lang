package hebrew

import "testing"

func TestStripToLetters(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"drops vowels and cantillation", "וַיְהִ֗י", "ויהי"},
		{"drops dagesh", "בִּימֵי֙", "בימי"},
		{"keeps shin dot", "שְׁפֹ֣ט", "שׁפט"},
		{"normalizes final mem", "אַבְרָהָ֖ם", "אברהמ"},
		{"normalizes final nun", "לָכֵ֜ן", "לכנ"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := StripToLetters(c.in)
			if got != c.want {
				t.Fatalf("StripToLetters(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

package tanach

import (
	"strings"
	"testing"
)

const sampleXML = `<?xml version="1.0" encoding="UTF-8"?>
<Tanach>
<teiHeader>
  <fileDesc><n>ignored, should never collide with body tags</n></fileDesc>
</teiHeader>
<tanach>
<book>
<c n="1">
<v n="1">
<w>וַיְהִ֗י</w>
<w>בִּימֵי֙</w>
</v>
<v n="2">
<w>שְׁפֹ֣ט</w>
<k>הגיאות</k>
<q>הַגֵּאָי֑וֹת</q>
</v>
</c>
</book>
</tanach>
</Tanach>`

func TestParseBook(t *testing.T) {
	words, err := ParseBook(strings.NewReader(sampleXML), "TestBook")
	if err != nil {
		t.Fatalf("ParseBook: %v", err)
	}

	want := []Word{
		{"TestBook", 1, 1, 1, "ויהי", "וַיְהִ֗י"},
		{"TestBook", 1, 1, 2, "בימי", "בִּימֵי֙"},
		{"TestBook", 1, 2, 1, "שׁפט", "שְׁפֹ֣ט"},
		{"TestBook", 1, 2, 2, "הגאיות", "הַגֵּאָי֑וֹת"}, // qere, not the kethiv
	}

	if len(words) != len(want) {
		t.Fatalf("got %d words, want %d: %+v", len(words), len(want), words)
	}
	for i, w := range want {
		if words[i] != w {
			t.Errorf("word %d = %+v, want %+v", i, words[i], w)
		}
	}
}

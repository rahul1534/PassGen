package web

import (
	"strings"
	"testing"
)

func TestSplitRunsGroupsByKind(t *testing.T) {
	got := splitRuns("ab12-cD")
	want := []run{
		{"ab", kindLetter},
		{"12", kindDigit},
		{"-", kindSymbol},
		{"cD", kindLetter},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d runs %v, want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("run %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestSplitRunsEmpty(t *testing.T) {
	if got := splitRuns(""); len(got) != 0 {
		t.Fatalf("expected no runs, got %v", got)
	}
}

// The runs must always reassemble to the exact original string, including
// spaces, multi-byte characters and repeated boundaries.
func TestSplitRunsRoundTrip(t *testing.T) {
	for _, s := range []string{
		"J.$,Kype*Fy^M39fC;kT",
		"correct-horse-battery-staple7!",
		"  spaced  out ",
		"héllo→wörld٣",
		"0000",
		"!!!!",
		"a",
	} {
		var b strings.Builder
		for _, r := range splitRuns(s) {
			if r.Text == "" {
				t.Errorf("%q produced an empty run", s)
			}
			b.WriteString(r.Text)
		}
		if b.String() != s {
			t.Errorf("round trip %q -> %q", s, b.String())
		}
	}
}

func TestKindClass(t *testing.T) {
	if kindLetter.class() != "" || kindDigit.class() != "ch-digit" || kindSymbol.class() != "ch-symbol" {
		t.Fatal("unexpected CSS classes")
	}
}

package web

import "unicode"

// charKind classifies a character for display so the UI can colour digits and
// symbols differently from letters. This file has no build constraint so the
// logic can be unit-tested natively.
type charKind int

const (
	kindLetter charKind = iota
	kindDigit
	kindSymbol
)

// class returns the CSS class for the kind ("" for plain letters).
func (k charKind) class() string {
	switch k {
	case kindDigit:
		return "ch-digit"
	case kindSymbol:
		return "ch-symbol"
	default:
		return ""
	}
}

// run is a maximal sequence of consecutive characters of the same kind.
type run struct {
	Text string
	Kind charKind
}

func classifyRune(r rune) charKind {
	switch {
	case r >= '0' && r <= '9':
		return kindDigit
	case unicode.IsLetter(r):
		return kindLetter
	default:
		return kindSymbol
	}
}

// splitRuns groups s into runs of same-kind characters. Concatenating the
// runs' Text always reproduces s exactly.
func splitRuns(s string) []run {
	var (
		runs []run
		cur  []rune
		kind charKind
	)
	flush := func() {
		if len(cur) > 0 {
			runs = append(runs, run{Text: string(cur), Kind: kind})
			cur = cur[:0]
		}
	}
	for _, r := range s {
		k := classifyRune(r)
		if len(cur) > 0 && k != kind {
			flush()
		}
		kind = k
		cur = append(cur, r)
	}
	flush()
	return runs
}

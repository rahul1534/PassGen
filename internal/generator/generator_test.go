package generator

import (
	"strings"
	"testing"
	"unicode"

	"github.com/rahultaneja/PassGen/internal/random"
)

func testWords() []string {
	words := make([]string, 4096)
	for i := range words {
		words[i] = "word" + string(rune('a'+i%26))
	}
	return words
}

func TestGeneratePasswordLength(t *testing.T) {
	src := random.NewDeterministicSource(1, 2, 3, 4, 5, 6, 7, 8)
	opts := DefaultPasswordOptions()

	for _, length := range []int{8, 16, 32, 64} {
		opts.Length = length
		password, err := GeneratePassword(src, opts)
		if err != nil {
			t.Fatalf("length %d: %v", length, err)
		}
		if len(password) != length {
			t.Fatalf("expected length %d, got %d (%q)", length, len(password), password)
		}
	}
}

func TestGeneratePasswordMinimums(t *testing.T) {
	src := random.NewDeterministicSource(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	opts := PasswordOptions{
		Length:         20,
		Uppercase:      true,
		Lowercase:      true,
		Numbers:        true,
		Symbols:        true,
		MinUppercase:   2,
		MinLowercase:   2,
		MinNumbers:     2,
		MinSymbols:     2,
		ExcludeAmbiguous: true,
	}

	password, err := GeneratePassword(src, opts)
	if err != nil {
		t.Fatal(err)
	}

	cs, _ := BuildCharSets(opts)
	upper, lower, nums, syms := CountRequirements(password, opts, cs)
	if upper < 2 || lower < 2 || nums < 2 || syms < 2 {
		t.Fatalf("minimums not met: upper=%d lower=%d nums=%d syms=%d password=%q", upper, lower, nums, syms, password)
	}
}

func TestValidatePasswordOptions(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.Uppercase, opts.Lowercase, opts.Numbers, opts.Symbols = false, false, false, false
	if err := ValidatePasswordOptions(opts); err != ErrNoCharacterSet {
		t.Fatalf("expected ErrNoCharacterSet, got %v", err)
	}

	opts = DefaultPasswordOptions()
	opts.Length = 8
	opts.MinUppercase = 5
	opts.MinLowercase = 5
	if err := ValidatePasswordOptions(opts); err != ErrImpossibleMinimums {
		t.Fatalf("expected ErrImpossibleMinimums, got %v", err)
	}
}

func TestExcludedCharacters(t *testing.T) {
	src := random.NewDeterministicSource(0, 1, 2, 3, 4, 5, 6, 7, 8, 9)
	opts := DefaultPasswordOptions()
	opts.Length = 16
	opts.ExcludedCharacters = "@#$"
	opts.Symbols = true

	password, err := GeneratePassword(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range "@#$" {
		if strings.ContainsRune(password, ch) {
			t.Fatalf("password contains excluded char %q: %q", ch, password)
		}
	}
}

func TestExcludeAmbiguous(t *testing.T) {
	cs, err := BuildCharSets(PasswordOptions{
		Lowercase: true, Uppercase: true, Numbers: true,
		ExcludeAmbiguous: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, ch := range "0Oo1Il" {
		if strings.ContainsRune(cs.Combined, ch) {
			t.Fatalf("ambiguous char %q in combined set %q", ch, cs.Combined)
		}
	}
}

func TestPreventRepeatedCharacters(t *testing.T) {
	src := random.NewDeterministicSource(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19)
	opts := DefaultPasswordOptions()
	opts.Length = 12
	opts.PreventRepeated = true

	password, err := GeneratePassword(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	seen := make(map[rune]bool)
	for _, r := range password {
		if seen[r] {
			t.Fatalf("repeated character %q in %q", r, password)
		}
		seen[r] = true
	}
}

func TestGeneratePassphrase(t *testing.T) {
	SetWordListForTests(testWords())
	opts := DefaultPassphraseOptions()

	src := random.NewDeterministicSource(100, 200, 300, 400, 500)
	passphrase, err := GeneratePassphrase(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(strings.TrimRightFunc(passphrase, func(r rune) bool {
		return unicode.IsDigit(r) || strings.ContainsRune(SymbolChars, r)
	}), opts.Separator)
	if len(parts) != opts.Words {
		t.Fatalf("expected %d words, got %q", opts.Words, passphrase)
	}
}

func TestGeneratePIN(t *testing.T) {
	src := random.NewDeterministicSource(1, 2, 3, 4, 5, 6)
	opts := DefaultPINOptions()

	pin, err := GeneratePIN(src, opts)
	if err != nil {
		t.Fatal(err)
	}
	if len(pin) != opts.Length {
		t.Fatalf("expected length %d, got %q", opts.Length, pin)
	}
	if !ContainsOnlyDigits(pin) {
		t.Fatalf("PIN contains non-digits: %q", pin)
	}
}

func TestStrengthLevels(t *testing.T) {
	cases := []struct {
		entropy float64
		level   StrengthLevel
	}{
		{10, StrengthVeryWeak},
		{30, StrengthWeak},
		{50, StrengthFair},
		{70, StrengthStrong},
		{100, StrengthVeryStrong},
	}
	for _, tc := range cases {
		result := StrengthFromEntropy(tc.entropy)
		if result.Level != tc.level {
			t.Fatalf("entropy %.0f: expected %s, got %s", tc.entropy, tc.level, result.Level)
		}
	}
}

func TestShuffleChangesOrder(t *testing.T) {
	src := random.NewDeterministicSource(3, 1, 2)
	items := []rune{'a', 'b', 'c', 'd', 'e'}
	original := string(items)
	if err := random.Shuffle(src, items); err != nil {
		t.Fatal(err)
	}
	if string(items) == original {
		t.Fatalf("shuffle did not change order")
	}
}

func TestDeterministicSourceNotIdentical(t *testing.T) {
	src1 := random.NewDeterministicSource(1, 5, 9, 2, 7)
	src2 := random.NewDeterministicSource(3, 8, 1, 6, 4)
	opts := DefaultPasswordOptions()
	opts.Length = 20

	p1, err := GeneratePassword(src1, opts)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := GeneratePassword(src2, opts)
	if err != nil {
		t.Fatal(err)
	}
	if p1 == p2 {
		t.Fatalf("expected different passwords from different sources")
	}
}

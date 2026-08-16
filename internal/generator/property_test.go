package generator

import (
	"strings"
	"testing"
	"unicode"

	"github.com/rahul1534/PassGen/internal/random"
)

// cryptoLikeSource uses crypto/rand via a simple counter-free wrapper for property tests.
// For non-WASM tests we use DeterministicSource with a large varied sequence.
func propertySource(seed int) *random.DeterministicSource {
	values := make([]int, 4096)
	x := uint32(seed)*1664525 + 1013904223
	for i := range values {
		x = x*1664525 + 1013904223
		values[i] = int(x)
	}
	return random.NewDeterministicSource(values...)
}

func TestPropertyPasswordRequirements(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.ExcludeAmbiguous = false
	opts.Length = 24
	opts.MinUppercase = 2
	opts.MinLowercase = 2
	opts.MinNumbers = 2
	opts.MinSymbols = 2

	cs, err := BuildCharSets(opts)
	if err != nil {
		t.Fatal(err)
	}

	for seed := 1; seed <= 500; seed++ {
		src := propertySource(seed)
		password, err := GeneratePassword(src, opts)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		if len([]rune(password)) != opts.Length {
			t.Fatalf("seed %d: length want %d got %d (%q)", seed, opts.Length, len([]rune(password)), password)
		}
		upper, lower, nums, syms := CountRequirements(password, opts, cs)
		if upper < opts.MinUppercase || lower < opts.MinLowercase || nums < opts.MinNumbers || syms < opts.MinSymbols {
			t.Fatalf("seed %d: mins not met u=%d l=%d n=%d s=%d password=%q", seed, upper, lower, nums, syms, password)
		}
		for _, r := range opts.ExcludedCharacters {
			if strings.ContainsRune(password, r) {
				t.Fatalf("seed %d: excluded %q present in %q", seed, r, password)
			}
		}
	}
}

func TestPropertyExclusionsAndAmbiguous(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.Length = 32
	opts.ExcludeAmbiguous = true
	opts.ExcludedCharacters = "@#"

	for seed := 1; seed <= 200; seed++ {
		src := propertySource(seed + 1000)
		password, err := GeneratePassword(src, opts)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		for _, r := range "0Oo1Il@#" {
			if strings.ContainsRune(password, r) {
				t.Fatalf("seed %d: forbidden %q in %q", seed, r, password)
			}
		}
	}
}

func TestPropertyNoRepeated(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.Length = 16
	opts.PreventRepeated = true
	opts.ExcludeAmbiguous = false

	for seed := 1; seed <= 200; seed++ {
		src := propertySource(seed + 2000)
		password, err := GeneratePassword(src, opts)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		seen := make(map[rune]bool)
		for _, r := range password {
			if seen[r] {
				t.Fatalf("seed %d: repeated %q in %q", seed, r, password)
			}
			seen[r] = true
		}
	}
}

func TestPropertyImpossibleConfigs(t *testing.T) {
	cases := []PasswordOptions{
		{Length: 8, Uppercase: true, Lowercase: true, MinUppercase: 5, MinLowercase: 5},
		{Length: 10, Uppercase: false, Lowercase: false, Numbers: false, Symbols: false},
		{Length: 3, Uppercase: true, Lowercase: true, Numbers: true, Symbols: true},
		{Length: 100, Uppercase: true, Lowercase: true, Numbers: true, Symbols: true, PreventRepeated: true, ExcludeAmbiguous: true},
	}
	src := propertySource(42)
	for i, opts := range cases {
		_, err := GeneratePassword(src, opts)
		if err == nil {
			t.Fatalf("case %d: expected error", i)
		}
	}
}

func TestPropertyPIN(t *testing.T) {
	opts := DefaultPINOptions()
	opts.Length = 8
	opts.AllowRepeatedDigits = false

	for seed := 1; seed <= 200; seed++ {
		src := propertySource(seed + 3000)
		pin, err := GeneratePIN(src, opts)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		if len(pin) != opts.Length || !ContainsOnlyDigits(pin) {
			t.Fatalf("seed %d: invalid pin %q", seed, pin)
		}
		seen := make(map[byte]bool)
		for i := 0; i < len(pin); i++ {
			if seen[pin[i]] {
				t.Fatalf("seed %d: repeated digit in %q", seed, pin)
			}
			seen[pin[i]] = true
		}
	}
}

func TestPropertyPassphraseWordCount(t *testing.T) {
	SetWordListForTests(testWords())
	opts := DefaultPassphraseOptions()
	opts.AddNumber = true
	opts.AddSymbol = false

	for seed := 1; seed <= 100; seed++ {
		src := propertySource(seed + 4000)
		phrase, err := GeneratePassphrase(src, opts)
		if err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
		base := strings.TrimRightFunc(phrase, func(r rune) bool {
			return unicode.IsDigit(r) || strings.ContainsRune(SymbolChars, r)
		})
		parts := strings.Split(base, opts.Separator)
		if len(parts) != opts.Words {
			t.Fatalf("seed %d: want %d words, got %q", seed, opts.Words, phrase)
		}
		if !unicode.IsDigit(rune(phrase[len(phrase)-1])) {
			t.Fatalf("seed %d: expected trailing digit in %q", seed, phrase)
		}
	}
}

func TestFuzzPasswordSmoke(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping fuzz smoke in short mode")
	}
	opts := DefaultPasswordOptions()
	opts.ExcludeAmbiguous = false
	src := propertySource(99)
	for i := 0; i < 2000; i++ {
		password, err := GeneratePassword(src, opts)
		if err != nil {
			t.Fatal(err)
		}
		if len(password) != opts.Length {
			t.Fatalf("unexpected length %d", len(password))
		}
	}
}

package generator

import (
	"strings"

	"github.com/rahultaneja/PassGen/internal/random"
)

// PasswordOptions configures random password generation.
type PasswordOptions struct {
	Length             int
	Uppercase          bool
	Lowercase          bool
	Numbers            bool
	Symbols            bool
	MinUppercase       int
	MinLowercase       int
	MinNumbers         int
	MinSymbols         int
	ExcludeAmbiguous   bool
	ExcludedCharacters string
	PreventRepeated    bool
}

// DefaultPasswordOptions returns quality defaults from the specification.
func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{
		Length:           20,
		Uppercase:        true,
		Lowercase:        true,
		Numbers:          true,
		Symbols:          true,
		MinUppercase:     1,
		MinLowercase:     1,
		MinNumbers:       1,
		MinSymbols:       1,
		ExcludeAmbiguous: true,
	}
}

// StrongPasswordOptions returns simplified strong password settings.
func StrongPasswordOptions(length int) PasswordOptions {
	opts := DefaultPasswordOptions()
	opts.Length = length
	opts.MinUppercase = 1
	opts.MinLowercase = 1
	opts.MinNumbers = 1
	opts.MinSymbols = 1
	opts.ExcludeAmbiguous = true
	return opts
}

// GeneratePassword creates a password satisfying all configured requirements.
func GeneratePassword(src random.Source, opts PasswordOptions) (string, error) {
	if err := ValidatePasswordOptions(opts); err != nil {
		return "", err
	}

	cs, err := BuildCharSets(opts)
	if err != nil {
		return "", err
	}

	used := make(map[rune]bool)
	password := make([]rune, 0, opts.Length)

	addFromSet := func(set string, count int) error {
		for i := 0; i < count; i++ {
			r, err := pickFromSet(src, set, used, opts.PreventRepeated)
			if err != nil {
				return err
			}
			password = append(password, r)
			used[r] = true
		}
		return nil
	}

	if opts.Uppercase {
		if err := addFromSet(cs.Uppercase, opts.MinUppercase); err != nil {
			return "", err
		}
	}
	if opts.Lowercase {
		if err := addFromSet(cs.Lowercase, opts.MinLowercase); err != nil {
			return "", err
		}
	}
	if opts.Numbers {
		if err := addFromSet(cs.Numbers, opts.MinNumbers); err != nil {
			return "", err
		}
	}
	if opts.Symbols {
		if err := addFromSet(cs.Symbols, opts.MinSymbols); err != nil {
			return "", err
		}
	}

	for len(password) < opts.Length {
		r, err := pickFromSet(src, cs.Combined, used, opts.PreventRepeated)
		if err != nil {
			return "", err
		}
		password = append(password, r)
		used[r] = true
	}

	if err := random.Shuffle(src, password); err != nil {
		return "", err
	}

	return string(password), nil
}

// UserMessage converts generator errors to human-readable messages.
func UserMessage(err error) string {
	switch err {
	case ErrInvalidLength:
		return "Password length must be between 4 and 128 characters."
	case ErrNoCharacterSet:
		return "Please select at least one character type."
	case ErrImpossibleMinimums:
		return "Minimum character requirements exceed password length."
	case ErrNoAvailableChars:
		return "Your excluded characters remove all available characters from a selected group."
	case ErrInvalidPINLength:
		return "PIN length must be between 4 and 32 digits."
	case ErrInvalidWordCount:
		return "Passphrase must contain between 3 and 20 words."
	case random.ErrRandomSourceFailure:
		return "Secure random source unavailable. Cannot generate passwords safely."
	default:
		if err != nil {
			return err.Error()
		}
		return ""
	}
}

// CountRequirements verifies a password meets minimum requirements.
func CountRequirements(password string, opts PasswordOptions, cs CharSets) (upper, lower, nums, syms int) {
	upper = countChars(password, cs.Uppercase)
	lower = countChars(password, cs.Lowercase)
	nums = countChars(password, cs.Numbers)
	syms = countChars(password, cs.Symbols)
	return
}

func pickFromSet(src random.Source, set string, used map[rune]bool, preventRepeated bool) (rune, error) {
	if set == "" {
		return 0, ErrNoAvailableChars
	}
	runes := []rune(set)
	if preventRepeated {
		var available []rune
		for _, r := range runes {
			if !used[r] {
				available = append(available, r)
			}
		}
		if len(available) == 0 {
			return 0, ErrNoAvailableChars
		}
		runes = available
	}
	idx, err := src.RandomInt(len(runes))
	if err != nil {
		return 0, err
	}
	return runes[idx], nil
}

// IsAmbiguousChar reports whether a rune is commonly confused with others.
func IsAmbiguousChar(r rune) bool {
	return ambiguousChars[r]
}

// SymbolSet returns the default symbol character set.
func SymbolSet() string {
	return SymbolChars
}

// ContainsOnlyDigits reports whether s contains only numeric characters.
func ContainsOnlyDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// NormalizeExcludedChars removes duplicates from excluded character input.
func NormalizeExcludedChars(s string) string {
	seen := make(map[rune]bool)
	var out strings.Builder
	for _, r := range s {
		if !seen[r] {
			seen[r] = true
			out.WriteRune(r)
		}
	}
	return out.String()
}

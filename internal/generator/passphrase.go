package generator

import (
	"strings"

	"github.com/rahul1534/PassGen/internal/random"
	"github.com/rahul1534/PassGen/wordlist"
)

// PassphraseOptions configures passphrase generation.
type PassphraseOptions struct {
	Words      int
	Separator  string
	Capitalize bool
	AddNumber  bool
	AddSymbol  bool
}

// DefaultPassphraseOptions returns recommended passphrase defaults.
func DefaultPassphraseOptions() PassphraseOptions {
	return PassphraseOptions{
		Words:      5,
		Separator:  "-",
		Capitalize: false,
		AddNumber:  true,
		AddSymbol:  false,
	}
}

var passphraseWords []string

// LoadWordList loads the bundled word list. Safe to call multiple times.
func LoadWordList() ([]string, error) {
	if len(passphraseWords) > 0 {
		return passphraseWords, nil
	}
	words, err := wordlist.Words()
	if err != nil {
		return nil, err
	}
	if len(words) < 2000 {
		return nil, ErrNoAvailableChars
	}
	passphraseWords = words
	return passphraseWords, nil
}

// SetWordListForTests replaces the word list in tests.
func SetWordListForTests(words []string) {
	passphraseWords = words
}

// WordListSize returns the number of loaded words.
func WordListSize() int {
	return len(passphraseWords)
}

// ValidatePassphraseOptions checks passphrase configuration.
func ValidatePassphraseOptions(opts PassphraseOptions) error {
	if opts.Words < MinPassphraseWords || opts.Words > MaxPassphraseWords {
		return ErrInvalidWordCount
	}
	if _, err := LoadWordList(); err != nil {
		return err
	}
	return nil
}

// GeneratePassphrase creates a random passphrase from the bundled word list.
func GeneratePassphrase(src random.Source, opts PassphraseOptions) (string, error) {
	if err := ValidatePassphraseOptions(opts); err != nil {
		return "", err
	}

	words, err := LoadWordList()
	if err != nil {
		return "", err
	}

	selected := make([]string, opts.Words)
	for i := 0; i < opts.Words; i++ {
		idx, err := src.RandomInt(len(words))
		if err != nil {
			return "", err
		}
		word := words[idx]
		if opts.Capitalize {
			word = strings.ToUpper(word[:1]) + word[1:]
		}
		selected[i] = word
	}

	result := strings.Join(selected, opts.Separator)

	if opts.AddNumber {
		numIdx, err := src.RandomInt(10)
		if err != nil {
			return "", err
		}
		result += string(NumberChars[numIdx])
	}

	if opts.AddSymbol {
		syms := []rune(SymbolChars)
		symIdx, err := src.RandomInt(len(syms))
		if err != nil {
			return "", err
		}
		result += string(syms[symIdx])
	}

	return result, nil
}

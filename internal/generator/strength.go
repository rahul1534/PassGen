package generator

import "math"

// StrengthLevel represents estimated password strength.
type StrengthLevel int

const (
	StrengthVeryWeak StrengthLevel = iota
	StrengthWeak
	StrengthFair
	StrengthStrong
	StrengthVeryStrong
)

func (s StrengthLevel) String() string {
	switch s {
	case StrengthVeryWeak:
		return "Very Weak"
	case StrengthWeak:
		return "Weak"
	case StrengthFair:
		return "Fair"
	case StrengthStrong:
		return "Strong"
	case StrengthVeryStrong:
		return "Very Strong"
	default:
		return "Unknown"
	}
}

// StrengthBarWidth returns a percentage width for the strength indicator (0-100).
func (s StrengthLevel) BarWidth() int {
	switch s {
	case StrengthVeryWeak:
		return 20
	case StrengthWeak:
		return 40
	case StrengthFair:
		return 60
	case StrengthStrong:
		return 80
	case StrengthVeryStrong:
		return 100
	default:
		return 0
	}
}

// StrengthResult holds entropy estimate and display level.
type StrengthResult struct {
	Level   StrengthLevel
	Entropy float64
}

// EstimatePasswordStrength calculates strength from password options and result.
func EstimatePasswordStrength(password string, opts PasswordOptions) StrengthResult {
	cs, err := BuildCharSets(opts)
	if err != nil {
		return StrengthFromEntropy(0)
	}
	poolSize := len(cs.Combined)
	if poolSize == 0 {
		return StrengthFromEntropy(0)
	}
	entropy := float64(len(password)) * math.Log2(float64(poolSize))
	return StrengthFromEntropy(entropy)
}

// EstimatePassphraseStrength calculates strength from word count and list size.
func EstimatePassphraseStrength(opts PassphraseOptions, wordListSize int) StrengthResult {
	if wordListSize <= 0 {
		wordListSize = WordListSize()
	}
	if wordListSize <= 0 {
		return StrengthFromEntropy(0)
	}
	entropy := float64(opts.Words) * math.Log2(float64(wordListSize))
	if opts.AddNumber {
		entropy += math.Log2(10)
	}
	if opts.AddSymbol {
		entropy += math.Log2(float64(len(SymbolChars)))
	}
	return StrengthFromEntropy(entropy)
}

// EstimatePINStrength calculates strength for numeric PINs.
func EstimatePINStrength(length int, uniqueDigits bool) StrengthResult {
	pool := 10
	if uniqueDigits && length <= pool {
		// Permutation: P(10, length)
		entropy := 0.0
		for i := 0; i < length; i++ {
			entropy += math.Log2(float64(pool - i))
		}
		return StrengthFromEntropy(entropy)
	}
	entropy := float64(length) * math.Log2(10)
	return StrengthFromEntropy(entropy)
}

// StrengthFromEntropy maps bits of entropy to a display level.
func StrengthFromEntropy(entropy float64) StrengthResult {
	var level StrengthLevel
	switch {
	case entropy < 28:
		level = StrengthVeryWeak
	case entropy < 36:
		level = StrengthWeak
	case entropy < 60:
		level = StrengthFair
	case entropy < 80:
		level = StrengthStrong
	default:
		level = StrengthVeryStrong
	}
	return StrengthResult{Level: level, Entropy: entropy}
}

package generator

import (
	"github.com/rahul1534/PassGen/internal/random"
)

// PINOptions configures numeric PIN generation.
type PINOptions struct {
	Length                 int
	AllowRepeatedDigits    bool
	AvoidAmbiguousPatterns bool
}

// DefaultPINOptions returns recommended PIN defaults.
func DefaultPINOptions() PINOptions {
	return PINOptions{
		Length:                 6,
		AllowRepeatedDigits:    true,
		AvoidAmbiguousPatterns: false,
	}
}

// ValidatePINOptions checks PIN configuration.
func ValidatePINOptions(opts PINOptions) error {
	if opts.Length < MinPINLength || opts.Length > MaxPINLength {
		return ErrInvalidPINLength
	}
	return nil
}

// GeneratePIN creates a numeric PIN.
func GeneratePIN(src random.Source, opts PINOptions) (string, error) {
	if err := ValidatePINOptions(opts); err != nil {
		return "", err
	}

	const maxAttempts = 100
	for attempt := 0; attempt < maxAttempts; attempt++ {
		pin, err := generatePINOnce(src, opts)
		if err != nil {
			return "", err
		}
		if !opts.AvoidAmbiguousPatterns || !hasAmbiguousPINPattern(pin) {
			return pin, nil
		}
	}
	return generatePINOnce(src, opts)
}

func generatePINOnce(src random.Source, opts PINOptions) (string, error) {
	digits := make([]byte, opts.Length)
	used := make(map[byte]bool)

	for i := 0; i < opts.Length; i++ {
		var available []byte
		for d := byte('0'); d <= '9'; d++ {
			if !opts.AllowRepeatedDigits && used[d] {
				continue
			}
			available = append(available, d)
		}
		if len(available) == 0 {
			return "", ErrImpossibleMinimums
		}
		idx, err := src.RandomInt(len(available))
		if err != nil {
			return "", err
		}
		digits[i] = available[idx]
		if !opts.AllowRepeatedDigits {
			used[available[idx]] = true
		}
	}
	return string(digits), nil
}

func hasAmbiguousPINPattern(pin string) bool {
	if len(pin) < 3 {
		return false
	}
	// Sequential ascending or descending (e.g. 123456, 654321)
	ascending := true
	descending := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			ascending = false
		}
		if pin[i] != pin[i-1]-1 {
			descending = false
		}
	}
	if ascending || descending {
		return true
	}
	// All same digit
	allSame := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[0] {
			allSame = false
			break
		}
	}
	return allSame
}

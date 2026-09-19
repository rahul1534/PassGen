package generator

import "testing"

func TestValidateRejectsNegativeMinimums(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*PasswordOptions)
	}{
		{"negative upper", func(o *PasswordOptions) { o.MinUppercase = -1 }},
		{"negative lower", func(o *PasswordOptions) { o.MinLowercase = -1 }},
		{"negative numbers", func(o *PasswordOptions) { o.MinNumbers = -1 }},
		{"negative symbols", func(o *PasswordOptions) { o.MinSymbols = -1 }},
		// A negative minimum used to cancel out an oversized one in the
		// minTotal > Length check, so validation passed.
		{"negative offsets oversized minimum", func(o *PasswordOptions) {
			o.Length = 10
			o.MinUppercase = 30
			o.MinLowercase = -20
			o.MinNumbers, o.MinSymbols = 0, 0
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := DefaultPasswordOptions()
			tc.mut(&opts)
			if err := ValidatePasswordOptions(opts); err != ErrInvalidMinimum {
				t.Fatalf("expected ErrInvalidMinimum, got %v", err)
			}
			pw, err := GeneratePassword(propertySource(1), opts)
			if err == nil || pw != "" {
				t.Fatalf("expected error and empty output, got %q, %v", pw, err)
			}
		})
	}
}

func TestNegativeMinimumOnDisabledClassIsIgnored(t *testing.T) {
	opts := DefaultPasswordOptions()
	opts.Symbols = false
	opts.MinSymbols = -3 // class not in use; must not matter
	if err := ValidatePasswordOptions(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Whatever the minimums are, a successful result must be exactly the requested
// length and satisfy every requested minimum.
func TestPropertyOutputLengthMatchesRequest(t *testing.T) {
	seed := 0
	for length := 4; length <= 40; length += 6 {
		for _, minU := range []int{-5, 0, 1, 12, 45} {
			for _, minL := range []int{-20, 0, 1, 12} {
				for _, minN := range []int{-1, 0, 3} {
					opts := DefaultPasswordOptions()
					opts.Length, opts.MinUppercase, opts.MinLowercase, opts.MinNumbers = length, minU, minL, minN
					seed++
					pw, err := GeneratePassword(propertySource(seed), opts)
					if err != nil {
						continue
					}
					if len([]rune(pw)) != length {
						t.Fatalf("length=%d minU=%d minL=%d minN=%d: got %d chars", length, minU, minL, minN, len([]rune(pw)))
					}
					cs, _ := BuildCharSets(opts)
					u, l, n, _ := CountRequirements(pw, opts, cs)
					if u < minU || l < minL || n < minN {
						t.Fatalf("minimums not met for length=%d: u=%d l=%d n=%d", length, u, l, n)
					}
				}
			}
		}
	}
}

func TestUserMessageInvalidMinimum(t *testing.T) {
	if got := UserMessage(ErrInvalidMinimum); got == "" || got == ErrInvalidMinimum.Error() {
		t.Fatalf("expected a friendly message, got %q", got)
	}
}

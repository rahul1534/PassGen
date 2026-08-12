package generator

const (
	LowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	UppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NumberChars    = "0123456789"
	SymbolChars    = "!@#$%^&*()-_=+[]{};:,.?"
)

var ambiguousChars = map[rune]bool{
	'0': true, 'O': true, 'o': true,
	'1': true, 'I': true, 'l': true,
}

// CharSets holds filtered character sets for password generation.
type CharSets struct {
	Lowercase string
	Uppercase string
	Numbers   string
	Symbols   string
	Combined  string
}

// BuildCharSets constructs available character sets from options.
func BuildCharSets(opts PasswordOptions) (CharSets, error) {
	excluded := make(map[rune]bool)
	for _, r := range opts.ExcludedCharacters {
		excluded[r] = true
	}

	filter := func(set string) string {
		var out []rune
		for _, r := range set {
			if excluded[r] {
				continue
			}
			if opts.ExcludeAmbiguous && ambiguousChars[r] {
				continue
			}
			out = append(out, r)
		}
		return string(out)
	}

	cs := CharSets{}
	if opts.Lowercase {
		cs.Lowercase = filter(LowercaseChars)
	}
	if opts.Uppercase {
		cs.Uppercase = filter(UppercaseChars)
	}
	if opts.Numbers {
		cs.Numbers = filter(NumberChars)
	}
	if opts.Symbols {
		cs.Symbols = filter(SymbolChars)
	}

	if cs.Lowercase == "" && cs.Uppercase == "" && cs.Numbers == "" && cs.Symbols == "" {
		return cs, ErrNoCharacterSet
	}

	cs.Combined = cs.Lowercase + cs.Uppercase + cs.Numbers + cs.Symbols
	if cs.Combined == "" {
		return cs, ErrNoAvailableChars
	}

	return cs, nil
}

// ValidatePasswordOptions checks configuration before generation.
func ValidatePasswordOptions(opts PasswordOptions) error {
	if opts.Length < MinPasswordLength || opts.Length > MaxPasswordLength {
		return ErrInvalidLength
	}

	cs, err := BuildCharSets(opts)
	if err != nil {
		return err
	}

	if opts.Lowercase && cs.Lowercase == "" {
		return ErrNoAvailableChars
	}
	if opts.Uppercase && cs.Uppercase == "" {
		return ErrNoAvailableChars
	}
	if opts.Numbers && cs.Numbers == "" {
		return ErrNoAvailableChars
	}
	if opts.Symbols && cs.Symbols == "" {
		return ErrNoAvailableChars
	}

	minTotal := 0
	if opts.Lowercase {
		minTotal += opts.MinLowercase
	}
	if opts.Uppercase {
		minTotal += opts.MinUppercase
	}
	if opts.Numbers {
		minTotal += opts.MinNumbers
	}
	if opts.Symbols {
		minTotal += opts.MinSymbols
	}

	if minTotal > opts.Length {
		return ErrImpossibleMinimums
	}

	if opts.PreventRepeated && len(cs.Combined) < opts.Length {
		return ErrImpossibleMinimums
	}

	return nil
}

func countChars(password string, set string) int {
	setMap := make(map[rune]bool)
	for _, r := range set {
		setMap[r] = true
	}
	count := 0
	for _, r := range password {
		if setMap[r] {
			count++
		}
	}
	return count
}

func charInSet(r rune, set string) bool {
	for _, c := range set {
		if c == r {
			return true
		}
	}
	return false
}

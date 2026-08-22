package generator

import "errors"

var (
	ErrInvalidLength      = errors.New("invalid password length")
	ErrNoCharacterSet     = errors.New("no character set selected")
	ErrImpossibleMinimums = errors.New("minimum requirements exceed length")
	ErrNoAvailableChars   = errors.New("no available characters")
	ErrInvalidPINLength   = errors.New("invalid PIN length")
	ErrInvalidWordCount   = errors.New("invalid word count")
	ErrUnableToSatisfyPINConstraints = errors.New("unable to satisfy PIN constraints")
)

const (
	MinPasswordLength  = 4
	MaxPasswordLength  = 128
	MinPINLength       = 4
	MaxPINLength       = 32
	MinPassphraseWords = 3
	MaxPassphraseWords = 20
)

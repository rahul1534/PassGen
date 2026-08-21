package random

import "errors"

var ErrRandomSourceFailure = errors.New("secure random source unavailable")

// Source provides cryptographically secure random integers and shuffling.
// Production implementations should use a crypto/rand-backed Source.
type Source interface {
	RandomInt(max int) (int, error)
}

// Shuffle randomly permutes items using the provided source.
func Shuffle[T any](src Source, items []T) error {
	n := len(items)
	for i := n - 1; i > 0; i-- {
		j, err := src.RandomInt(i + 1)
		if err != nil {
			return err
		}
		items[i], items[j] = items[j], items[i]
	}
	return nil
}

// RandomIndex selects a uniform random index in [0, len(items)).
func RandomIndex(src Source, items []string) (int, error) {
	if len(items) == 0 {
		return 0, errors.New("empty items")
	}
	return src.RandomInt(len(items))
}

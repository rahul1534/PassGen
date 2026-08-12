package random

// DeterministicSource provides predictable randomness for tests only.
type DeterministicSource struct {
	values []int
	index  int
}

// NewDeterministicSource returns a test random source with the given values.
func NewDeterministicSource(values ...int) *DeterministicSource {
	return &DeterministicSource{values: values}
}

func (d *DeterministicSource) RandomInt(max int) (int, error) {
	if max <= 0 {
		return 0, ErrRandomSourceFailure
	}
	if len(d.values) == 0 {
		return 0, nil
	}
	v := d.values[d.index%len(d.values)]
	d.index++
	return v % max, nil
}

// Reset rewinds the deterministic sequence.
func (d *DeterministicSource) Reset() {
	d.index = 0
}

// SetValues replaces the sequence of values.
func (d *DeterministicSource) SetValues(values ...int) {
	d.values = values
	d.index = 0
}

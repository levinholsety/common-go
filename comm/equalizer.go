package comm

// Equalizer provides Equal method.
type Equalizer interface {
	// Equal returns true if self and the argument are equal.
	Equal(a Equalizer) bool
}

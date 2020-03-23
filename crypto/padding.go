package crypto

// Padding is the interface that wraps the padding methods.
type Padding interface {
	AddPadding(block []byte, size int)
	RemovePadding(block []byte) ([]byte, error)
}

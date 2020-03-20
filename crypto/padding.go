package crypto

// Padding is the interface that wraps the padding methods.
type Padding interface {
	AddPadding(data []byte, blockSize int) []byte
	RemovePadding(data []byte, blockSize int) (int, error)
}

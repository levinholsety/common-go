package crypto

// Padding provides methods for padding.
type Padding interface {
	AddPadding(data []byte, blockSize int) []byte
	RemovePadding(data []byte, blockSize int) (int, error)
}

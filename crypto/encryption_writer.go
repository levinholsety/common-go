package crypto

import (
	"io"
)

// NewEncryptionWriter wraps w and returns an encryption writer to encrypt data.
// w holds the data to be encrypted.
// The data will be encrypted after it has been written to the encryption writer.
// Remember to close the encryption writer at the end.
func NewEncryptionWriter(w io.Writer, encryptor Encryptor, paddingAlg PaddingAlgorithm) io.WriteCloser {
	return &blockWriter{
		writer: &blockEncryptionWriter{
			writer:      w,
			encryptor:   encryptor,
			cipherBlock: make([]byte, encryptor.CipherBlockSize()),
		},
		block:      make([]byte, encryptor.BlockSize()),
		paddingAlg: paddingAlg,
	}
}

type blockWriter struct {
	writer     io.Writer
	closed     bool
	block      []byte
	offset     int
	paddingAlg PaddingAlgorithm
}

func (w *blockWriter) Write(p []byte) (n int, err error) {
	if w.closed {
		return
	}
	blockSize := len(w.block)
	dataLen := len(p)
	for n < dataLen {
		count := copy(w.block[w.offset:], p[n:])
		n += count
		w.offset += count
		if w.offset == blockSize {
			_, err = w.writer.Write(w.block)
			if err != nil {
				return
			}
			w.offset = 0
		}
	}
	return
}

func (w *blockWriter) Close() (err error) {
	if w.closed {
		return
	}
	w.closed = true
	if w.paddingAlg == nil {
		w.block = w.block[:w.offset]
	} else {
		err = w.paddingAlg.AddPadding(w.block, w.offset)
		if err != nil {
			return
		}
	}
	_, err = w.writer.Write(w.block)
	return
}

type blockEncryptionWriter struct {
	writer      io.Writer
	encryptor   Encryptor
	cipherBlock []byte
}

func (w *blockEncryptionWriter) Write(p []byte) (n int, err error) {
	err = w.encryptor.Encrypt(w.cipherBlock, p)
	if err != nil {
		return
	}
	_, err = w.writer.Write(w.cipherBlock)
	return
}

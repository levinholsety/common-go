package crypto

import (
	"bufio"
	"io"
)

// NewDecryptionReader wraps r and returns a decryption reader to decrypt data.
// r holds the data to be decrypted.
// The data will be decrypted after it has been read from the decryption reader.
func NewDecryptionReader(r io.Reader, decryptor Decryptor, paddingAlg PaddingAlgorithm) (result io.Reader, err error) {
	rd, err := newBlockDecryptionReader(r, decryptor, paddingAlg)
	if err != nil {
		return
	}
	result = &decryptionReader{
		reader: rd,
	}
	return
}

type decryptionReader struct {
	reader *blockDecryptionReader
}

func (r *decryptionReader) Read(p []byte) (n int, err error) {
	n, err = io.ReadFull(r.reader, p)
	if err == io.ErrUnexpectedEOF {
		err = io.EOF
	}
	return
}

func newBlockDecryptionReader(r io.Reader, decryptor Decryptor, paddingAlg PaddingAlgorithm) (result *blockDecryptionReader, err error) {
	if _, ok := r.(*bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	reader := &blockDecryptionReader{
		reader:      r,
		decryptor:   decryptor,
		cipherBlock: make([]byte, decryptor.CipherBlockSize()),
		block:       make([]byte, decryptor.BlockSize()),
		paddingAlg:  paddingAlg,
	}
	err = reader.readCipherBlock()
	if err != nil {
		return
	}
	result = reader
	return
}

type blockDecryptionReader struct {
	reader      io.Reader
	eof         bool
	cipherBlock []byte
	block       []byte
	decryptor   Decryptor
	paddingAlg  PaddingAlgorithm
}

func (r *blockDecryptionReader) Read(p []byte) (n int, err error) {
	if r.eof {
		err = io.EOF
		return
	}
	n, err = r.decryptor.Decrypt(r.block, r.cipherBlock)
	if err != nil {
		return
	}
	err = r.readCipherBlock()
	if err != nil {
		return
	}
	if !r.eof || r.paddingAlg == nil {
		n = copy(p, r.block[:n])
		return
	}
	r.block, err = r.paddingAlg.RemovePadding(r.block)
	if err != nil {
		return
	}
	n = copy(p, r.block)
	return
}

func (r *blockDecryptionReader) readCipherBlock() (err error) {
	_, err = io.ReadFull(r.reader, r.cipherBlock)
	if err == io.EOF {
		err = nil
		r.eof = true
	} else if err == io.ErrUnexpectedEOF {
		err = ErrBadPadding
		r.eof = true
	}
	return
}

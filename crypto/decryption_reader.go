package crypto

import (
	"bufio"
	"io"
)

type decryptionReader struct {
	reader          io.Reader
	decryptor       Decryptor
	padding         Padding
	cipherBlockSize int
	cipherBlock     []byte
	eof             bool
}

func (p *decryptionReader) Read(buffer []byte) (n int, err error) {
	if p.eof {
		err = io.EOF
		return
	}
	n, err = p.decryptor.Decrypt(buffer, p.cipherBlock)
	if _, err = io.ReadFull(p.reader, p.cipherBlock); err != nil {
		if err == io.EOF {
			p.eof = true
			if p.padding != nil {
				if buffer, err = p.padding.RemovePadding(buffer[:n]); err != nil {
					return
				}
				n = len(buffer)
			}
		}
		return
	}
	return
}

// NewDecryptionReader wraps r and returns a decryption reader to decrypt data.
// r holds the data to be decrypted.
// The data will be decrypted after it has been read from the decryption reader.
func NewDecryptionReader(r io.Reader, decryptor Decryptor, padding Padding) (dr io.Reader, err error) {
	cipherBlockSize := decryptor.CipherBlockSize()
	cipherBlock := make([]byte, cipherBlockSize)
	eof := false
	if _, err = io.ReadFull(r, cipherBlock); err != nil {
		if err == io.EOF {
			eof = true
			err = nil
		} else {
			return
		}
	}
	dr = bufio.NewReaderSize(&decryptionReader{
		reader:          r,
		decryptor:       decryptor,
		padding:         padding,
		cipherBlockSize: cipherBlockSize,
		cipherBlock:     cipherBlock,
		eof:             eof,
	}, cipherBlockSize)
	return
}

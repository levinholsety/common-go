package crypto

import (
	"bufio"
	"crypto/cipher"
	"io"
)

// CipherWriter represents a cipher writer.
type CipherWriter struct {
	writer  io.Writer
	cb      cipher.Block
	dst     []byte
	buf     []byte
	bufLen  int
	padding Padding
	closed  bool
}

// NewCipherWriter creates a cipher writer.
func NewCipherWriter(w io.Writer, cipherBlock cipher.Block, padding Padding) io.WriteCloser {
	return &CipherWriter{
		writer:  w,
		cb:      cipherBlock,
		dst:     make([]byte, cipherBlock.BlockSize()),
		buf:     make([]byte, cipherBlock.BlockSize()),
		bufLen:  0,
		padding: padding,
		closed:  false,
	}
}

func (w *CipherWriter) Write(p []byte) (n int, err error) {
	var size int
	if w.bufLen > 0 {
		size = copy(w.buf[w.bufLen:], p)
		w.bufLen += size
		if w.bufLen < w.cb.BlockSize() {
			return
		}
		p = p[size:]
		w.cb.Encrypt(w.dst, w.buf)
		w.bufLen = 0
		size, err = w.writer.Write(w.dst)
		if err != nil {
			return
		}
		n += size
	}
	dataLength := len(p)
	count := dataLength / w.cb.BlockSize()
	for i := 0; i < count; i++ {
		w.cb.Encrypt(w.dst, p[i*w.cb.BlockSize():i*w.cb.BlockSize()+w.cb.BlockSize()])
		size, err = w.writer.Write(w.dst)
		if err != nil {
			return
		}
		n += size
	}
	w.bufLen = dataLength % w.cb.BlockSize()
	if w.bufLen > 0 {
		copy(w.buf, p[dataLength-w.bufLen:])
	}
	return
}

// Close writes last block with padding to destination.
func (w *CipherWriter) Close() (err error) {
	if w.closed {
		return
	}
	w.cb.Encrypt(w.dst, w.padding.AddPadding(w.buf[:w.bufLen], w.cb.BlockSize()))
	_, err = w.writer.Write(w.dst)
	if err != nil {
		return
	}
	w.closed = true
	return
}

// CipherReader represents a cipher reader.
type CipherReader struct {
	reader  io.Reader
	cb      cipher.Block
	src     []byte
	dst     []byte
	buf     []byte
	bufLen  int
	bufOff  int
	padding Padding
	eof     bool
}

// NewCipherReader creates a cipher reader.
func NewCipherReader(r io.Reader, cipherBlock cipher.Block, padding Padding) io.Reader {
	if _, ok := r.(interface{}).(bufio.Reader); !ok {
		r = bufio.NewReader(r)
	}
	return &CipherReader{
		reader:  r,
		cb:      cipherBlock,
		src:     make([]byte, cipherBlock.BlockSize()),
		dst:     nil,
		buf:     make([]byte, cipherBlock.BlockSize()),
		bufLen:  0,
		bufOff:  0,
		padding: padding,
		eof:     false,
	}
}

func (r *CipherReader) Read(p []byte) (n int, err error) {
	if r.bufLen-r.bufOff > 0 {
		n = copy(p, r.buf[r.bufOff:r.bufLen])
		r.bufOff += n
	}
	for n < len(p) {
		r.bufLen, err = r.readBlock(r.buf)
		if err != nil && err != io.EOF {
			return
		}
		if r.bufLen == 0 && err == io.EOF {
			return
		}
		r.bufOff = copy(p[n:], r.buf[:r.bufLen])
		n += r.bufOff
	}
	return
}

func (r *CipherReader) readBlock(p []byte) (n int, err error) {
	if r.eof {
		return 0, io.EOF
	}
	if r.dst == nil {
		n, err = r.reader.Read(r.src)
		if err != nil && err != io.EOF {
			return
		}
		if n == 0 && err == io.EOF {
			r.eof = true
			return
		}
		r.dst = make([]byte, r.cb.BlockSize())
		r.cb.Decrypt(r.dst, r.src)
	}
	size, err := r.reader.Read(r.src)
	if err != nil && err != io.EOF {
		return
	}
	if size == 0 && err == io.EOF {
		r.eof = true
		n = copy(p, r.padding.RemovePadding(r.dst))
		return
	}
	n = copy(p, r.dst)
	r.cb.Decrypt(r.dst, r.src[:size])
	return
}

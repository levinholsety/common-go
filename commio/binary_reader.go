package commio

import (
	"encoding/binary"
	"io"
)

// NewBinaryReader creates an instance of BinaryReader.
func NewBinaryReader(r io.Reader, order binary.ByteOrder) *BinaryReader {
	br := &BinaryReader{
		reader:    r,
		ByteOrder: order,
	}
	if seeker, ok := r.(io.Seeker); ok {
		br.seeker = seeker
		br.offset, _ = br.SeekCurrent(0)
	}
	return br
}

// BinaryReader provides methods to read data.
type BinaryReader struct {
	reader    io.Reader
	seeker    io.Seeker
	ByteOrder binary.ByteOrder
	offset    int64
}

// Reader returns internal reader.
func (p *BinaryReader) Reader() io.Reader {
	return p.reader
}

// SeekStart seeks from the start of the file.
func (p *BinaryReader) SeekStart(offset int64) error {
	_, err := p.seeker.Seek(offset, io.SeekStart)
	return err
}

// SeekCurrent seeks from the current offset.
func (p *BinaryReader) SeekCurrent(offset int64) (int64, error) {
	n, err := p.seeker.Seek(offset, io.SeekCurrent)
	return n, err
}

// SeekEnd seeks from the end of the file.
func (p *BinaryReader) SeekEnd(offset int64) (int64, error) {
	n, err := p.seeker.Seek(offset, io.SeekEnd)
	return n, err
}

// Offset returns offset from beginning.
func (p *BinaryReader) Offset() (n int64, err error) {
	if n, err = p.SeekCurrent(0); err != nil {
		return
	}
	n -= p.offset
	return
}

// Length returns total byte count can be read.
func (p *BinaryReader) Length() (n int64, err error) {
	current, err := p.SeekCurrent(0)
	if err != nil {
		return
	}
	if n, err = p.SeekEnd(0); err != nil {
		return
	}
	n -= p.offset
	err = p.SeekStart(current)
	return
}

// Read reads content into data.
func (p *BinaryReader) Read(data interface{}) error {
	return binary.Read(p.reader, p.ByteOrder, data)
}

// ReadUInt64 reads a uint64 value.
func (p *BinaryReader) ReadUInt64() (v uint64, err error) {
	err = p.Read(&v)
	return
}

// ReadInt64 reads a int64 value.
func (p *BinaryReader) ReadInt64() (v int64, err error) {
	err = p.Read(&v)
	return
}

// ReadUInt32 reads a uint32 value.
func (p *BinaryReader) ReadUInt32() (v uint32, err error) {
	err = p.Read(&v)
	return
}

// ReadInt32 reads a int32 value.
func (p *BinaryReader) ReadInt32() (v int32, err error) {
	err = p.Read(&v)
	return
}

// ReadUInt16 reads a uint16 value.
func (p *BinaryReader) ReadUInt16() (v uint16, err error) {
	err = p.Read(&v)
	return
}

// ReadInt16 reads a int16 value.
func (p *BinaryReader) ReadInt16() (v int16, err error) {
	err = p.Read(&v)
	return
}

// ReadUInt8 reads a uint8 value.
func (p *BinaryReader) ReadUInt8() (v uint8, err error) {
	err = p.Read(&v)
	return
}

// ReadInt8 reads a int8 value.
func (p *BinaryReader) ReadInt8() (v int8, err error) {
	err = p.Read(&v)
	return
}

// ReadByteArray reads byte array in specified size.
func (p *BinaryReader) ReadByteArray(size int) (v []byte, err error) {
	v = make([]byte, size)
	err = p.Read(&v)
	return
}

// ReadUInt reads byte array in specified size and returns uint.
func (p *BinaryReader) ReadUInt(size int) (v uint, err error) {
	buf, err := p.ReadByteArray(size)
	if err != nil {
		return
	}
	for i, b := range buf {
		v |= (uint(b) << ((size - i - 1) * 8))
	}
	return
}

// ReadString reads string in specified size.
func (p *BinaryReader) ReadString(size int) (v string, err error) {
	buf, err := p.ReadByteArray(size)
	if err != nil {
		return
	}
	v = string(buf)
	return
}

package commio

import (
	"encoding/binary"
	"io"
)

// NewBinaryReader creates an instance of BinaryReader.
func NewBinaryReader(r io.ReadSeeker, order binary.ByteOrder) *BinaryReader {
	br := &BinaryReader{
		reader: r,
		order:  order,
	}
	br.offset = br.SeekCurrent(0)
	return br
}

// BinaryReader provides methods to read data.
type BinaryReader struct {
	reader io.ReadSeeker
	order  binary.ByteOrder
	offset int64
}

// SeekStart seeks from the start of the file.
func (p *BinaryReader) SeekStart(offset int64) int64 {
	return p.seek(offset, io.SeekStart)
}

// SeekCurrent seeks from the current offset.
func (p *BinaryReader) SeekCurrent(offset int64) int64 {
	return p.seek(offset, io.SeekCurrent)
}

// SeekEnd seeks from the end of the file.
func (p *BinaryReader) SeekEnd(offset int64) int64 {
	return p.seek(offset, io.SeekEnd)
}

func (p *BinaryReader) seek(offset int64, whence int) int64 {
	n, err := p.reader.Seek(offset, whence)
	if err != nil {
		panic(err)
	}
	return n
}

// Offset returns offset from beginning.
func (p *BinaryReader) Offset() int64 {
	return p.SeekCurrent(0) - p.offset
}

// Length returns total byte count can be read.
func (p *BinaryReader) Length() int64 {
	current := p.SeekCurrent(0)
	length := p.SeekEnd(0) - p.offset
	p.SeekStart(current)
	return length
}

// Read reads content into data.
func (p *BinaryReader) Read(data interface{}) {
	if err := binary.Read(p.reader, p.order, data); err != nil {
		panic(err)
	}
}

// ReadUInt64 reads a uint64 value.
func (p *BinaryReader) ReadUInt64() (v uint64) {
	p.Read(&v)
	return
}

// ReadInt64 reads a int64 value.
func (p *BinaryReader) ReadInt64() (v int64) {
	p.Read(&v)
	return
}

// ReadUInt32 reads a uint32 value.
func (p *BinaryReader) ReadUInt32() (v uint32) {
	p.Read(&v)
	return
}

// ReadInt32 reads a int32 value.
func (p *BinaryReader) ReadInt32() (v int32) {
	p.Read(&v)
	return
}

// ReadUInt16 reads a uint16 value.
func (p *BinaryReader) ReadUInt16() (v uint16) {
	p.Read(&v)
	return
}

// ReadInt16 reads a int16 value.
func (p *BinaryReader) ReadInt16() (v int16) {
	p.Read(&v)
	return
}

// ReadUInt8 reads a uint8 value.
func (p *BinaryReader) ReadUInt8() (v uint8) {
	p.Read(&v)
	return
}

// ReadInt8 reads a int8 value.
func (p *BinaryReader) ReadInt8() (v int8) {
	p.Read(&v)
	return
}

// ReadByteArray reads byte array in specified size.
func (p *BinaryReader) ReadByteArray(size int) (v []byte) {
	v = make([]byte, size)
	p.Read(&v)
	return
}

// ReadUInt reads byte array in specified size and returns uint.
func (p *BinaryReader) ReadUInt(size int) (v uint) {
	buf := p.ReadByteArray(size)
	for i, b := range buf {
		v |= (uint(b) << ((size - i - 1) * 8))
	}
	return
}

// ReadString reads string in specified size.
func (p *BinaryReader) ReadString(size int) (v string) {
	v = string(p.ReadByteArray(size))
	return
}

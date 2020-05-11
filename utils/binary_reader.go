package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

// NewBinaryReader creates an instance of BinaryReader.
func NewBinaryReader(r io.Reader, o binary.ByteOrder) *BinaryReader {
	return &BinaryReader{
		reader:    r,
		ByteOrder: o,
		OnError:   func(err error) { panic(err) },
	}
}

// BinaryReader provides methods to read data.
type BinaryReader struct {
	reader    io.Reader
	ByteOrder binary.ByteOrder
	OnError   func(err error)
}

// Read reads structured binary data from internal reader into data.
func (p *BinaryReader) Read(data interface{}) (err error) {
	if buf, ok := data.([]byte); ok {
		_, err = io.ReadFull(p.reader, buf)
	} else {
		err = binary.Read(p.reader, p.ByteOrder, data)
	}
	if err != nil {
		p.OnError(err)
	}
	return
}

// ReadInt64 reads a int64 value.
func (p *BinaryReader) ReadInt64() (result int64, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt64 reads a uint64 value.
func (p *BinaryReader) ReadUInt64() (result uint64, err error) {
	err = p.Read(&result)
	return
}

// ReadFloat64 reads a float64 value.
func (p *BinaryReader) ReadFloat64() (result float64, err error) {
	err = p.Read(&result)
	return
}

// ReadInt32 reads a int32 value.
func (p *BinaryReader) ReadInt32() (result int32, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt32 reads a uint32 value.
func (p *BinaryReader) ReadUInt32() (result uint32, err error) {
	err = p.Read(&result)
	return
}

// ReadFloat32 reads a float32 value.
func (p *BinaryReader) ReadFloat32() (result float32, err error) {
	err = p.Read(&result)
	return
}

// ReadInt16 reads a int16 value.
func (p *BinaryReader) ReadInt16() (result int16, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt16 reads a uint16 value.
func (p *BinaryReader) ReadUInt16() (result uint16, err error) {
	err = p.Read(&result)
	return
}

// ReadInt8 reads a int8 value.
func (p *BinaryReader) ReadInt8() (result int8, err error) {
	err = p.Read(&result)
	return
}

// ReadUInt8 reads a uint8 value.
func (p *BinaryReader) ReadUInt8() (result uint8, err error) {
	err = p.Read(&result)
	return
}

// ReadByte reads a byte.
func (p *BinaryReader) ReadByte() (result byte, err error) {
	err = p.Read(&result)
	return
}

// ReadBytes reads bytes in length.
func (p *BinaryReader) ReadBytes(length uint) (result []byte, err error) {
	buf := make([]byte, length)
	err = p.Read(buf)
	if err != nil {
		return
	}
	result = buf
	return
}

// ReadBytesUntil reads bytes until reaches delim byte.
func (p *BinaryReader) ReadBytesUntil(delim byte) (result []byte, err error) {
	var b byte
	buf := &bytes.Buffer{}
	for b, err = p.ReadByte(); err == nil && b != delim; b, err = p.ReadByte() {
		buf.WriteByte(b)
	}
	if err != nil {
		return
	}
	result = buf.Bytes()
	return
}

// ReadString reads string until reaches null character.
func (p *BinaryReader) ReadString() (result string, err error) {
	data, err := p.ReadBytesUntil(0)
	if err != nil {
		return
	}
	result = string(data)
	return
}

package tiff

import (
	"encoding/binary"
	"io"
)

// ReadHeader reads a TIFF header from r.
func ReadHeader(r *io.SectionReader) (result *Header, err error) {
	var header struct {
		ByteOrder   [2]byte
		Number42    [2]byte
		OffsetOfIFD [4]byte
	}
	if err = binary.Read(r, binary.BigEndian, &header); err != nil {
		return
	}
	var order binary.ByteOrder
	switch header.ByteOrder {
	case [2]byte{0x49, 0x49}:
		order = binary.LittleEndian
	case [2]byte{0x4d, 0x4d}:
		order = binary.BigEndian
	default:
		err = ErrInvalidTIFFHeader
		return
	}
	if order.Uint16(header.Number42[:]) != 0x2a {
		err = ErrInvalidTIFFHeader
		return
	}
	result = &Header{
		Reader:      r,
		Order:       order,
		OffsetOfIFD: order.Uint32(header.OffsetOfIFD[:]),
	}
	return
}

// Header represents header of Tag Image File Format.
type Header struct {
	Reader      *io.SectionReader
	Order       binary.ByteOrder
	OffsetOfIFD uint32
}

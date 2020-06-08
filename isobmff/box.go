package isobmff

import (
	"encoding/binary"
	"io"
)

var order = binary.BigEndian

// Box represents ISO-BMFF box.
type Box interface {
	BoxReader
	Size() int64
	Type() string
	UserType() string
}

type boxBase struct {
	reader   *io.SectionReader
	boxSize  int64
	boxType  string
	userType string
}

func (p *boxBase) Size() int64 {
	return p.boxSize
}

func (p *boxBase) Type() string {
	return p.boxType
}

func (p *boxBase) UserType() string {
	return p.userType
}

func (p *boxBase) ReadBox() (result Box, err error) {
	off, err := p.reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}
	limit := p.boxSize - off
	result, err = NewBoxReader(io.NewSectionReader(p.reader, off, limit)).ReadBox()
	if err != nil {
		return
	}
	_, err = p.reader.Seek(result.Size(), io.SeekCurrent)
	return
}

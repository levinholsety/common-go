package isobmff

import (
	"io"
	"os"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/utils"
)

var boxBuilderMap = map[string]func(boxBase *boxBase) (Box, error){
	"ftyp": buildFileTypeBox,
	"meta": buildMetaBox,
	"iinf": buildItemInfoBox,
	"infe": buildItemInfoEntry,
	"iloc": buildItemLocationBox,
	"moov": buildMovieBox,
	"mvhd": buildMovieHeaderBox,
}

// BoxReader provides methods to read box.
type BoxReader interface {
	ReadBox() (Box, error)
}

// NewBoxReader creates a BoxReader from *io.SectionReader.
func NewBoxReader(r *io.SectionReader) BoxReader {
	return &boxReader{
		reader: r,
	}
}

type boxReader struct {
	reader *io.SectionReader
}

func (p *boxReader) ReadBox() (result Box, err error) {
	off, err := p.reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}
	box := &boxBase{}
	r := utils.NewBinaryReader(p.reader, order)
	boxSize, err := r.ReadUInt32()
	if err != nil {
		return
	}
	box.boxType, err = r.ReadStringFixed(4)
	if err != nil {
		return
	}
	if boxSize == 0 {
		box.boxSize = p.reader.Size() - off
	} else if boxSize == 1 {
		box.boxSize, err = r.ReadInt64()
		if err != nil {
			return
		}
	} else {
		box.boxSize = int64(boxSize)
	}
	if box.boxType == "uuid" {
		box.userType, err = r.ReadStringFixed(16)
		if err != nil {
			return
		}
	}
	skip, err := p.reader.Seek(0, io.SeekCurrent)
	if err != nil {
		return
	}
	skip -= off
	box.reader = io.NewSectionReader(p.reader, off, box.boxSize)
	_, err = box.reader.Seek(skip, io.SeekStart)
	if err != nil {
		return
	}
	if buildBox, ok := boxBuilderMap[box.boxType]; ok {
		result, err = buildBox(box)
	} else {
		result = box
	}
	_, err = p.reader.Seek(off+result.Size(), io.SeekStart)
	return
}

// NewBoxFile creates a *BoxFile and returns it.
func NewBoxFile(file *os.File) (result BoxReader, err error) {
	r, err := comm.FileToSectionReader(file)
	if err != nil {
		return
	}
	result = NewBoxReader(r)
	return
}

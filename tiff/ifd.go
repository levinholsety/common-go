package tiff

import (
	"io"

	"github.com/levinholsety/common-go/utils"
)

// DE tags.
const (
	TagCompression                 uint16 = 0x103
	TagMake                        uint16 = 0x10f
	TagModel                       uint16 = 0x110
	TagOrientation                 uint16 = 0x112
	TagSoftware                    uint16 = 0x131
	TagDateTime                    uint16 = 0x132
	TagJPEGInterchangeFormat       uint16 = 0x201
	TagJPEGInterchangeFormatLength uint16 = 0x202
)

// ReadIFD reads an image file directory.
func ReadIFD(tiffHeader *Header, offset uint32) (result *IFD, err error) {
	result = &IFD{
		tiff: tiffHeader,
	}
	r := utils.NewBinaryReader(io.NewSectionReader(tiffHeader.Reader, int64(offset), tiffHeader.Reader.Size()-int64(offset)), tiffHeader.Order)
	if err = r.Read(&result.entryCount); err != nil {
		return
	}
	result.entries = map[uint16]*DE{}
	for i := uint16(0); i < result.entryCount; i++ {
		var entry struct {
			Tag         uint16
			Type        uint16
			Count       uint32
			ValueOffset [4]byte
		}
		if err = r.Read(&entry); err != nil {
			return
		}
		result.entries[entry.Tag] = &DE{
			ifd:         result,
			Tag:         entry.Tag,
			ValueType:   ValueType(entry.Type),
			ValueCount:  entry.Count,
			valueOffset: entry.ValueOffset,
		}
	}
	err = r.Read(&result.OffsetOfNextIFD)
	return
}

// IFD represents Image File Directory.
type IFD struct {
	tiff            *Header
	entryCount      uint16
	entries         map[uint16]*DE
	OffsetOfNextIFD uint32
}

// Entry returns the directory entry which tag is specified.
func (p *IFD) Entry(tag uint16) (result *DE, ok bool) {
	result, ok = p.entries[tag]
	return
}

package isobmff

import (
	"io"

	"github.com/levinholsety/common-go/utils"
)

func buildFileTypeBox(boxBase *boxBase) (result Box, err error) {
	box := &FileTypeBox{boxBase: boxBase}
	r := utils.NewBinaryReader(box.reader, order)
	if box.MajorBrand, err = r.ReadStringFixed(4); err != nil {
		return
	}
	if box.MinorVersion, err = r.ReadUInt32(); err != nil {
		return
	}
	for {
		var compatibleBrand string
		compatibleBrand, err = r.ReadStringFixed(4)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		box.CompatibleBrands = append(box.CompatibleBrands, compatibleBrand)
	}
	result = box
	return
}

// FileTypeBox represents ISO-BMFF FileTypeBox.
type FileTypeBox struct {
	*boxBase
	MajorBrand       string
	MinorVersion     uint32
	CompatibleBrands []string
}

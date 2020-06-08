package isobmff

import (
	"github.com/levinholsety/common-go/utils"
)

func newFullBox(boxBase *boxBase) (result *FullBox, err error) {
	result = &FullBox{boxBase: boxBase}
	r := utils.NewBinaryReader(result.reader, order)
	value, err := r.ReadUInt32()
	if err != nil {
		return
	}
	result.Version = value >> 24
	result.Flags = value & 0xffffff
	return
}

// FullBox represents ISO-BMFF FullBox.
type FullBox struct {
	*boxBase
	Version uint32
	Flags   uint32
}

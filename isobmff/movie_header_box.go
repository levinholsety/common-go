package isobmff

import (
	"github.com/levinholsety/common-go/utils"
)

func buildMovieHeaderBox(boxBase *boxBase) (result Box, err error) {
	fullBox, err := newFullBox(boxBase)
	if err != nil {
		return
	}
	box := &MovieHeaderBox{FullBox: fullBox}
	r := utils.NewBinaryReader(box.reader, order)
	if box.Version == 1 {
		var data struct {
			CreationTime     uint64
			ModificationTime uint64
			Timescale        uint32
			Duration         uint64
		}
		if err = r.Read(&data); err != nil {
			return
		}
		box.CreationTime = uint(data.CreationTime)
		box.ModificationTime = uint(data.ModificationTime)
		box.Timescale = uint(data.Timescale)
		box.Duration = uint(data.Duration)
	} else {
		var data struct {
			CreationTime     uint32
			ModificationTime uint32
			Timescale        uint32
			Duration         uint32
		}
		if err = r.Read(&data); err != nil {
			return
		}
		box.CreationTime = uint(data.CreationTime)
		box.ModificationTime = uint(data.ModificationTime)
		box.Timescale = uint(data.Timescale)
		box.Duration = uint(data.Duration)
	}
	result = box
	return
}

// MovieHeaderBox represents ISO-BMFF MovieHeaderBox.
type MovieHeaderBox struct {
	*FullBox
	CreationTime     uint
	ModificationTime uint
	Timescale        uint
	Duration         uint
}

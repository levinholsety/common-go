package isobmff

import (
	"io"

	"github.com/levinholsety/common-go/utils"
)

func buildItemLocationBox(boxBase *boxBase) (result Box, err error) {
	fullBox, err := newFullBox(boxBase)
	if err != nil {
		return
	}
	box := &ItemLocationBox{FullBox: fullBox}
	r := utils.NewBinaryReader(box.reader, order)
	var sizes [2]byte
	if err = r.Read(&sizes); err != nil {
		return
	}
	box.OffsetSize = sizes[0] >> 4
	box.LengthSize = sizes[0] & 0b1111
	box.BaseOffsetSize = sizes[1] >> 4
	if box.Version == 1 || box.Version == 2 {
		box.IndexSize = sizes[1] & 0b1111
	}
	if box.Version < 2 {
		var itemCount uint16
		itemCount, err = r.ReadUInt16()
		if err != nil {
			return
		}
		box.ItemCount = uint32(itemCount)
	} else if box.Version == 2 {
		box.ItemCount, err = r.ReadUInt32()
		if err != nil {
			return
		}
	}
	for i := uint32(0); i < box.ItemCount; i++ {
		var item ItemLocationBoxItem
		if box.Version < 2 {
			var itemID uint16
			itemID, err = r.ReadUInt16()
			if err != nil {
				return
			}
			item.ItemID = uint32(itemID)
		} else if box.Version == 2 {
			item.ItemID, err = r.ReadUInt32()
			if err != nil {
				return
			}
		}
		if box.Version == 1 || box.Version == 2 {
			if _, err = box.reader.Seek(2, io.SeekCurrent); err != nil {
				return
			}
		}
		length := 2 + int64(box.BaseOffsetSize)
		if _, err = box.reader.Seek(length, io.SeekCurrent); err != nil {
			return
		}
		var extentCount uint16
		extentCount, err = r.ReadUInt16()
		if err != nil {
			return
		}
		for j := uint16(0); j < extentCount; j++ {
			var extent ItemLocationBoxItemExtent
			if (box.Version == 1 || box.Version == 2) && box.IndexSize > 0 {
				if _, err = box.reader.Seek(int64(box.IndexSize), io.SeekCurrent); err != nil {
					return
				}
			}
			if extent.ExtentOffset, err = readBySize(r, box.OffsetSize); err != nil {
				return
			}
			if extent.ExtentLength, err = readBySize(r, box.LengthSize); err != nil {
				return
			}
			item.Extents = append(item.Extents, extent)
		}
		box.Items = append(box.Items, item)
	}
	result = box
	return
}

func readBySize(r *utils.BinaryReader, size uint8) (result uint, err error) {
	switch size {
	case 0:
		result = 0
	case 4:
		var value uint32
		if err = r.Read(&value); err != nil {
			return
		}
		result = uint(value)
	case 8:
		var value uint64
		if err = r.Read(&value); err != nil {
			return
		}
		result = uint(value)
	default:
		result = 0
	}
	return
}

// ItemLocationBox represents ISO-BMFF ItemLocationBox.
type ItemLocationBox struct {
	*FullBox
	OffsetSize     uint8
	LengthSize     uint8
	BaseOffsetSize uint8
	IndexSize      uint8
	ItemCount      uint32
	Items          []ItemLocationBoxItem
}

// ItemLocationBoxItem represents ISO-BMFF ItemLocationBoxItem
type ItemLocationBoxItem struct {
	ItemID  uint32
	Extents []ItemLocationBoxItemExtent
}

// ItemLocationBoxItemExtent represents ISO-BMFF ItemLocationBoxItemExtent
type ItemLocationBoxItemExtent struct {
	ExtentOffset uint
	ExtentLength uint
}

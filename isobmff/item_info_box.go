package isobmff

import (
	"github.com/levinholsety/common-go/utils"
)

func buildItemInfoBox(boxBase *boxBase) (result Box, err error) {
	fullBox, err := newFullBox(boxBase)
	if err != nil {
		return
	}
	box := &ItemInfoBox{FullBox: fullBox}
	r := utils.NewBinaryReader(box.reader, order)
	if box.Version == 0 {
		var entryCount uint16
		entryCount, err = r.ReadUInt16()
		if err != nil {
			return
		}
		box.EntryCount = uint32(entryCount)
	} else {
		box.EntryCount, err = r.ReadUInt32()
	}
	result = box
	return
}

// ItemInfoBox represents ISO-BMFF ItemInfoBox.
type ItemInfoBox struct {
	*FullBox
	EntryCount uint32
	ItemInfos  []*ItemInfoEntry
}

func buildItemInfoEntry(boxBase *boxBase) (result Box, err error) {
	fullBox, err := newFullBox(boxBase)
	if err != nil {
		return
	}
	box := &ItemInfoEntry{FullBox: fullBox}
	r := utils.NewBinaryReader(box.reader, order)
	if box.Version >= 2 {
		if box.Version == 2 {
			var itemID uint16
			itemID, err = r.ReadUInt16()
			if err != nil {
				return
			}
			box.ItemID = uint32(itemID)
		} else if box.Version == 3 {
			box.ItemID, err = r.ReadUInt32()
			if err != nil {
				return
			}
		}
		if _, err = r.ReadUInt16(); err != nil {
			return
		}
		box.ItemType, err = r.ReadStringFixed(4)
	}
	result = box
	return
}

// ItemInfoEntry represents ISO-BMFF ItemInfoEntry.
type ItemInfoEntry struct {
	*FullBox
	ItemID   uint32
	ItemType string
}

package isobmff

func buildMetaBox(boxBase *boxBase) (result Box, err error) {
	fullBox, err := newFullBox(boxBase)
	if err != nil {
		return
	}
	result = &MetaBox{FullBox: fullBox}
	return
}

// MetaBox represents ISO-BMFF MetaBox.
type MetaBox struct {
	*FullBox
}

package isobmff

func buildMovieBox(boxBase *boxBase) (result Box, err error) {
	result = &MovieBox{boxBase: boxBase}
	return
}

// MovieBox represents ISO-BMFF MovieBox.
type MovieBox struct {
	*boxBase
}

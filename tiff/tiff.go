// Package tiff provides methods for operating Tag Image File Format.
package tiff

import (
	"errors"
)

// errors
var (
	ErrInvalidTIFFHeader = errors.New("invalid tiff header")
)

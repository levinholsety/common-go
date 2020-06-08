package exif

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/levinholsety/common-go/comm"
	"github.com/levinholsety/common-go/isobmff"
	"github.com/levinholsety/common-go/num"
	"github.com/levinholsety/common-go/tiff"
	"github.com/levinholsety/common-go/utils"
)

// Entry tags.
const (
	TagExposureTime          uint16 = 0x829a
	TagFNumber               uint16 = 0x829d
	TagExifIFDPointer        uint16 = 0x8769
	TagGPSInfoIFDPointer     uint16 = 0x8825
	TagISOSpeedRatings       uint16 = 0x8827
	TagDateTimeOriginal      uint16 = 0x9003
	TagDateTimeDigitized     uint16 = 0x9004
	TagFocalLength           uint16 = 0x920a
	TagMakerNote             uint16 = 0x927c
	TagSubsecTime            uint16 = 0x9290
	TagSubSecTimeOriginal    uint16 = 0x9291
	TagSubSecTimeDigitized   uint16 = 0x9292
	TagFocalLengthIn35mmFilm uint16 = 0xa405
	TagGPSLatitude           uint16 = 0x2
	TagGPSLongitude          uint16 = 0x4
)

var (
	soi               uint16 = 0xffd8
	app1              uint16 = 0xffe1
	sos               uint16 = 0xffda
	exifMarker               = [6]byte{0x45, 0x78, 0x69, 0x66, 0x00, 0x00}
	nikonType1Marker         = []byte{0x4e, 0x69, 0x6b, 0x6f, 0x6e, 0x00, 0x01, 0x00}
	nikonType3AMarker        = []byte{0x4e, 0x69, 0x6b, 0x6f, 0x6e, 0x00, 0x02, 0x10, 0x00, 0x00}
	nikonType3BMarker        = []byte{0x4e, 0x69, 0x6b, 0x6f, 0x6e, 0x00, 0x02, 0x00, 0x00, 0x00}
)

// Parse parses EXIF from reader.
func Parse(file *os.File) (result *Info, err error) {
	r, err := comm.FileToSectionReader(file)
	if err != nil {
		return
	}
	var off int64
	for _, f := range []func(r *io.SectionReader) (int64, int64, error){
		findExifInJPEG,
		findExifInHEIC,
	} {
		off, _, err = f(r)
		if err == io.EOF {
			continue
		}
		if err != nil {
			return
		}
		break
	}
	tiffHeader, err := tiff.ReadHeader(io.NewSectionReader(r, off, r.Size()-off))
	if err != nil {
		return
	}
	result, err = parse(tiffHeader)
	return
}

func findExifInJPEG(r *io.SectionReader) (offset, length int64, err error) {
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	br := utils.NewBinaryReader(r, binary.BigEndian)
	var marker uint16
	if err = br.Read(&marker); err != nil {
		return
	}
	if marker == soi {
		for true {
			var appMarker struct {
				ID   uint16
				Size uint16
			}
			if err = br.Read(&appMarker); err != nil {
				return
			}
			if appMarker.ID == app1 {
				var buf6 [6]byte
				if err = br.Read(&buf6); err != nil {
					return
				}
				if exifMarker == buf6 {
					if offset, err = r.Seek(0, io.SeekCurrent); err != nil {
						return
					}
					length = int64(appMarker.Size) - 8
					return
				}
			} else if appMarker.ID == sos {
				return
			} else {
				if _, err = r.Seek(int64(appMarker.Size)-2, io.SeekCurrent); err != nil {
					return
				}
			}
		}
	}
	err = io.EOF
	return
}

func findExifInHEIC(r *io.SectionReader) (offset, length int64, err error) {
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return
	}
	boxFile := isobmff.NewBoxReader(r)
	var box isobmff.Box
	if box, err = boxFile.ReadBox(); err != nil {
		return
	}
	ftyp, ok := box.(*isobmff.FileTypeBox)
	if ok && ftyp.MajorBrand == "heic" {
		if box, err = boxFile.ReadBox(); err != nil {
			return
		}
		var meta *isobmff.MetaBox
		if meta, ok = box.(*isobmff.MetaBox); ok {
			var itemID uint32
			for {
				if box, err = meta.ReadBox(); err != nil {
					if err == io.EOF {
						err = nil
						break
					}
					return
				}
				var (
					iinf *isobmff.ItemInfoBox
					iloc *isobmff.ItemLocationBox
				)
				if iinf, ok = box.(*isobmff.ItemInfoBox); ok {
					for i := uint32(0); i < iinf.EntryCount; i++ {
						if box, err = iinf.ReadBox(); err != nil {
							return
						}
						infe := box.(*isobmff.ItemInfoEntry)
						if infe.ItemType == "Exif" {
							itemID = infe.ItemID
						}
					}
				} else if iloc, ok = box.(*isobmff.ItemLocationBox); ok {
					for _, item := range iloc.Items {
						if item.ItemID == itemID {
							extent := item.Extents[0]
							offset = int64(extent.ExtentOffset) + 10
							length = int64(extent.ExtentLength) - 10
							return
						}
					}
				}
			}
		}
	}
	err = io.EOF
	return
}

// Info represents Exchangeable image file format information.
type Info struct {
	tiffHeader   *tiff.Header
	ifd0         *tiff.IFD
	ifd1         *tiff.IFD
	exifIFD      *tiff.IFD
	gpsInfoIFD   *tiff.IFD
	makerNoteIFD *tiff.IFD
}

const (
	makeCanon = "Canon"
	makeNikon = "NIKON CORPORATION"
)

func parse(tiffHeader *tiff.Header) (result *Info, err error) {
	result = &Info{tiffHeader: tiffHeader}
	result.ifd0, err = tiff.ReadIFD(tiffHeader, tiffHeader.OffsetOfIFD)
	if err != nil {
		return
	}
	if result.ifd0.OffsetOfNextIFD > 0 {
		if result.ifd1, err = tiff.ReadIFD(tiffHeader, result.ifd0.OffsetOfNextIFD); err != nil {
			return
		}
	}
	if exifIFDPointer, ok := result.ifd0.Entry(TagExifIFDPointer); ok && exifIFDPointer.ValueOffset() > 0 {
		if result.exifIFD, err = tiff.ReadIFD(tiffHeader, exifIFDPointer.ValueOffset()); err != nil {
			return
		}
		if makerNoteOffset, ok := result.exifIFD.Entry(TagMakerNote); ok && makerNoteOffset.ValueOffset() > 0 {
			if makeEntry, ok := result.ifd0.Entry(tiff.TagMake); ok {
				var makeValue string
				makeValue, err = makeEntry.StringValue()
				if err != nil {
					return
				}
				if makeValue == makeCanon {
					result.makerNoteIFD, err = tiff.ReadIFD(tiffHeader, makerNoteOffset.ValueOffset())
					if err != nil {
						return
					}
				} else if makeValue == makeNikon {
					var buf = make([]byte, 10)
					if _, err = tiffHeader.Reader.ReadAt(buf, int64(makerNoteOffset.ValueOffset())); err != nil {
						return
					}
					if err == nil {
						if bytes.Equal(buf[:8], nikonType1Marker) {
							// Nikon Type 1 Makernote
							result.makerNoteIFD, err = tiff.ReadIFD(tiffHeader, makerNoteOffset.ValueOffset()+8)
							if err != nil {
								return
							}
						} else if bytes.Equal(buf, nikonType3AMarker) || bytes.Equal(buf, nikonType3BMarker) {
							// Nikon Type 3 Makernote
							var nikonTIFFHeader *tiff.Header
							off := int64(makerNoteOffset.ValueOffset()) + 10
							if nikonTIFFHeader, err = tiff.ReadHeader(io.NewSectionReader(tiffHeader.Reader, off, tiffHeader.Reader.Size()-off)); err != nil {
								return
							}
							if err == nil {
								result.makerNoteIFD, err = tiff.ReadIFD(nikonTIFFHeader, nikonTIFFHeader.OffsetOfIFD)
								if err != nil {
									return
								}
							}
						} else {
							// Nikon Type 2 Makernote
							result.makerNoteIFD, err = tiff.ReadIFD(tiffHeader, makerNoteOffset.ValueOffset())
							if err != nil {
								return
							}
						}
					}
				}
			}
		}
	}
	if gpsInfoIFDPointer, ok := result.ifd0.Entry(TagGPSInfoIFDPointer); ok && gpsInfoIFDPointer.ValueOffset() > 0 {
		if result.gpsInfoIFD, err = tiff.ReadIFD(tiffHeader, gpsInfoIFDPointer.ValueOffset()); err != nil {
			return
		}
	}
	err = nil
	return
}

// Make returns Make.
func (p *Info) Make() (result string, err error) {
	if entry, ok := p.ifd0.Entry(tiff.TagMake); ok {
		return entry.StringValue()
	}
	return
}

// Model returns Model.
func (p *Info) Model() (result string, err error) {
	if entry, ok := p.ifd0.Entry(tiff.TagModel); ok {
		return entry.StringValue()
	}
	return
}

// DateTime returns DateTime.
func (p *Info) DateTime() (result string, err error) {
	if entry, ok := p.ifd0.Entry(tiff.TagDateTime); ok {
		return entry.StringValue()
	}
	return
}

// ExposureTime returns ExposureTime.
func (p *Info) ExposureTime() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagExposureTime); ok {
		var value interface{}
		value, err = entry.Value()
		if err != nil {
			return
		}
		if v, ok := value.(num.Fraction); ok {
			if v.Numerator() < v.Denominator() {
				result = fmt.Sprintf("1/%d", int(v.Reciprocal().Float64()))
			} else {
				result = fmt.Sprintf("%g", v.Float64())
			}
		}
		return
	}
	return
}

// FNumber returns FNumber.
func (p *Info) FNumber() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagFNumber); ok {
		var v num.Fraction
		v, err = entry.RationalValue()
		if err != nil {
			return
		}
		result = fmt.Sprintf("%g", v.Float64())
	}
	return
}

// DateTimeOriginal returns date and time of original data generation.
func (p *Info) DateTimeOriginal() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagDateTimeOriginal); ok {
		return entry.StringValue()
	}
	return
}

// DateTimeDigitized returns date and time of digital data generation.
func (p *Info) DateTimeDigitized() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagDateTimeDigitized); ok {
		return entry.StringValue()
	}
	return
}

// SubsecTime returns DateTime subseconds.
func (p *Info) SubsecTime() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagSubsecTime); ok {
		return entry.StringValue()
	}
	return
}

// SubSecTimeOriginal returns DateTimeOriginal subseconds.
func (p *Info) SubSecTimeOriginal() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagSubSecTimeOriginal); ok {
		return entry.StringValue()
	}
	return
}

// SubSecTimeDigitized returns DateTimeDigitized subseconds.
func (p *Info) SubSecTimeDigitized() (result string, err error) {
	if p.exifIFD == nil {
		return
	}
	if entry, ok := p.exifIFD.Entry(TagSubSecTimeDigitized); ok {
		return entry.StringValue()
	}
	return
}

func dms(array []num.Fraction) float64 {
	return array[0].Float64() + array[1].Float64()/60 + array[2].Float64()/3600
}

// GPSLatitude returns Latitude.
func (p *Info) GPSLatitude() (result float64, err error) {
	if p.gpsInfoIFD == nil {
		return
	}
	if entry, ok := p.gpsInfoIFD.Entry(TagGPSLatitude); ok {
		var value interface{}
		value, err = entry.Value()
		if err != nil {
			return
		}
		if array, ok := value.([]num.Fraction); ok && len(array) == 3 {
			result = dms(array)
		}
	}
	return
}

// GPSLongitude returns Longitude.
func (p *Info) GPSLongitude() (result float64, err error) {
	if p.gpsInfoIFD == nil {
		return
	}
	if entry, ok := p.gpsInfoIFD.Entry(TagGPSLongitude); ok {
		var value interface{}
		value, err = entry.Value()
		if err != nil {
			return
		}
		if array, ok := value.([]num.Fraction); ok && len(array) == 3 {
			result = dms(array)
		}
	}
	return
}

// ShutterCount returns shutter count.
func (p *Info) ShutterCount() (result uint32, err error) {
	if p.makerNoteIFD == nil {
		return
	}
	makeValue, err := p.Make()
	if err != nil {
		return
	}
	if makeValue == makeCanon {
		if entry, ok := p.makerNoteIFD.Entry(0x93); ok {
			var array []uint16
			array, err = entry.ShortValues()
			if err != nil {
				return
			}
			if len(array) > 2 && array[2] > 0 {
				result = uint32(array[2])
			}
		}
	} else if makeValue == makeNikon {
		if entry, ok := p.makerNoteIFD.Entry(0xa7); ok {
			return entry.LongValue()
		}
	}
	return
}

// Lens returns lens.
func (p *Info) Lens() (result string, err error) {
	if p.makerNoteIFD == nil {
		return
	}
	makeValue, err := p.Make()
	if err != nil {
		return
	}
	if makeValue == makeCanon {
		if entry, ok := p.makerNoteIFD.Entry(0x95); ok {
			result, err = entry.StringValue()
		}
	} else if makeValue == makeNikon {
		var value interface{}
		if entry, ok := p.makerNoteIFD.Entry(0x84); ok {
			value, err = entry.Value()
			if err != nil {
				return
			}
			var array []num.Fraction
			if array, ok = value.([]num.Fraction); ok && len(array) == 4 {
				v1 := int(array[0].Float64())
				v2 := int(array[1].Float64())
				v3 := array[2].Float64()
				v4 := array[3].Float64()
				var s1, s2 string
				if v1 == v2 {
					s1 = fmt.Sprintf("%d", v1)
				} else {
					s1 = fmt.Sprintf("%d-%d", v1, v2)
				}
				if v3 == v4 {
					s2 = fmt.Sprintf("%g", v3)
				} else {
					s2 = fmt.Sprintf("%g-%g", v3, v4)
				}
				result = fmt.Sprintf("%smm f/%s", s1, s2)
				return
			}
		}
	}
	return
}

// CanonImageNumber returns image number of canon device.
func (p *Info) CanonImageNumber() (result string, err error) {
	if p.makerNoteIFD == nil {
		return
	}
	if entry, ok := p.makerNoteIFD.Entry(0x8); ok {
		var value interface{}
		value, err = entry.Value()
		if err != nil {
			return
		}
		var v uint32
		if v, ok = value.(uint32); ok {
			result = fmt.Sprintf("%d-%d", v/10000, v%10000)
			return
		}
	}
	return
}

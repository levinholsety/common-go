package tiff

import (
	"bytes"
	"io"
	"math"

	"github.com/levinholsety/common-go/num"
)

// ValueType represents value type of directory entry.
type ValueType uint16

// Directory Entry types.
const (
	TypeBYTE ValueType = iota + 1
	TypeASCII
	TypeSHORT
	TypeLONG
	TypeRATIONAL
	TypeSBYTE
	TypeUNDEFINED
	TypeSSHORT
	TypeSLONG
	TypeSRATIONAL
	TypeFLOAT
	TypeDOUBLE
)

// Size returns the size of the value type.
func (v ValueType) Size() (size uint32) {
	switch v {
	case TypeBYTE, TypeASCII, TypeSBYTE, TypeUNDEFINED:
		size = 1
	case TypeSHORT, TypeSSHORT:
		size = 2
	case TypeLONG, TypeSLONG, TypeFLOAT:
		size = 4
	case TypeRATIONAL, TypeSRATIONAL, TypeDOUBLE:
		size = 8
	}
	return
}

// DE represents Directory Entry
type DE struct {
	ifd         *IFD
	Tag         uint16
	ValueType   ValueType
	ValueCount  uint32
	valueOffset [4]byte
}

// ValueOffset returns value offset.
func (p *DE) ValueOffset() uint32 {
	return p.ifd.tiff.Order.Uint32(p.valueOffset[:])
}

// BytesValue returns value in bytes.
func (p *DE) BytesValue() (result []byte, err error) {
	valueSize := p.ValueType.Size()
	if valueSize == 0 {
		return
	}
	size := int(valueSize * p.ValueCount)
	result = make([]byte, size)
	valueOffset := p.valueOffset[:]
	if size <= 4 {
		copy(result, valueOffset)
		return
	}
	offset := int64(p.ifd.tiff.Order.Uint32(valueOffset))
	n, err := p.ifd.tiff.Reader.ReadAt(result, offset)
	if n < size {
		err = io.ErrUnexpectedEOF
		return
	}
	if err == io.EOF {
		err = nil
	}
	return
}

// StringValues returns value in string array.
func (p *DE) StringValues() (result []string, err error) {
	data, err := p.BytesValue()
	if err != nil {
		return
	}
	for len(data) > 0 {
		index := bytes.IndexByte(data, 0)
		if index < 0 {
			result = append(result, string(data))
			return
		}
		if index > 0 {
			result = append(result, string(data[:index]))
			data = data[index+1:]
			continue
		}
		data = data[1:]
	}
	return
}

// StringValue returns value in string.
func (p *DE) StringValue() (result string, err error) {
	values, err := p.StringValues()
	if err != nil {
		return
	}
	if len(values) == 0 {
		return
	}
	result = values[0]
	return
}

func (p *DE) numberValues(valueType ValueType, f func(value []byte) bool) (err error) {
	data, err := p.BytesValue()
	if err != nil {
		return
	}
	valueSize := p.ValueType.Size()
	for len(data) > 0 {
		if !f(data[:valueSize]) {
			break
		}
		data = data[valueSize:]
	}
	return
}

// ShortValues returns value in uint16 array.
func (p *DE) ShortValues() (result []uint16, err error) {
	err = p.numberValues(TypeSHORT, func(value []byte) bool {
		result = append(result, p.ifd.tiff.Order.Uint16(value))
		return true
	})
	return
}

// ShortValue returns value in uint16.
func (p *DE) ShortValue() (result uint16, err error) {
	err = p.numberValues(TypeSHORT, func(value []byte) bool {
		result = p.ifd.tiff.Order.Uint16(value)
		return false
	})
	return
}

// SShortValues returns value in int16 array.
func (p *DE) SShortValues() (result []int16, err error) {
	err = p.numberValues(TypeSSHORT, func(value []byte) bool {
		result = append(result, int16(p.ifd.tiff.Order.Uint16(value)))
		return true
	})
	return
}

// SShortValue returns value in int16.
func (p *DE) SShortValue() (result int16, err error) {
	err = p.numberValues(TypeSSHORT, func(value []byte) bool {
		result = int16(p.ifd.tiff.Order.Uint16(value))
		return false
	})
	return
}

// LongValues returns value in uint32 array.
func (p *DE) LongValues() (result []uint32, err error) {
	err = p.numberValues(TypeLONG, func(value []byte) bool {
		result = append(result, p.ifd.tiff.Order.Uint32(value))
		return true
	})
	return
}

// LongValue returns value in uint32.
func (p *DE) LongValue() (result uint32, err error) {
	err = p.numberValues(TypeLONG, func(value []byte) bool {
		result = p.ifd.tiff.Order.Uint32(value)
		return false
	})
	return
}

// SLongValues returns value in int32 array.
func (p *DE) SLongValues() (result []int32, err error) {
	err = p.numberValues(TypeSLONG, func(value []byte) bool {
		result = append(result, int32(p.ifd.tiff.Order.Uint32(value)))
		return true
	})
	return
}

// SLongValue returns value in int32.
func (p *DE) SLongValue() (result int32, err error) {
	err = p.numberValues(TypeSLONG, func(value []byte) bool {
		result = int32(p.ifd.tiff.Order.Uint32(value))
		return false
	})
	return
}

// FloatValues returns value in float32 array.
func (p *DE) FloatValues() (result []float32, err error) {
	err = p.numberValues(TypeFLOAT, func(value []byte) bool {
		result = append(result, math.Float32frombits(p.ifd.tiff.Order.Uint32(value)))
		return true
	})
	return
}

// FloatValue returns value in float32.
func (p *DE) FloatValue() (result float32, err error) {
	err = p.numberValues(TypeFLOAT, func(value []byte) bool {
		result = math.Float32frombits(p.ifd.tiff.Order.Uint32(value))
		return false
	})
	return
}

// RationalValues returns value in fraction array.
func (p *DE) RationalValues() (result []num.Fraction, err error) {
	err = p.numberValues(TypeRATIONAL, func(value []byte) bool {
		result = append(result, num.NewFraction(
			int(p.ifd.tiff.Order.Uint32(value[:4])),
			int(p.ifd.tiff.Order.Uint32(value[4:])),
		))
		return true
	})
	return
}

// RationalValue returns value in fraction.
func (p *DE) RationalValue() (result num.Fraction, err error) {
	err = p.numberValues(TypeRATIONAL, func(value []byte) bool {
		result = num.NewFraction(
			int(p.ifd.tiff.Order.Uint32(value[:4])),
			int(p.ifd.tiff.Order.Uint32(value[4:])),
		)
		return false
	})
	return
}

// SRationalValues returns value in signed fraction array.
func (p *DE) SRationalValues() (result []num.Fraction, err error) {
	err = p.numberValues(TypeSRATIONAL, func(value []byte) bool {
		result = append(result, num.NewFraction(
			int(int32(p.ifd.tiff.Order.Uint32(value[:4]))),
			int(int32(p.ifd.tiff.Order.Uint32(value[4:]))),
		))
		return true
	})
	return
}

// SRationalValue returns value in signed fraction.
func (p *DE) SRationalValue() (result num.Fraction, err error) {
	err = p.numberValues(TypeSRATIONAL, func(value []byte) bool {
		result = num.NewFraction(
			int(int32(p.ifd.tiff.Order.Uint32(value[:4]))),
			int(int32(p.ifd.tiff.Order.Uint32(value[4:]))),
		)
		return false
	})
	return
}

// DoubleValues returns value in float64 array.
func (p *DE) DoubleValues() (result []float64, err error) {
	err = p.numberValues(TypeDOUBLE, func(value []byte) bool {
		result = append(result, math.Float64frombits(p.ifd.tiff.Order.Uint64(value)))
		return true
	})
	return
}

// DoubleValue returns value in float64.
func (p *DE) DoubleValue() (result float64, err error) {
	err = p.numberValues(TypeDOUBLE, func(value []byte) bool {
		result = math.Float64frombits(p.ifd.tiff.Order.Uint64(value))
		return true
	})
	return
}

// Value returns value in its original format.
func (p *DE) Value() (result interface{}, err error) {
	switch p.ValueType {
	case TypeBYTE, TypeSBYTE, TypeUNDEFINED:
		result, err = p.BytesValue()
	case TypeASCII:
		result, err = p.StringValues()
	case TypeSHORT:
		if p.ValueCount == 1 {
			result, err = p.ShortValue()
		} else {
			result, err = p.ShortValues()
		}
	case TypeSSHORT:
		if p.ValueCount == 1 {
			result, err = p.SShortValue()
		} else {
			result, err = p.SShortValues()
		}
	case TypeLONG:
		if p.ValueCount == 1 {
			result, err = p.LongValue()
		} else {
			result, err = p.LongValues()
		}
	case TypeSLONG:
		if p.ValueCount == 1 {
			result, err = p.SLongValue()
		} else {
			result, err = p.SLongValues()
		}
	case TypeFLOAT:
		if p.ValueCount == 1 {
			result, err = p.FloatValue()
		} else {
			result, err = p.FloatValues()
		}
	case TypeRATIONAL:
		if p.ValueCount == 1 {
			result, err = p.RationalValue()
		} else {
			result, err = p.RationalValues()
		}
	case TypeSRATIONAL:
		if p.ValueCount == 1 {
			result, err = p.SRationalValue()
		} else {
			result, err = p.SRationalValues()
		}
	case TypeDOUBLE:
		if p.ValueCount == 1 {
			result, err = p.DoubleValue()
		} else {
			result, err = p.DoubleValues()
		}
	}
	return
}

package isobmff_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/levinholsety/common-go/isobmff"
)

func Test_Video(t *testing.T) {
	dir := `D:\Temp\Videos`
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, fileInfo := range fileInfos {
		func() {
			file, err := os.Open(filepath.Join(dir, fileInfo.Name()))
			if err != nil {
				panic(err)
			}
			defer file.Close()
			fmt.Println(readVideoInfo(file))
		}()
	}

}

func Test_HEIC(t *testing.T) {
	filename := `d:\Temp\test-rename\20180513_185604_iPhone 8 Plus_451.heic`
	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	fmt.Println(readHEICInfo(file))
}

var unixTime1904 = time.Date(1904, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

func readVideoInfo(file *os.File) (creationTime time.Time, duration uint, err error) {
	boxFile, err := isobmff.NewBoxFile(file)
	if err != nil {
		return
	}
	var box isobmff.Box
	if box, err = boxFile.ReadBox(); err != nil {
		return
	}
	if ftyp, ok := box.(*isobmff.FileTypeBox); ok &&
		(ftyp.MajorBrand == "isom" || ftyp.MajorBrand == "qt  ") {
		for {
			if box, err = boxFile.ReadBox(); err != nil {
				return
			}
			if moov, ok := box.(*isobmff.MovieBox); ok {
				for {
					if box, err = moov.ReadBox(); err != nil {
						return
					}
					var mvhd *isobmff.MovieHeaderBox
					if mvhd, ok = box.(*isobmff.MovieHeaderBox); ok {
						creationTime, duration = mvhdVideoInfo(mvhd)
						return
					}
				}
			}
		}
	}
	err = fmt.Errorf("video info not found")
	return
}

func mvhdVideoInfo(mvhd *isobmff.MovieHeaderBox) (creationTime time.Time, duration uint) {
	if mvhd.CreationTime != 0 {
		creationTime = time.Unix(unixTime1904+int64(mvhd.CreationTime), 0)
	} else if mvhd.ModificationTime != 0 {
		creationTime = time.Unix(unixTime1904+int64(mvhd.ModificationTime), 0)
	}
	duration = mvhd.Duration / mvhd.Timescale
	return
}

func readHEICInfo(file *os.File) (extentOffset, extentLength uint, ok bool, err error) {
	boxFile, err := isobmff.NewBoxFile(file)
	if err != nil {
		return
	}
	var box isobmff.Box
	if box, err = boxFile.ReadBox(); err != nil {
		return
	}
	var ftyp *isobmff.FileTypeBox
	if ftyp, ok = box.(*isobmff.FileTypeBox); ok && (ftyp.MajorBrand == "heic") {
		if box, err = boxFile.ReadBox(); err != nil {
			return
		}
		var meta *isobmff.MetaBox
		if meta, ok = box.(*isobmff.MetaBox); ok {
			var itemID uint32
			for {
				if box, err = meta.ReadBox(); err != nil {
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
							extentOffset = extent.ExtentOffset
							extentLength = extent.ExtentLength
							return
						}
					}
				}
			}
		}
	}
	ok = false
	return
}

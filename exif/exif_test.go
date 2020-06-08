package exif_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/levinholsety/common-go/exif"
)

func TestEXIF(t *testing.T) {
	dir := `d:\Temp\test-rename\`
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	for _, fileInfo := range fileInfos {
		ext := strings.ToLower(filepath.Ext(fileInfo.Name()))
		if ext == ".jpg" || ext == ".heic" || ext == ".nef" {
			if err = readExifInfo(filepath.Join(dir, fileInfo.Name())); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func readExifInfo(filename string) (err error) {
	fmt.Println(filename)
	var file *os.File
	if file, err = os.Open(filename); err != nil {
		return
	}
	defer file.Close()
	exifInfo, err := exif.Parse(file)
	if err != nil {
		return
	}
	jsonData, _ := json.MarshalIndent(data{
		Make:                value(exifInfo.Make()),
		Model:               value(exifInfo.Model()),
		DateTime:            value(exifInfo.DateTime()),
		SubSecTime:          value(exifInfo.SubsecTime()),
		DateTimeOriginal:    value(exifInfo.DateTimeOriginal()),
		SubSecTimeOriginal:  value(exifInfo.SubSecTimeOriginal()),
		DateTimeDigitized:   value(exifInfo.DateTimeDigitized()),
		SubSecTimeDigitized: value(exifInfo.SubSecTimeDigitized()),
		ExposureTime:        value(exifInfo.ExposureTime()),
		FNumber:             value(exifInfo.FNumber()),
		GPSLatitude:         value(exifInfo.GPSLatitude()),
		GPSLongitude:        value(exifInfo.GPSLongitude()),
		ShutterCount:        value(exifInfo.ShutterCount()),
		Lens:                value(exifInfo.Lens()),
		CanonImageNumber:    value(exifInfo.CanonImageNumber()),
	}, "", "    ")
	fmt.Println(string(jsonData))
	return
}

func value(v interface{}, err error) interface{} {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		debug.PrintStack()
	}
	return v
}

type data struct {
	Make                interface{} `json:"make,omitempty"`
	Model               interface{} `json:"model,omitempty"`
	DateTime            interface{} `json:"dateTime,omitempty"`
	SubSecTime          interface{} `json:"subSecTime,omitempty"`
	DateTimeOriginal    interface{} `json:"dateTimeOriginal,omitempty"`
	SubSecTimeOriginal  interface{} `json:"subSecTimeOriginal,omitempty"`
	DateTimeDigitized   interface{} `json:"dateTimeDigitized,omitempty"`
	SubSecTimeDigitized interface{} `json:"subSecTimeDigitized,omitempty"`
	ExposureTime        interface{} `json:"exposureTime,omitempty"`
	FNumber             interface{} `json:"fNumber,omitempty"`
	GPSLatitude         interface{} `json:"gpsLatitude,omitempty"`
	GPSLongitude        interface{} `json:"gpsLongitude,omitempty"`
	ShutterCount        interface{} `json:"shutterCount,omitempty"`
	Lens                interface{} `json:"lens,omitempty"`
	CanonImageNumber    interface{} `json:"canonImageNumber,omitempty"`
}

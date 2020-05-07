package utils

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// NewFileSaver creates and returns an instance of FileSaver.
func NewFileSaver() *FileSaver {
	return &FileSaver{
		MarshalMapping: map[string]func(v interface{}) ([]byte, error){
			".json": json.Marshal,
			".xml":  xml.Marshal,
		},
		PermissionMapping: map[string]os.FileMode{
			".sh": 0755,
		},
	}
}

// FileSaver can marshal content automatically and save it to file.
// It will marshal content with the function specified in MarshalMapping.
// The key of MarshalMapping is file extension and the value of it is marshal function.
// We can also specify file permission in PermissionMapping.
// The key of PermissionMapping is file extension too and the value of it is os.FileMode.
// If file permission is not specified, it will use os.FileMode(0644).
// By default, it already supports '.json', '.xml' and '.sh' files.
type FileSaver struct {
	MarshalMapping    map[string]func(v interface{}) ([]byte, error)
	PermissionMapping map[string]os.FileMode
}

// SaveFile saves content to file.
func (p *FileSaver) SaveFile(filename string, content interface{}) (err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	var data []byte
	if v, ok := content.([]byte); ok {
		data = v
	} else if v, ok := content.(string); ok {
		data = []byte(v)
	} else {
		if marshal, ok := p.MarshalMapping[ext]; ok {
			data, err = marshal(content)
			if err != nil {
				return
			}
		} else {
			err = errors.New("cannot marshal content")
			return
		}
	}
	var perm os.FileMode
	if v, ok := p.PermissionMapping[ext]; ok {
		perm = v
	} else {
		perm = 0644
	}
	return ioutil.WriteFile(filename, data, perm)
}

// SaveFiles saves each content to each file in the map.
func (p *FileSaver) SaveFiles(files map[string]interface{}) (err error) {
	for filename, content := range files {
		if err = p.SaveFile(filename, content); err != nil {
			return
		}
	}
	return
}

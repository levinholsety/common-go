package highlight_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/levinholsety/common-go/highlight"
)

func Test_highlight(t *testing.T) {
	data, _ := ioutil.ReadFile(`D:\a.sql`)
	result := highlight.Parse(string(data), highlight.SQLConfig)
	for _, text := range result {
		data, _ = json.Marshal(text)
		fmt.Println(string(data))
	}
}

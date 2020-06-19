package highlight_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/highlight"
)

func Test_highlight(t *testing.T) {
	data, _ := json.Marshal(highlight.SQLConfig)
	cfg := &highlight.Config{}
	err := json.Unmarshal(data, cfg)
	assert.NoError(t, err)
	data, _ = ioutil.ReadFile(`D:\a.sql`)
	result := highlight.Parse(string(data), cfg)
	for _, text := range result {
		fmt.Println(text)
	}

}

func Benchmark_highlight(b *testing.B) {
	b.StopTimer()
	data, _ := ioutil.ReadFile(`D:\a.sql`)
	text := string(data)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		highlight.Parse(text, highlight.SQLConfig)
	}
}

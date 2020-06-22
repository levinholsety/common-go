package dbutil_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/dbutil"
)

func Test_Split(t *testing.T) {
	data, err := ioutil.ReadFile(`D:\a.sql`)
	assert.NoError(t, err)
	result := dbutil.Split(string(data))
	for i, str := range result {
		fmt.Println(i, str)
	}
}

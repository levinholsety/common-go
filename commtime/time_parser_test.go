package commtime_test

import (
	"testing"

	"github.com/levinholsety/common-go/assert"
	"github.com/levinholsety/common-go/commtime"
)

func Test_ParseTime(t *testing.T) {
	assert.StringEqual(t, "2019-11-01 13:02:03", parse("2019-11-1 13:2:3"))
	assert.StringEqual(t, "2019-11-01 13:02:03", parse("2019-11-01 13:02:03"))
	assert.StringEqual(t, "2019-11-01 00:00:00", parse("2019-11-01"))
	assert.StringEqual(t, "0000-01-01 13:02:03", parse("13:02:03"))
	assert.StringEqual(t, "2019-11-01 13:02:03", parse("20191101130203"))
	assert.StringEqual(t, "2019-11-01 00:00:00", parse("20191101"))
	assert.StringEqual(t, "0000-01-01 13:02:03", parse("130203"))
	assert.StringEqual(t, "2019-11-01 21:02:03", parse("2019-11-01T13:02:03Z"))
	assert.StringEqual(t, "2019-11-01 20:02:03", parse("2019-11-01T13:02:03+01:00"))
	assert.StringEqual(t, "2019-11-01 21:02:03", parse("20191101T130203Z"))
	assert.StringEqual(t, "2019-11-01 20:02:03", parse("20191101T130203+0100"))
	assert.StringEqual(t, "2019-11-01 20:02:03", parse("20191101T130203+01"))
}

func parse(value string) string {
	t, err := commtime.ParseTime(value)
	if err != nil {
		panic(err)
	}
	return t.Local().Format("2006-01-02 15:04:05")
}

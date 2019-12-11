package commtime

import (
	"fmt"
	"regexp"
	"time"
)

var timeParserMap = map[*regexp.Regexp]string{
	regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`):               "2006-1-2T15:4:5Z",
	regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}[+-]\d{2}:\d{2}$`): "2006-1-2T15:4:5-07:00",
	regexp.MustCompile(`^\d{8}T\d{6}Z$`):                                       "20060102T150405Z",
	regexp.MustCompile(`^\d{8}T\d{6}[+-]\d{4}$`):                               "20060102T150405-0700",
	regexp.MustCompile(`^\d{8}T\d{6}[+-]\d{2}$`):                               "20060102T150405-07",
}

var localTimeParserMap = map[*regexp.Regexp]string{
	regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}$`): "2006-1-2 15:4:5",
	regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`):                         "2006-1-2",
	regexp.MustCompile(`^\d{1,2}:\d{1,2}:\d{1,2}$`):                       "15:4:5",
	regexp.MustCompile(`^\d{14}$`):                                        "20060102150405",
	regexp.MustCompile(`^\d{8}$`):                                         "20060102",
	regexp.MustCompile(`^\d{6}$`):                                         "150405",
}

// ParseTime parses string value to time.Time.
func ParseTime(value string) (result time.Time, err error) {
	for re, layout := range localTimeParserMap {
		if re.MatchString(value) {
			result, err = time.ParseInLocation(layout, value, time.Local)
			return
		}
	}
	for re, layout := range timeParserMap {
		if re.MatchString(value) {
			result, err = time.Parse(layout, value)
			return
		}
	}
	err = fmt.Errorf("parse failed: %s", value)
	return
}

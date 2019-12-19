package gps_test

import (
	"fmt"
	"testing"

	"github.com/levinholsety/common-go/gps"
)

func TestGPS(t *testing.T) {
	p := gps.GeodeticPoint{31.993536111111112, 118.79508333333334}
	fmt.Println(p.String())
	fmt.Println(p.WGS84ToGCJ02().String())
	fmt.Println(p.WGS84ToGCJ02().GCJ02ToBD09().String())
	fmt.Println(p.WGS84ToGCJ02().GCJ02ToBD09().BD09ToGCJ02().String())
	fmt.Println(p.WGS84ToGCJ02().GCJ02ToBD09().BD09ToGCJ02().GCJ02ToWGS84().String())
}

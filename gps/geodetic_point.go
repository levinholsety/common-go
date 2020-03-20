package gps

import (
	"fmt"
	"math"
)

const (
	pi  = math.Pi
	xpi = math.Pi * 3000 / 180
	a   = 6378245
	ee  = 0.00669342162296594323
)

func lat(x, y float64) (ret float64) {
	ret = -100 + 2*x + 3*y + 0.2*y*y + 0.1*x*y + 0.2*math.Sqrt(math.Abs(x))
	ret += (20*math.Sin(6*x*pi) + 20*math.Sin(2*x*pi)) * 2 / 3
	ret += (20*math.Sin(y*pi) + 40*math.Sin(y/3*pi)) * 2 / 3
	ret += (160*math.Sin(y/12*pi) + 320*math.Sin(y*pi/30)) * 2 / 3
	return
}

func lon(x, y float64) (ret float64) {
	ret = 300 + x + 2*y + 0.1*x*x + 0.1*x*y + 0.1*math.Sqrt(math.Abs(x))
	ret += (20*math.Sin(6*x*pi) + 20*math.Sin(2*x*pi)) * 2 / 3
	ret += (20*math.Sin(x*pi) + 40*math.Sin(x/3*pi)) * 2 / 3
	ret += (150*math.Sin(x/12*pi) + 300*math.Sin(x/30*pi)) * 2 / 3
	return ret
}

// GeodeticPoint represents a geodetic point.
//
// Example:
//
//	p := gps.GeodeticPoint{31.993536111111112, 118.79508333333334}
//
// WGS－84原始坐标系，一般用国际GPS纪录仪记录下来的经纬度，通过GPS定位拿到的原始经纬度，
// Google和高德地图定位的的经纬度（国外）都是基于WGS－84坐标系的；但是在国内是不允许直接用WGS84坐标系标注的，必须经过加密后才能使用；
//
// GCJ－02坐标系，又名“火星坐标系”，是我国国测局独创的坐标体系，由WGS－84加密而成，
// 在国内，必须至少使用GCJ－02坐标系，或者使用在GCJ－02加密后再进行加密的坐标系，
// 如百度坐标系。高德和Google在国内都是使用GCJ－02坐标系，可以说，GCJ－02是国内最广泛使用的坐标系；
//
// 百度坐标系:bd-09，百度坐标系是在GCJ－02坐标系的基础上再次加密偏移后形成的坐标系，只适用于百度地图。
type GeodeticPoint struct {
	Latitude  float64
	Longitude float64
}

func (v GeodeticPoint) outsideChina() bool {
	return (v.Latitude < 0.8293 || v.Latitude > 55.8271) || (v.Longitude < 72.004 || v.Longitude > 137.8347)
}

// String outputs position text, such as:
//
//	118.795083,31.993536
func (v GeodeticPoint) String() string {
	return fmt.Sprintf("%f,%f", v.Longitude, v.Latitude)
}

// WGS84ToGCJ02 converts WGS84 coordinates to GCJ02 coordinates.
func (v GeodeticPoint) WGS84ToGCJ02() GeodeticPoint {
	if v.outsideChina() {
		return v
	}
	lat := lat(v.Longitude-105, v.Latitude-35)
	lon := lon(v.Longitude-105, v.Latitude-35)
	radLat := v.Latitude / 180 * pi
	magic := math.Sin(radLat)
	magic = 1 - ee*magic*magic
	sqrtMagic := math.Sqrt(magic)
	lat = (lat * 180) / ((a * (1 - ee)) / (magic * sqrtMagic) * pi)
	lon = (lon * 180) / (a / sqrtMagic * math.Cos(radLat) * pi)
	return GeodeticPoint{v.Latitude + lat, v.Longitude + lon}
}

// GCJ02ToBD09 converts GCJ02 coordinates to BD09 coordinates.
func (v GeodeticPoint) GCJ02ToBD09() GeodeticPoint {
	z := math.Sqrt(v.Latitude*v.Latitude+v.Longitude*v.Longitude) + 0.00002*math.Sin(v.Latitude*xpi)
	theta := math.Atan2(v.Latitude, v.Longitude) + 0.000003*math.Cos(v.Longitude*xpi)
	return GeodeticPoint{z*math.Sin(theta) + 0.006, z*math.Cos(theta) + 0.0065}
}

// BD09ToGCJ02 converts BD09 coordinates to GCJ02 coordinates.
func (v GeodeticPoint) BD09ToGCJ02() GeodeticPoint {
	lat := v.Latitude - 0.006
	lon := v.Longitude - 0.0065
	z := math.Sqrt(lat*lat+lon*lon) - 0.00002*math.Sin(lat*xpi)
	theta := math.Atan2(lat, lon) - 0.000003*math.Cos(lon*xpi)
	return GeodeticPoint{z * math.Sin(theta), z * math.Cos(theta)}
}

// GCJ02ToWGS84 converts GCJ02 coordinates to WGS84 coordinates.
func (v GeodeticPoint) GCJ02ToWGS84() GeodeticPoint {
	p := v.WGS84ToGCJ02()
	return GeodeticPoint{v.Latitude*2 - p.Latitude, v.Longitude*2 - p.Longitude}
}

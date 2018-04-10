package rtbus

import (
	"math"
)

func Distance(lat1, lng1, lat2, lng2 float64) int {
	var radLat1 = lat1 * math.Pi / 180.0
	var radLat2 = lat2 * math.Pi / 180.0
	var a = radLat1 - radLat2
	var b = lng1*math.Pi/180.0 - lng2*math.Pi/180.0
	var s = 2 * math.Asin(math.Sqrt(math.Pow(math.Sin(a/2), 2)+math.Cos(radLat1)*math.Cos(radLat2)*math.Pow(math.Sin(b/2), 2)))
	s = s * 6378.137
	return int(s * 1000)
}

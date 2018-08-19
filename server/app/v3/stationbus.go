package v3

import (
	"github.com/xuebing1110/rtbus/server/handlers"
)

func init() {
	router.GET("/station-buses/bylocation/:location", handlers.StationBusesByLocation)
	router.GET("/station-buses/byline/:city/:sn", handlers.StationBusesByStation)
}

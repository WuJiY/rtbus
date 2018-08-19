package v3

import "github.com/xuebing1110/rtbus/server/handlers"

func init() {
	router.GET("/lines/:city/:line/:dir", handlers.BusLineDir)
	router.GET("/lines/:city/:line/:dir/buses", handlers.BusLineDirBuses)
	router.GET("/search/:city/:keyword", handlers.SearchBusLineDir)

}

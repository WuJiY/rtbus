package handlers

import (
	//"log"

	"github.com/gin-gonic/gin"
	"github.com/xuebing1110/rtbus/pkg/apis"
)

func StationBusesByLocation(ctx *gin.Context) {
	loc := ctx.Param("location")

	lazy := true
	lazy_str, found := ctx.GetQuery("lazy")
	if !found || lazy_str == "False" || lazy_str == "false" || lazy_str == "0" {
		lazy = false
	}

	sbs, err := apis.ListLocalStationBuses(loc, lazy, 6, 3)
	if err != nil {
		sendBadResponse(ctx, 500, "PoiAroundSearchFailed", err.Error())
		return
	}

	sendResponse(ctx, sbs)
}

func StationBusesByStation(ctx *gin.Context) {
	city := ctx.Param("city")
	sn := ctx.Param("sn")

	sb, err := apis.GetStationBusesByStation(city, sn, 3)
	if err != nil {
		sendBadResponse(ctx, 500, "PoiTextSearchFailed", err.Error())
		return
	}

	sendResponse(ctx, sb)
}

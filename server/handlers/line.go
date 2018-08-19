package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/xuebing1110/rtbus/pkg/apis"
	"log"
	"strconv"
)

func BusLineDir(ctx *gin.Context) {
	city := ctx.Param("city")
	line := ctx.Param("line")
	dir := ctx.Param("dir")

	withBus := true
	bus, found := ctx.GetQuery("bus")
	if !found || bus == "False" || bus == "false" || bus == "0" {
		withBus = false
	}

	ld, err := apis.GetBusLineDir(city, line, dir, withBus)
	if err != nil {
		sendBadResponse(ctx, 500, "LoadLineDirFailed", err.Error())
		return
	}

	sendResponse(ctx, ld)
}

func BusLineDirBuses(ctx *gin.Context) {
	city := ctx.Param("city")
	line := ctx.Param("line")
	dir := ctx.Param("dir")

	order := 0
	order_str := ctx.Query("order")
	if order_str != "" {
		var err error
		order, err = strconv.Atoi(order_str)
		if err != nil {
			sendBadResponse(ctx, 400, "ParamParseFailed", "order param should be number:"+err.Error())
			return
		}
	}

	buses, err := apis.GetBusLineDirBuses(city, line, dir, order, 3)
	if err != nil {
		log.Printf("get %s %s[%s] running buses failed: %v", city, line, dir, err)
		sendBadResponse(ctx, 500, "LoadLineDirFailed", err.Error())
		return
	}

	sendResponse(ctx, buses)
}

func SearchBusLineDir(ctx *gin.Context) {
	city := ctx.Param("city")
	keyword := ctx.Param("keyword")

	//withBus := true
	//bus, found := ctx.GetQuery("bus")
	//if !found || bus == "False" || bus == "false" || bus == "0" {
	//	withBus = false
	//}

	blds, err := apis.SearchBusLineDir(city, keyword, false)
	if err != nil {
		log.Printf("search %s %s busline failed: %v", city, keyword, err)
		sendBadResponse(ctx, 500, "SearchLineDirFailed", err.Error())
		return
	}

	sendResponse(ctx, blds)

}

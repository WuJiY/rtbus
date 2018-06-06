package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

func BusLineSearch(params martini.Params, r render.Render) {
	cityid := params["city"]
	keyword := params["keyword"]

	bids, err := RTBusClient.Search(cityid, keyword)
	if err != nil {
		r.JSON(
			502,
			&Response{502, err.Error(), nil},
		)
		return
	}

	bids_resp := make([]rtbus.BusDirInfo, len(bids))
	for i, bid := range bids {
		bids_resp[i] = *bid
		bids_resp[i].Stations = nil
		bids_resp[i].RunningBuses = nil
	}

	r.JSON(200,
		&Response{
			0,
			"OK",
			bids_resp,
		},
	)
}

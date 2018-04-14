package handler

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
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

	r.JSON(200,
		&Response{
			0,
			"OK",
			bids,
		},
	)
}

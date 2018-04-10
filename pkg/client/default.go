package client

import (
	"github.com/xuebing1110/rtbus/pkg/rtbus/beijing"
	"github.com/xuebing1110/rtbus/pkg/rtbus/chell"
)

var (
	DefaultRTBus *RTBus = newRTBus()
)

func init() {
	// chell
	citys, err := chell.GetAllCitys()
	if err != nil {
		panic(err)
	}
	for _, city := range citys {
		DefaultRTBus.Register(chell.NewChellRTBusApi(city))
	}

	// beijing
	cba, err := beijing.NewBJRTBusApi()
	if err != nil {
		panic(err)
	}
	DefaultRTBus.MustRegister(cba)
}

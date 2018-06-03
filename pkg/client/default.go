package client

import (
	"github.com/xuebing1110/rtbus/pkg/rtbus/beijing"
	"github.com/xuebing1110/rtbus/pkg/rtbus/chell"
	"github.com/xuebing1110/rtbus/pkg/rtbus/guangzhou"
	"github.com/xuebing1110/rtbus/pkg/rtbus/jinan"
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

	// jinan
	cba_jinan := jinan.NewJinanRTBusApi()
	DefaultRTBus.MustRegister(cba_jinan)

	// guangzhou
	cba_guangzhou := guangzhou.NewRTBusApi()
	DefaultRTBus.MustRegister(cba_guangzhou)
}

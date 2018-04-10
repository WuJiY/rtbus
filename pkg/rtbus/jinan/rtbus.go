package jinan

import (
	"fmt"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

type JinanRTBusApi struct {
	city *rtbus.CityInfo
}

func NewJinanRTBusApi() *JinanRTBusApi {
	return &JinanRTBusApi{
		city: &rtbus.CityInfo{
			Code:   "0531",
			ID:     "0531",
			Name:   "济南",
			Hot:    1,
			PinYin: "jinan",
			Subway: 0,
		},
	}
}

func (cba *JinanRTBusApi) City() *rtbus.CityInfo {
	return cba.city
}
func (cba *JinanRTBusApi) GetBusLine(lineno string) (*rtbus.BusLine, error) {
	bl, err := getBusLine(lineno)
	if err != nil {
		return nil, err
	}

	for _, bdi := range bl.Directions {
		bdi.RunningBuses, err = getRunningBus(bdi)
		if err != nil {
			return bl, err
		}
	}
	return bl, nil
}
func (cba *JinanRTBusApi) GetBusLineDir(lineno, dirname string) (*rtbus.BusDirInfo, error) {
	bl, err := cba.GetBusLine(lineno)
	if err != nil {
		return nil, err
	}

	bdi, found := bl.GetBusDirInfo(dirname)
	if !found {
		return nil, fmt.Errorf("not found")
	}

	bdi.RunningBuses, err = getRunningBus(bdi)
	return bdi, err
}

func (cba *JinanRTBusApi) GetRunningBus(lineno, dirname string) (rbus []*rtbus.RunningBus, err error) {
	bl, err := cba.GetBusLine(lineno)
	if err != nil {
		return rbus, err
	}

	bdi, found := bl.GetBusDirInfo(dirname)
	if !found {
		return rbus, fmt.Errorf("not found")
	}

	return getRunningBus(bdi)
}

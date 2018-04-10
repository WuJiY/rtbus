package chell

import (
	"fmt"

	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

var (
	LOGGER = logs.GetBlogger()
)

type ChellRTBusApi struct {
	city *rtbus.CityInfo
}

func NewChellRTBusApi(city *rtbus.CityInfo) *ChellRTBusApi {
	return &ChellRTBusApi{
		city: city,
	}
}

func (cba *ChellRTBusApi) City() *rtbus.CityInfo {
	return cba.city
}
func (cba *ChellRTBusApi) GetBusLine(lineno string) (*rtbus.BusLine, error) {
	return loadBusline(cba.city.ID, lineno)
}
func (cba *ChellRTBusApi) GetBusLineDir(lineno, dirname string) (*rtbus.BusDirInfo, error) {
	bl, err := cba.GetBusLine(lineno)
	if err != nil {
		return nil, err
	}

	bdi, found := bl.GetBusDirInfo(dirname)
	if !found {
		return nil, fmt.Errorf("can't found the direction %s of line %s in %s", dirname, lineno, cba.city.Name)
	}

	return bdi, nil
}

func (cba *ChellRTBusApi) GetRunningBus(lineno, dirname string) (rbus []*rtbus.RunningBus, err error) {
	bdi, err := cba.GetBusLineDir(lineno, dirname)
	if err != nil {
		return rbus, err
	}

	return bdi.RunningBuses, nil
}

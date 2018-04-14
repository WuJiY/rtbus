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
func (cba *ChellRTBusApi) Search(keyword string) (bdis []*rtbus.BusDirInfo, err error) {
	csls, err := search(cba.city.ID, keyword)
	if err != nil {
		return bdis, err
	}

	bdis = make([]*rtbus.BusDirInfo, 0, len(csls)*2)
	for _, csl := range csls {
		bdi, err := getNewestCllBusDirInfo(cba.City().ID, csl.LineId, csl.LineNo)
		if err == nil {
			bdis = append(bdis, bdi)

			// other directions
			if len(bdi.OtherDirIDs) > 0 {
				for _, dir_id := range bdi.OtherDirIDs {
					bdi_other, err_other := getNewestCllBusDirInfo(cba.City().ID, dir_id, csl.LineNo)
					if err_other == nil {
						bdis = append(bdis, bdi_other)
					}
				}
			}
		} else {

		}
	}
	return bdis, nil
}
func (cba *ChellRTBusApi) GetBusLine(lineno string, with_running_bus bool) (*rtbus.BusLine, error) {
	return loadBusline(cba.city.ID, lineno)
}
func (cba *ChellRTBusApi) GetBusLineDir(lineno, dirname string) (*rtbus.BusDirInfo, error) {
	bl, err := cba.GetBusLine(lineno, true)
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

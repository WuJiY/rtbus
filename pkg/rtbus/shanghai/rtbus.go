package shanghai

import (
	"fmt"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
	"sort"
)

type RTBusApi struct {
	city  *rtbus.CityInfo
	lines sort.StringSlice
}

func NewRTBusApi() (*RTBusApi, error) {
	rtb := &RTBusApi{
		city: &rtbus.CityInfo{
			Code:   "021",
			ID:     "021",
			Name:   "上海",
			Hot:    1,
			PinYin: "shanghai",
			Subway: 1,
		},
		lines: slines,
	}

	return rtb, nil
}

func (cba *RTBusApi) City() *rtbus.CityInfo {
	return cba.city
}

func (cba *RTBusApi) Search(keyword string) (bdis []*rtbus.BusDirInfo, err error) {
	names := linePrefixSearch(keyword)

	bdis = make([]*rtbus.BusDirInfo, 0, len(names)*2)
	for _, name := range names {
		bl, err := getBusLine(name)
		if err != nil {
			fmt.Printf("error:%v", err)
			continue
		}
		for _, bdi := range bl.Directions {
			bdis = append(bdis, bdi)
		}
	}
	return
}

func (cba *RTBusApi) GetBusLine(lineno string, with_running_bus bool) (bl *rtbus.BusLine, err error) {
	bl, err = getBusLine(lineno)
	if with_running_bus {
		for _, bdi := range bl.Directions {
			bdi.RunningBuses, err = getRunningBus(bdi, -1)
			if err != nil {
				fmt.Printf("get %s %d error: %v", lineno, bdi.Direction, err)
				continue
			}
		}
	}
	return
}

func (cba *RTBusApi) GetBusLineDir(lineno, dirname string) (*rtbus.BusDirInfo, error) {
	return getBusLineDir(lineno, dirname)
}

func (cba *RTBusApi) GetRunningBus(lineno, dirname string) (rbus []*rtbus.RunningBus, err error) {
	bdi, err := getBusLineDir(lineno, dirname)
	if err != nil {
		return rbus, err
	}
	return getRunningBus(bdi, -1)
}

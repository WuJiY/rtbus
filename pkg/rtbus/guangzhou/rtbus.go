package guangzhou

import (
	"fmt"

	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

type RTBusApi struct {
	city *rtbus.CityInfo
}

func NewRTBusApi() *RTBusApi {
	return &RTBusApi{
		city: &rtbus.CityInfo{
			Code:   "020",
			ID:     "020",
			Name:   "广州",
			Hot:    1,
			PinYin: "guangzhou",
			Subway: 0,
		},
	}
}

func (cba *RTBusApi) City() *rtbus.CityInfo {
	return cba.city
}
func (cba *RTBusApi) Search(keyword string) (bdis []*rtbus.BusDirInfo, err error) {
	lines, err := searchBusLine(keyword)
	if err != nil {
		return bdis, err
	}

	bdis = make([]*rtbus.BusDirInfo, 0, len(lines)*2)
	for _, line := range lines {
		for i := 0; i < 2; i++ {
			bdi, err := getBusLineDirByRouteId(line.Name, line.RouteId, i)
			if err == nil {
				bdis = append(bdis, bdi)
			}
		}
	}

	return bdis, nil
}
func (cba *RTBusApi) GetBusLine(lineno string, with_running_bus bool) (*rtbus.BusLine, error) {
	bl := &rtbus.BusLine{
		LineNum:    lineno,
		LineName:   lineno,
		Directions: make(map[string]*rtbus.BusDirInfo),
	}

	routeid := GetBusLineRouteId(lineno)
	for i := 0; i < 2; i++ {
		bdi, err := getBusLineDirByRouteId(lineno, routeid, i)
		if err != nil {
			return nil, err
		}
		bdi.Direction = i
		//bdi.Name = bdi.GetDirName()
		bdi.Name = lineno

		if i == 0 {
			bdi.OtherDirIDs = []string{"1"}
		} else {
			bdi.OtherDirIDs = []string{"0"}
		}
		bl.PutDir(bdi)
	}

	return bl, nil
}
func (cba *RTBusApi) GetBusLineDir(lineno, dirname string) (*rtbus.BusDirInfo, error) {
	bl, err := cba.GetBusLine(lineno, true)
	if err != nil {
		return nil, err
	}

	bdi, found := bl.GetBusDirInfo(dirname)
	if !found {
		return nil, fmt.Errorf("not found")
	}
	return bdi, nil
}

func (cba *RTBusApi) GetRunningBus(lineno, dirname string) (rbus []*rtbus.RunningBus, err error) {
	bdi, err := cba.GetBusLineDir(lineno, dirname)
	if err != nil {
		return rbus, err
	}

	return bdi.RunningBuses, nil
}

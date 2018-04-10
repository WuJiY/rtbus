package beijing

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

var (
	LOGGER = logs.GetBlogger()
)

const (
	BTIC_BUS_PATH = `/ssgj/bus.php`
)

type BJRTBusApi struct {
	city        *rtbus.CityInfo
	byShortName *sync.Map
}

func NewBJRTBusApi() (*BJRTBusApi, error) {
	bj_api := &BJRTBusApi{
		city: &rtbus.CityInfo{
			Code:   "010",
			ID:     "010",
			Name:   "北京",
			Hot:    1,
			PinYin: "beijing",
			Subway: 0,
		},
		byShortName: new(sync.Map),
	}

	version := "0"
	bus_dir_chan, err := LoadBusLines(version)
	if err != nil {
		return nil, err
	}

	for bus_dir := range bus_dir_chan {
		//LOGGER.Info("%+v", bus_dir)

		v, found := bj_api.byShortName.Load(bus_dir.Name)
		if found {
			v.(*rtbus.BusLine).Directions[bus_dir.GetDirName()] = bus_dir
		} else {
			bj_api.byShortName.Store(bus_dir.Name, newBusLineByABLine(bus_dir))
		}
	}

	return bj_api, nil
}

func (brt *BJRTBusApi) GetBusLine(lineno string) (bl *rtbus.BusLine, err error) {
	bus_line, found := brt.byShortName.Load(lineno)
	if !found {
		return nil, fmt.Errorf("can't found the line %s in beijing", lineno)
	} else {
		bl = bus_line.(*rtbus.BusLine)
		// refresh running bus
		for _, bdi := range bl.Directions {
			//fmt.Printf("%+v", bdi)
			bdi.RunningBuses, err = brt.GetRunningBus(lineno, bdi.ID)
			if err != nil {
				return bl, err
			}
		}
		return bl, nil
	}
}

func (brt *BJRTBusApi) GetBusLineDir(lineno, dirname string) (bdi *rtbus.BusDirInfo, err error) {
	bdi, err = brt.getBusLineDir(lineno, dirname)
	if err != nil {
		return
	}

	// refresh running bus
	bdi.RunningBuses, err = brt.GetRunningBus(lineno, bdi.ID)
	if err.Error() == "获取数据失败" {
		return bdi, nil
	}
	return
}

func (brt *BJRTBusApi) getBusLineDir(lineno, dirname string) (bdi *rtbus.BusDirInfo, err error) {
	bus_line, found := brt.byShortName.Load(lineno)
	if !found {
		return nil, fmt.Errorf("can't found the line %s in beijing", lineno)
	} else {
		bl := bus_line.(*rtbus.BusLine)
		bdi, found = bl.GetBusDirInfo(dirname)
		if found {
			return bdi, nil
		} else {
			return nil, fmt.Errorf("can't found the direction %s of line %s in beijing", dirname, lineno)
		}
	}
}

func (brt *BJRTBusApi) City() *rtbus.CityInfo {
	return brt.city
}

func (brt *BJRTBusApi) GetRunningBus(lineno, dirname string) (rbus []*rtbus.RunningBus, err error) {
	bdi, err := brt.getBusLineDir(lineno, dirname)
	if err != nil {
		return rbus, err
	}

	// BTIC_BUS_URL_PARAM_FMT = `versionid=5&encrypt=1&datatype=json&no=1&type=0&id=%s&city=%E5%8C%97%E4%BA%AC`
	curtime := time.Now().Unix()
	var params *url.Values = &url.Values{}
	params.Set("city", "北京")
	params.Set("id", bdi.ID)
	params.Set("no", "1")
	params.Set("type", "0")
	params.Set("datatype", "json")
	params.Set("encrypt", "1")
	params.Set("versionid", "5")

	linert_resp := &BticLineRTResp{}
	err = bticRequest(BTIC_BUS_PATH, params, linert_resp)
	if err != nil {
		return
	}
	linert := linert_resp.Root
	if linert == nil {
		//LOGGER.Error("can't found root from %s", ToJsonString(linert_resp))
		err = errors.New("can't found root")
		return
	} else if linert.Status != 200 {
		//LOGGER.Error("%s", ToJsonString(linert))
		err = errors.New(linert.Message)
		return
	} else {

	}

	rbus = make([]*rtbus.RunningBus, len(linert.Data.Bus))
	for i, bus := range linert.Data.Bus {
		//fmt.Printf("bus:%+v\n", bus)
		rawkey := fmt.Sprintf(BTIC_RAW_KEY_FMT, bus.GT)
		bus.Lat, _ = Rc4DecodeString(rawkey, bus.Lat)
		bus.Lon, _ = Rc4DecodeString(rawkey, bus.Lon)
		bus.NextStationName, _ = Rc4DecodeString(rawkey, bus.NextStationName)
		bus.NextStationNum, _ = Rc4DecodeString(rawkey, bus.NextStationNum)
		//bus.NextStationDistance, _ = Rc4DecodeString(rawkey, bus.NextStationDistance)
		//bus.NextStationArrTime, _ = Rc4DecodeString(rawkey, bus.NextStationArrTime)
		bus.StationDistance, _ = Rc4DecodeString(rawkey, bus.StationDistance)
		bus.StationArrTime, _ = Rc4DecodeString(rawkey, bus.StationArrTime)

		// LOGGER.Warn(ToJsonString(bus))

		rbus[i] = &rtbus.RunningBus{Name: bus.NextStationName, BusID: bus.ID, SyncTime: curtime}
		rbus[i].No, _ = strconv.Atoi(bus.NextStationNum)
		rbus[i].Lat, _ = strconv.ParseFloat(bus.Lat, 10)
		rbus[i].Lng, _ = strconv.ParseFloat(bus.Lon, 10)

		if bus.NextStationDistance == "0" || bus.NextStationDistance == "-1" {
			rbus[i].Status = rtbus.BUS_ARRIVING_STATUS
			rbus[i].Distance = 0
		} else {
			rbus[i].Status = rtbus.BUS_ARRIVING_FUTURE_STATUS
			rbus[i].Distance, _ = strconv.Atoi(bus.NextStationDistance)
		}
	}

	sort.Sort(SortRunningBus(rbus))
	return
}

// SortRunningBus
type SortRunningBus []*rtbus.RunningBus

func (rb SortRunningBus) Len() int      { return len(rb) }
func (rb SortRunningBus) Swap(i, j int) { rb[i], rb[j] = rb[j], rb[i] }
func (rb SortRunningBus) Less(i, j int) bool {
	if rb[i].No == rb[j].No {
		return rb[i].Status < rb[j].Status
	} else {
		return rb[i].No < rb[j].No
	}
}

package jinan

import (
	"fmt"
	"time"

	"github.com/xuebing1110/location"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
	"sort"
)

const (
	FMT_URL_SUGGEST    = "/server-ue2/rest/buslines/simple/370100/%s/0/20"
	FMT_URL_BUSLINE    = "/server-ue2/rest/buslines/370100/%s"
	FMT_URL_RUNNINGBUS = "/server-ue2/rest/buses/busline/370100/%s"
)

func searchBusLine(keyword string) ([]*SearchResponseResult, error) {
	path := fmt.Sprintf(FMT_URL_SUGGEST, keyword)
	resp := new(SearchResponse)

	err := doRequest(path, resp)
	if err != nil {
		return []*SearchResponseResult{}, err
	}

	return resp.Result.Result, nil
}

func searchSpecBusline(lineno string) ([]*SearchResponseResult, error) {
	results, err := searchBusLine(lineno)
	if err != nil {
		return []*SearchResponseResult{}, err
	}

	ssr := make([]*SearchResponseResult, 0, 2)
	for _, ret := range results {
		if ret.LineName == lineno {
			ssr = append(ssr, ret)
			if len(ssr) == 2 {
				break
			}
		}
	}

	if len(ssr) == 0 {
		return []*SearchResponseResult{}, fmt.Errorf("not found")
	}

	return ssr, nil
}

func getBusLine(lineno string) (*rtbus.BusLine, error) {
	ssr, err := searchSpecBusline(lineno)
	if err != nil {
		return nil, err
	}

	bl := &rtbus.BusLine{
		LineNum:    lineno,
		LineName:   lineno,
		Directions: make(map[string]*rtbus.BusDirInfo),
	}

	for i, ret := range ssr {
		bdi, err := getBusLineDirByLineid(ret.ID)
		if err != nil {
			return nil, err
		}
		bdi.Direction = i
		//bdi.Name = bdi.GetDirName()
		bdi.Name = lineno

		if i == 0 {
			if len(ssr) == 2 {
				bdi.OtherDirIDs = []string{"1"}
			}
		} else {
			bdi.OtherDirIDs = []string{"0"}
		}
		bl.PutDir(bdi)
	}

	return bl, nil
}

func getBusLineDirByLineid(lineid string) (*rtbus.BusDirInfo, error) {
	path := fmt.Sprintf(FMT_URL_BUSLINE, lineid)
	resp := new(BusLineResponse)
	err := doRequest(path, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status.Code != 0 {
		return nil, fmt.Errorf(resp.Status.Msg)
	}

	stations := make([]*rtbus.BusStation, 0, len(resp.Result.Stations))
	for i, station := range resp.Result.Stations {
		lat, lng := location.BdDencrypt(station.Lat, station.Lng)
		stations = append(stations, &rtbus.BusStation{
			No:   i + 1,
			Name: station.StationName,
			Lat:  lat,
			Lon:  lng,
		})
	}

	bdi := &rtbus.BusDirInfo{
		ID:       lineid,
		StartSn:  resp.Result.StartStationName,
		EndSn:    resp.Result.EndStationName,
		Price:    resp.Result.TicketPrice,
		SnNum:    len(resp.Result.Stations),
		Stations: stations,
	}
	return bdi, nil
}

func getRunningBus(bdi *rtbus.BusDirInfo) ([]*rtbus.RunningBus, error) {
	rbuses, err := getRunningBusByID(bdi.ID)
	if err != nil {
		return rbuses, err
	}

	for _, rbus := range rbuses {
		if rbus.No > len(bdi.Stations) {
			rbus.No = len(bdi.Stations)
		}
		station := bdi.Stations[rbus.No-1]
		rbus.Distance = int(location.Distance(station.Lat, station.Lon, rbus.Lat, rbus.Lng))
		if rbus.Distance < 30 {
			rbus.Status = rtbus.BUS_ARRIVING_STATUS
		}
		//fmt.Printf("%f,%f <-> %f,%f : %d\n", station.Lat, station.Lon, rbus.Lat, rbus.Lng, rbus.Distance)
	}

	// sort
	rbs := rtbus.RunBuses(rbuses)
	sort.Sort(rbs)

	return []*rtbus.RunningBus(rbs), nil
}

func getRunningBusByID(lineid string) ([]*rtbus.RunningBus, error) {
	path := fmt.Sprintf(FMT_URL_RUNNINGBUS, lineid)
	resp := new(RunningBusResponse)
	err := doRequest(path, resp)
	if err != nil {
		return nil, err
	}
	if resp.Status.Code != 0 {
		return nil, fmt.Errorf(resp.Status.Msg)
	}

	running_buses := make([]*rtbus.RunningBus, 0, len(resp.Result))
	for _, bus := range resp.Result {
		running_bus := &rtbus.RunningBus{
			No:       bus.StationSeqNum + 1,
			Status:   rtbus.BUS_ARRIVING_FUTURE_STATUS,
			BusID:    bus.BusID,
			Lat:      bus.Lat,
			Lng:      bus.Lng,
			Distance: -1,
		}

		update_time, err := time.ParseInLocation("Jan 02, 2006 03:04:05 PM", bus.ActTime, time.Local)
		if err == nil {
			running_bus.SyncTime = update_time.Unix()
		}

		running_buses = append(running_buses, running_bus)
	}
	return running_buses, nil
}

package guangzhou

import (
	"fmt"
	"time"

	"github.com/xuebing1110/location"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

const (
	FMT_URL_BUSLINE = "/xxt-min-api/bus/route/getByRouteIdAndDirection.do?direction=%d&routeId=%s"
	FMT_URL_SUGGEST = "/xxt-min-api/bus/getByName.do?name=%s"
)

func searchBusLine(keyword string) ([]Line, error) {
	path := fmt.Sprintf(FMT_URL_SUGGEST, keyword)
	sr := new(SearchResult)
	resp := &Response{Result: sr}

	err := doRequest(path, resp)
	if err != nil {
		return []Line{}, err
	}

	err = checkResponse(resp)
	if err != nil {
		return []Line{}, err
	}

	return sr.Routes, nil
}

func GetBusLineRouteId(lineid string) string {
	var routeid string
	lines, err := searchBusLine(lineid)
	if err != nil {
		return ""
	}
	if len(lines) == 0 {
		return ""
	}

	for _, line := range lines {
		if line.Name == lineid || line.Name == lineid+"路" {
			routeid = line.RouteId
			break
		}
	}
	if routeid == "" {
		routeid = lines[0].RouteId
	}

	return routeid
}

func getBusLineDirByRouteId(lineid, routeid string, dirid int) (*rtbus.BusDirInfo, error) {
	path := fmt.Sprintf(FMT_URL_BUSLINE, dirid, routeid)
	br := new(BusLineResult)
	resp := &Response{Result: br}
	err := doRequest(path, resp)
	if err != nil {
		return nil, err
	}
	err = checkResponse(resp)
	if err != nil {
		return nil, err
	}
	if br.BusLine == nil {
		return nil, fmt.Errorf("read rb failed")
	}

	stations := make([]*rtbus.BusStation, 0, len(br.BusLine.Stations))
	for i, station := range br.BusLine.Stations {
		// lat, lng := location.BdDencrypt(station.Lat, station.Lon)
		lat, lng := location.GCJEncrypt(station.Lat, station.Lon)
		// lat, lng := station.Lat, station.Lon
		stations = append(stations, &rtbus.BusStation{
			No:   i + 1,
			Name: station.Name,
			Lat:  lat,
			Lon:  lng,
		})
	}
	if len(stations) == 0 {
		return nil, fmt.Errorf("read stations failed")
	}

	cur_time := time.Now().Unix()
	rbs := make([]*rtbus.RunningBus, 0)
	for i, b := range br.RunningBusInfo {
		for _, bl_array := range [][]RunningBus{b.Bl, b.Bbl} {
			for _, bl := range bl_array {
				lat, lng := location.GCJEncrypt(bl.Lat, bl.Lon)
				distance := int(location.Distance(lat, lng, stations[i].Lat, stations[i].Lon))

				var no int
				if i == 0 {
					no = 1
				} else {
					no = i + 2
				}
				if no > len(stations) {
					no = len(stations)
				}

				rbus := &rtbus.RunningBus{
					No:       no,
					Name:     stations[no-1].Name,
					Status:   rtbus.BUS_ARRIVING_FUTURE_STATUS,
					BusID:    bl.No,
					Lat:      lat,
					Lng:      lng,
					Distance: distance,
					SyncTime: cur_time,
				}
				if rbus.Distance < 30 {
					rbus.Status = rtbus.BUS_ARRIVING_STATUS
				}
				rbs = append(rbs, rbus)
			}
		}
	}

	bdi := &rtbus.BusDirInfo{
		ID:           lineid,
		Name:         lineid,
		StartSn:      br.BusLine.Stations[0].Name,
		EndSn:        br.BusLine.Stations[len(br.BusLine.Stations)-1].Name,
		Price:        "",
		SnNum:        len(br.BusLine.Stations),
		Stations:     stations,
		RunningBuses: rbs,
	}
	return bdi, nil
}

func getBusLineDirByLineid(lineid string, dirid int) (*rtbus.BusDirInfo, error) {
	routeid := GetBusLineRouteId(lineid)
	if routeid == "" {
		return nil, fmt.Errorf("%s not found", lineid)
	}

	return getBusLineDirByRouteId(lineid, routeid, dirid)
}

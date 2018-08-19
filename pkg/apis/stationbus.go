package apis

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/xuebing1110/location/amap"
)

func GetStationBusesByStation(city, sn string, buslimit int) (sb *StationBus, err error) {
	resp, err2 := amap.NewPoiSearchRequest(amapClient, sn).
		SetPageSize(1).
		SetType("150700").
		SetCity(city).
		Do()
	if err2 != nil {
		err = err2
		return
	}

	if len(resp.Pois) == 0 {
		err = fmt.Errorf("search %s %s location failed", city, sn)
		return
	}

	poi := resp.Pois[0]
	lns := strings.Split(poi.Address, ";")
	sb = &StationBus{
		StationName: sn,
		Location:    poi.Location,
		LineCount:   len(lns),
	}
	sb.GetStationBus(poi.CityCode, lns, true, buslimit)

	return
}

func ListLocalStationBuses(loc string, lazy bool, stationSize int, buslimit int) (sbs []StationBus, err error) {
	resp, err2 := amap.NewPoiSearchRequest(amapClient, "").
		SetAroundSearch(loc).
		SetPageSize(stationSize).
		SetType("150700").
		Do()
	if err2 != nil {
		err = err2
		return
	}
	sbs = make([]StationBus, len(resp.Pois))

	var globalWg sync.WaitGroup
	for i, poi := range resp.Pois {
		globalWg.Add(1)

		sn := strings.TrimRight(poi.Name, "(公交站)")
		lns := strings.Split(poi.Address, ";")
		sbs[i] = StationBus{
			StationName: sn,
			Location:    poi.Location,
			LineCount:   len(lns),
		}

		withBuses := false
		if i == 0 {
			withBuses = true
		}
		if !lazy {
			withBuses = true
		}

		go func(sb *StationBus, city string, lns []string, withBuses bool) {
			defer globalWg.Done()

			sb.GetStationBus(poi.CityCode, lns, withBuses, buslimit)
		}(&sbs[i], poi.CityCode, lns, withBuses)
	}

	globalWg.Wait()

	return
}

func (sb *StationBus) GetStationBus(city string, lns []string, withBuses bool, buslimit int) {
	var wg sync.WaitGroup
	c := make(chan *BaseLineDir)
	for _, ln := range lns {
		wg.Add(1)

		go func(ln string) {
			defer wg.Done()

			ln = ParseLineName(ln)
			dir := "0"
			bl, err := getBaseLineDir(city, ln, dir, sb.StationName, withBuses, buslimit)
			if err != nil {
				log.Printf("read %s %s %s busline failed: %v", city, ln, dir, err)
				return
			}

			c <- bl
		}(ln)
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	sb.Lines = make([]BaseLineDir, 0, len(lns))
	for bl := range c {
		sb.Lines = append(sb.Lines, *bl)
		sb.SupportLineCount++
	}

	return
}

func getBaseLineDir(city, ln, dir, sn string, withBuses bool, buslimit int) (*BaseLineDir, error) {
	ld, err := getLineDir(city, ln, dir, withBuses)
	if err != nil {
		return nil, err
	}
	bl := ld.BaseLineDir

	// get station order && next station name
	if sn != "" {
		for i, s := range ld.Stations {
			if s.Sn == sn {
				bl.Order = s.Order
				if i == len(ld.Stations)-1 {
					bl.NextSn = sn
				} else {
					bl.NextSn = ld.Stations[i+1].Sn
				}
				break
			}
		}

		bl.Buses = limitBuses(bl.Buses, bl.Order, buslimit)
	}

	return bl, nil
}

func limitBuses(buses []Bus, order, limit int) []Bus {
	if limit <= 0 {
		limit = len(buses)
	}

	rbs := make([]Bus, 0, len(buses))
	for _, rbus := range buses {
		if rbus.Order <= order || order == 0 {
			rbs = append(rbs, rbus)
		} else {
			break
		}
	}

	ret := make([]Bus, 0, limit)
	for i := len(rbs) - 1; i >= 0; i-- {
		ret = append(ret, Bus{
			Order:    rbs[i].Order,
			Status:   rbs[i].Status,
			Location: rbs[i].Location,
			Distance: rbs[i].Distance,
		})

		if len(ret) >= limit {
			break
		}
	}
	return ret
}

func ParseLineName(n string) string {
	lineno_1 := strings.SplitN(n, `/`, -1)
	lineno := strings.TrimRight(lineno_1[0], "专线车")
	lineno = strings.TrimRight(lineno, "线")
	lineno = strings.TrimRight(lineno, "路")
	lineno = strings.Replace(lineno, "路内环", "内", 1)
	lineno = strings.Replace(lineno, "路外环", "外", 1)
	lineno = strings.Replace(lineno, "路大站快车", "大站", 1)
	lineno = strings.Replace(lineno, "路大站快", "大站", 1)
	lineno = strings.Replace(lineno, "路快线", "快线", 1)
	lineno = strings.Replace(lineno, "路", "", 1)

	return lineno
}

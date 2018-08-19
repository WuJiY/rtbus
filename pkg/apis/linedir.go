package apis

import (
	"fmt"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
	"log"
)

func GetBusLineDir(city, ln, dir string, withBuses bool) (*LineDirection, error) {
	return getLineDir(city, ln, dir, withBuses)
}

func GetBusLineDirBuses(city, ln, dir string, order, buslimit int) ([]Bus, error) {
	buses, err := getBuses(city, ln, dir)
	if err != nil {
		return buses, err
	}

	if order > 0 {
		buses = limitBuses(buses, order, buslimit)
	}
	return buses, nil
}

func SearchBusLineDir(city, keyword string, withBuses bool) ([]*LineDirection, error) {
	// read busline direction
	bdis, err := RTBusClient.Search(city, keyword)
	if err != nil {
		return nil, err
	}

	lds := make([]*LineDirection, 0, len(bdis))
	for _, bdi := range bdis {
		ld := toLineDirection(bdi)

		// running buses
		if withBuses {
			ld.BaseLineDir.Buses, err = getBuses(city, ld.No, ld.Dir)
			if err != nil {
				log.Printf("get %s %s[%s] running buses failed: %v", city, ld.No, ld.Dir, err)
			}
		}

		lds = append(lds, ld)
	}
	return lds, nil
}

func getLineDir(city, ln, dir string, withBuses bool) (*LineDirection, error) {
	// read busline direction
	bdi, err := RTBusClient.GetBusLineDir(city, ln, dir)
	if err != nil {
		return nil, err
	}

	ld := toLineDirection(bdi)

	// running buses
	if withBuses {
		ld.BaseLineDir.Buses, err = getBuses(city, ln, dir)
		if err != nil {
			log.Printf("get %s %s[%s] running buses failed: %v", city, ln, dir, err)
		}
	}

	return ld, nil
}

func toLineDirection(bdi *rtbus.BusDirInfo) *LineDirection {
	// base line
	bl := &BaseLineDir{
		No:      bdi.Name,
		Dir:     fmt.Sprintf("%d", bdi.Direction),
		StartSn: bdi.StartSn,
		EndSn:   bdi.EndSn,
	}

	// another direction id
	if len(bdi.OtherDirIDs) > 0 {
		bl.AnotherDir = bdi.OtherDirIDs[0]
	}

	ld := &LineDirection{
		BaseLineDir: bl,
		ID:          bdi.ID,
		Price:       bdi.Price,
		FirstTime:   bdi.FirstTime,
		LastTime:    bdi.LastTime,
		Stations:    make([]Station, len(bdi.Stations)),
	}
	for i, s := range bdi.Stations {
		ld.Stations[i] = Station{s.No, s.Name, fmt.Sprintf("%f,%f", s.Lon, s.Lat)}
	}

	return ld
}

package apis

import (
	"fmt"
)

func getBuses(city, line, dir string) ([]Bus, error) {
	rbuses, err := RTBusClient.GetRunningBus(city, line, dir)
	if err != nil {
		return []Bus{}, err
	}

	buses := make([]Bus, 0, len(rbuses))
	for _, rbus := range rbuses {
		buses = append(buses, Bus{
			Order:    rbus.No,
			Status:   rbus.Status,
			Location: fmt.Sprintf("%f,%f", rbus.Lng, rbus.Lat),
			Distance: rbus.Distance,
		})
	}

	return buses, nil
}

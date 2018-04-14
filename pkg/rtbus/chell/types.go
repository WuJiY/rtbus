package chell

import (
	"time"

	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

type CllLineDirResp struct {
	Data *CllLineDirData `json:"data"`
}

type CllLineDirData struct {
	Line       *rtbus.BusDirInfo   `json:"line"`
	Bus        []*rtbus.RunningBus `json:"buses"`
	Stations   []*rtbus.BusStation `json:"stations"`
	Otherlines []struct {
		LineId string `json:"lineid"`
	} `json:"otherlines"`
}

type CllLineSearchResp struct {
	ErrMsg   string `json:"errmsg"`
	SVersion string `json:"sversion"`
	Data     struct {
		Lines []CllLineSearchLine `json:"lines"`
	} `json:"data"`
}

type CllLineSearchLine struct {
	EndSn  string `json:"endSn"`
	LineId string `json:"lineId"`
	LineNo string `json:"lineNo"`
	Name   string `json:"name"`
}

func (cdd *CllLineDirData) getBusDirInfo() (bdi *rtbus.BusDirInfo) {
	bdi = cdd.Line
	if bdi == nil {
		return
	}

	bdi.Stations = cdd.Stations
	bdi.RunningBuses = cdd.Bus
	bdi.OtherDirIDs = make([]string, 0)

	curtime := time.Now().Unix()
	for _, rb := range bdi.RunningBuses {
		if rb.Distance == 0 {
			rb.Status = rtbus.BUS_ARRIVING_STATUS
		} else {
			rb.Status = rtbus.BUS_ARRIVING_FUTURE_STATUS
		}

		if rb.No <= 0 {
			LOGGER.Warn("the running bus order is le zero: %d!", rb.No)
		} else if rb.No > len(bdi.Stations) {
			LOGGER.Warn("the running bus number is too large: %d!", rb.No)
		} else {
			rb.Name = bdi.Stations[rb.No-1].Name
		}

		rb.SyncTime = curtime - rb.SyncTime
	}

	return
}

package rtbus

import (
	"fmt"
	"sync"
)

type CityInfo struct {
	Code   string `json:"-"`
	ID     string `json:"cityId"`
	Name   string `json:"cityName"`
	Hot    int    `json:"hot"`
	PinYin string `json:"pinyin"`
	Subway int    `json:"supportSubway"`
}

type CityBusLines struct {
	l          sync.Mutex
	Source     string
	CityInfo   *CityInfo
	ByLineName map[string]*BusLine
}

type BusLine struct {
	l          sync.Mutex
	LineNum    string                 `json:"linenum"`
	LineName   string                 `json:"lineName"`
	Directions map[string]*BusDirInfo `json:"direction"`
}

func (bl *BusLine) GetBusDirInfo(dirname string) (*BusDirInfo, bool) {
	for dirkey, bdi := range bl.Directions {
		//fmt.Printf("%+v\n", bdi)
		if dirname == fmt.Sprintf("%d", bdi.Direction) ||
			dirname == bdi.GetDirName() ||
			dirname == dirkey ||
			dirname == bdi.ID ||
			(bdi.did != "" && dirname == bdi.did) {
			return bdi, true
		}
	}

	return nil, false
}

type BusDirInfo struct {
	l         sync.Mutex
	freshTime int64

	did          string
	ID           string        `json:"id"`
	OtherDirIDs  []string      `json:"otherDirIds"`
	Direction    int           `json:"direction"`
	Name         string        `json:"name"`
	StartSn      string        `json:"startsn,omitempty"`
	EndSn        string        `json:"endsn,omitempty"`
	Price        string        `json:"price,omitempty"`
	SnNum        int           `json:"stationsNum,omitempty"`
	FirstTime    string        `json:"firstTime,omitempty"`
	LastTime     string        `json:"lastTime,omitempty"`
	Stations     []*BusStation `json:"stations"`
	RunningBuses []*RunningBus `json:"buses,omitempty"`
}

func (bdi *BusDirInfo) GetDirName() string {
	return fmt.Sprintf("%s-%s", bdi.StartSn, bdi.EndSn)
}

type BusStation struct {
	No   int     `json:"order"`
	Name string  `json:"sn,omitempty"`
	Lat  float64 `json:"lat,omitempty"`
	Lon  float64 `json:"lon,omitempty"`
}

type RunningBus struct {
	No       int     `json:"order"`
	Name     string  `json:"-"`
	Status   float64 `json:"status"`
	BusID    string  `json:"busid,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Distance int     `json:"distanceToSc"`
	SyncTime int64   `json:"syncTime"`
}

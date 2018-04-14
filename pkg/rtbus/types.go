package rtbus

import (
	"fmt"
	"sort"
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
	Source     string
	CityInfo   *CityInfo
	ByLineName map[string]*BusLine
}

type BusLine struct {
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
func (bl *BusLine) PutDir(bd *BusDirInfo) {
	bl.Directions[bd.GetDirName()] = bd
}

type BusDirInfo struct {
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

type RunBuses []*RunningBus

func (rb RunBuses) Len() int {
	return len(rb)
}
func (rb RunBuses) Swap(i, j int) {
	rb[i], rb[j] = rb[j], rb[i]
}
func (rb RunBuses) Less(i, j int) bool {
	return rb[i].No < rb[j].No || (rb[i].No == rb[j].No && rb[i].Distance > rb[j].Distance)
}

func (bdi *BusDirInfo) Sort() {
	rbs := RunBuses(bdi.RunningBuses)
	sort.Sort(rbs)
	bdi.RunningBuses = []*RunningBus(rbs)
}

type BusStation struct {
	No   int     `json:"order"`
	Name string  `json:"sn,omitempty"`
	Lat  float64 `json:"lat,omitempty"`
	Lon  float64 `json:"lon,omitempty"`
}

type RunningBus struct {
	No       int     `json:"order"`
	Name     string  `json:"name,omitempty"`
	Status   float64 `json:"status"`
	BusID    string  `json:"busid,omitempty"`
	Lat      float64 `json:"lat,omitempty"`
	Lng      float64 `json:"lng,omitempty"`
	Distance int     `json:"distanceToSc"`
	SyncTime int64   `json:"syncTime"`
}

package beijing

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/xuebing1110/rtbus/pkg/rtbus"
	"strconv"
)

const (
	BTIC_RAW_KEY_FMT = `aibang%s`
	//BTIC_KEY_SECRET     = `bjjw_jtcx`
	BTIC_LISTLINES_PATH = `/ssgj/v1.0.0/checkupdate`
	BTIC_BUSLINE_PATH   = `/ssgj/v1.0.0/update`
)

func LoadBusLines(version string) (<-chan *rtbus.BusDirInfo, error) {
	linedir_chan := make(chan *rtbus.BusDirInfo, 0)

	allline := &BticAllLineResp{}
	var params *url.Values = &url.Values{}
	params.Set("version", version)
	err := bticRequest(BTIC_LISTLINES_PATH, params, allline)
	if err != nil {
		return linedir_chan, err
	}
	if allline.ErrCode != "200" {
		err = errors.New(allline.ErrMsg)
		return linedir_chan, err
	}

	go func(lines []*BticBasicLine) {
		var rateLimit chan bool = make(chan bool, 20)
		var wg sync.WaitGroup
		for _, line := range lines {
			rateLimit <- true

			wg.Add(1)
			go func(id string) {
				defer func() {
					<-rateLimit
					wg.Done()
				}()

				bd, err := loadLineDir(id)
				if err != nil {
					LOGGER.Errorf("init BJ lineid %s failed:%v", id, err)
					return
				}
				linedir_chan <- bd
			}(line.ID)
		}

		//wait complete
		wg.Wait()
		close(rateLimit)
		close(linedir_chan)
	}(allline.Lines.Line)

	return linedir_chan, nil
}

func loadLineDir(id string) (*rtbus.BusDirInfo, error) {
	abline, err := getBticLineDir(id)
	if err != nil {
		LOGGER.Errorf("init BJ lineid %s failed:%v", id, err)
		return nil, err
	}

	return parseBticLineToBuslineDir(abline)
}

func parseBticLineToBuslineDir(line *BticLine) (bdi *rtbus.BusDirInfo, err error) {
	sNum := len(line.Stations.Station)
	if sNum == 0 {
		err = fmt.Errorf("can't find the any station of line %s from aibang", line.LineName)
		return
	}

	firstS := line.Stations.Station[0]
	lastS := line.Stations.Station[sNum-1]

	var price string
	if line.TotalPrice == 0 {
		price = line.Ticket
	} else {
		price = fmt.Sprintf("%.0f", line.TotalPrice)
	}

	var firsrt_time, last_time string
	start_end_time := strings.SplitN(line.Time, "-", 2)
	firsrt_time = start_end_time[0]
	if len(start_end_time) > 1 {
		last_time = start_end_time[1]
	}

	bdi = &rtbus.BusDirInfo{
		ID:          line.ID,
		Direction:   1,
		OtherDirIDs: []string{"0"},
		StartSn:     firstS.Name,
		EndSn:       lastS.Name,
		Price:       price,
		SnNum:       sNum,
		FirstTime:   firsrt_time,
		LastTime:    last_time,
		Stations:    convertABLineStation(line.Stations.Station),
	}

	bdi.Name = line.ShortName
	return
}

func newBusLineByABLine(bdi *rtbus.BusDirInfo) (bl *rtbus.BusLine) {
	bdi.Direction = 0
	bdi.OtherDirIDs = []string{"1"}

	return &rtbus.BusLine{
		LineNum:  bdi.Name,
		LineName: bdi.Name,
		Directions: map[string]*rtbus.BusDirInfo{
			bdi.GetDirName(): bdi,
		},
	}
}

func convertABLineStation(abss []*BticLineStation) []*rtbus.BusStation {
	var err error

	bss := make([]*rtbus.BusStation, len(abss))
	for i, abs := range abss {
		bs := &rtbus.BusStation{
			Name: abs.Name,
		}

		//StationNo
		bs.No, err = strconv.Atoi(abs.No)
		if err != nil {
			bs.No = i + 1
		}

		//lat lon
		bs.Lat, _ = strconv.ParseFloat(abs.Lat, 10)
		bs.Lon, _ = strconv.ParseFloat(abs.Lon, 10)

		bss[i] = bs
	}

	return bss
}

func getBticLineDir(id string) (line *BticLine, err error) {
	params := &url.Values{}
	params.Set("id", id)

	lineResp := &BticLineResp{}
	err = bticRequest(BTIC_BUSLINE_PATH, params, lineResp)
	if err != nil {
		return
	}
	if lineResp.BusLine == nil || len(lineResp.BusLine) == 0 {
		err = errors.New("busline info is null")
		return
	}

	//decrypt
	line = lineResp.BusLine[0]
	key := fmt.Sprintf(BTIC_RAW_KEY_FMT, id)
	line.ShortName, _ = Rc4DecodeString(key, line.ShortName)
	line.LineName, _ = Rc4DecodeString(key, line.LineName)
	line.Coord, _ = Rc4DecodeString(key, line.Coord)

	for _, station := range line.Stations.Station {
		station.Name, _ = Rc4DecodeString(key, station.Name)
		station.No, _ = Rc4DecodeString(key, station.No)
		station.Lon, _ = Rc4DecodeString(key, station.Lon)
		station.Lat, _ = Rc4DecodeString(key, station.Lat)
	}

	//fmt.Printf("%+v\n", line.Stations)

	return
}

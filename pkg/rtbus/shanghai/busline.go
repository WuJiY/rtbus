package shanghai

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	path_listlines      = "/interface/GetYBLineList.ashx?sign="
	path_fmt_getbase    = "/interface/getBase.ashx?sign=&name=%s"
	path_fmt_liststop   = "/interface/getStopList.ashx?name=%s&lineid=%s&dir=%d"
	path_fmt_runningbus = "/interface/getCarmonitor.ashx?name=%s&lineid=%s&stopid=%d&dir=%d"
)

var (
	slines     sort.StringSlice
	REG_NUMBER = regexp.MustCompile(`^\d+$`)
)

func init() {
	lines, err := ListLines()
	if err != nil {
		panic(err)
	}

	slines = sort.StringSlice(lines)
	slines.Sort()
}

func ListLines() ([]string, error) {
	list := make([]string, 0)
	err := doRequest(path_listlines, &list)
	return list, err
}

func linePrefixSearch(x string) []string {
	return prefixSearch(slines, x)
}

func lineSearch(x string) string {
	index := slines.Search(x)
	if index < 0 {
		return ""
	}

	return slines[index]
}

func prefixSearch(p []string, x string) []string {
	input_len := len(x)
	index := sort.Search(len(p), func(i int) bool {
		return bytes.Compare([]byte(p[i])[:input_len], []byte(x)) >= 0
	})

	le_strings := make([]string, 0)
	for i := index; i >= 0; i-- {
		if strings.HasPrefix(p[i], x) {
			le_strings = append(le_strings, p[i])
		}
	}

	gt_strings := make([]string, 0)
	for i := index + 1; i < len(p); i++ {
		if strings.HasPrefix(p[i], x) {
			gt_strings = append(gt_strings, p[i])
		}
	}

	ret := make([]string, 0, len(le_strings)+len(gt_strings))
	for i := len(le_strings) - 1; i >= 0; i-- {
		ret = append(ret, le_strings[i])
	}
	ret = append(ret, gt_strings...)

	return ret
}

type BaseLine struct {
	XMLName xml.Name `xml:"linedetail"`

	LineId         string `xml:"line_id"`
	LineName       string `xml:"line_name"`
	StartSN        string `xml:"start_stop"`
	EndSN          string `xml:"end_stop"`
	StartEarlytime string `xml:"start_earlytime"`
	StartLatetime  string `xml:"start_latetime"`
	EndEarlytime   string `xml:"end_earlytime"`
	EndLatetime    string `xml:"end_latetime"`
}

func _getlineNo(l string) string {
	if REG_NUMBER.Match([]byte(l)) {
		l = l + "路"
	}

	l = strings.Replace(strconv.QuoteToASCII(l), `\`, "%", -1)
	l, _ = strconv.Unquote(l)

	return l
}

func getBusLine(lineno string) (*rtbus.BusLine, error) {
	lineno = _getlineNo(lineno)

	baseline, err := _getBaseLine(lineno)
	if err != nil {
		return nil, err
	}

	bl := &rtbus.BusLine{
		LineNum:    baseline.LineId,
		LineName:   baseline.LineName,
		Directions: make(map[string]*rtbus.BusDirInfo),
	}
	dirname := "0"
	bdi, err := getBusLineDir(lineno, dirname)
	if err != nil {
		return nil, err
	}
	bl.PutDir(bdi)

	for _, dirname = range bdi.OtherDirIDs {
		bdi, err = getBusLineDir(lineno, dirname)
		if err != nil {
			continue
		}

		bl.PutDir(bdi)
	}

	return bl, nil
}

func _getBaseLine(lineno string) (*BaseLine, error) {
	// base info && get line id
	path := fmt.Sprintf(path_fmt_getbase, lineno)
	bl := new(BaseLine)
	if err := doRequest(path, bl); err != nil {
		return nil, fmt.Errorf("get %s error: %v", path, err)
	}

	bl.StartSN = strings.TrimSpace(bl.StartSN)
	bl.EndSN = strings.TrimSpace(bl.EndSN)

	bl.LineId = strings.TrimSpace(bl.LineId)
	bl.LineName = strings.TrimSpace(bl.LineName)
	bl.StartSN = strings.TrimSpace(bl.StartSN)
	bl.EndSN = strings.TrimSpace(bl.EndSN)
	bl.StartEarlytime = strings.TrimSpace(bl.StartEarlytime)
	bl.StartLatetime = strings.TrimSpace(bl.StartLatetime)
	bl.EndEarlytime = strings.TrimSpace(bl.EndEarlytime)
	bl.EndLatetime = strings.TrimSpace(bl.EndLatetime)

	return bl, nil
}

func getBusLineDir(lineno, dirname string) (*rtbus.BusDirInfo, error) {
	lineno = _getlineNo(lineno)

	baseline, err := _getBaseLine(lineno)
	if err != nil {
		return nil, err
	}

	dir := 0
	switch dirname {
	case "0":
		dir = 0
	case "1":
		dir = 1
	default:
		if dirname != "" && dirname != baseline.StartSN+"-"+baseline.EndSN {
			dir = 1
		} else {
			dir = 0
		}
	}

	var bdi *rtbus.BusDirInfo
	if dir == 0 {
		bdi = &rtbus.BusDirInfo{
			ID:        baseline.LineId,
			Name:      baseline.LineName,
			Direction: dir,
			StartSn:   baseline.StartSN,
			EndSn:     baseline.EndSN,
			Price:     "--",
			FirstTime: baseline.StartEarlytime,
			LastTime:  baseline.StartLatetime,
		}

	} else {
		bdi = &rtbus.BusDirInfo{
			ID:        baseline.LineId,
			Name:      baseline.LineName,
			Direction: dir,
			StartSn:   baseline.EndSN,
			EndSn:     baseline.StartSN,
			Price:     "--",
			FirstTime: baseline.EndEarlytime,
			LastTime:  baseline.EndLatetime,
		}
	}

	if bdi.StartSn == bdi.EndSn {
		bdi.OtherDirIDs = []string{}
	} else {
		if dir == 0 {
			bdi.OtherDirIDs = []string{"1"}
		} else {
			bdi.OtherDirIDs = []string{"0"}
		}
	}

	// get station
	stops := make([]Stop, 0)
	path := fmt.Sprintf(path_fmt_liststop, lineno, baseline.LineId, dir)
	resp := &Response{Data: &stops}
	err = doRequest(path, &resp)
	if err != nil {
		return bdi, fmt.Errorf("get %s error: %v", path, err)
	}

	bdi.SnNum = len(stops)
	bdi.Stations = make([]*rtbus.BusStation, len(stops))
	for i, s := range stops {
		//fmt.Printf("%+v\n", s)
		bdi.Stations[i] = &rtbus.BusStation{No: s.Id, Name: s.Name}
	}

	return bdi, nil
}

type Stop struct {
	Name string `json:"name"`
	Id   int    `json:"id,string"`
}

func getRunningBus(bdi *rtbus.BusDirInfo, stopid int) (rb []*rtbus.RunningBus, err error) {
	lineno := _getlineNo(bdi.Name)
	if stopid <= 0 {
		stopid = bdi.Stations[bdi.SnNum-1].No
	}

	path := fmt.Sprintf(path_fmt_runningbus, lineno, bdi.ID, stopid, bdi.Direction)
	data := new(RunningBus)
	resp := &Response{Data: data}
	err = doRequest(path, resp)
	if err != nil {
		return
	}

	if data.Terminal == "" || data.Terminal == "null" {
		return
	}

	rb = make([]*rtbus.RunningBus, 0)

	//fmt.Printf("prepare to get %d-%d+1 running bus...\n", stopid, data.StopDis)
	stopid_bus := stopid - data.StopDis + 1

	// distance
	var distance = -1
	if stopid > 1 {
		path = fmt.Sprintf(path_fmt_runningbus, lineno, bdi.ID, stopid_bus, bdi.Direction)
		data_2 := new(RunningBus)
		resp_2 := &Response{Data: data_2}
		err = doRequest(path, resp_2)
		if err == nil {
			distance = data_2.Distance
		}
	}

	// status
	var status = rtbus.BUS_ARRIVING_FUTURE_STATUS
	if distance < 50 {
		status = rtbus.BUS_ARRIVING_STATUS
	}

	rb = append(rb, &rtbus.RunningBus{
		No:       stopid_bus,
		Name:     bdi.Stations[stopid_bus-1].Name,
		Status:   status,
		BusID:    data.Terminal,
		Lat:      0,
		Lng:      0,
		Distance: distance,
		Time:     data.Time,
		SyncTime: time.Now().Unix(),
	})

	if stopid_bus-1 > 0 {
		//fmt.Printf("get %d running bus...\n", stopid_bus-1)
		rb_2, err := getRunningBus(bdi, stopid_bus-1)
		if err != nil {
			return rb, err
		}

		for _, rb_tmp := range rb_2 {
			if rb_tmp.BusID == rb[len(rb)-1].BusID {
				continue
			}
			rb = append(rb, rb_tmp)
		}
	}

	return rb, nil
}

type RunningBus struct {
	Terminal string `json:"terminal"`
	StopDis  int    `json:"stopdis,string"`
	Time     int    `json:"time,string"`
	Distance int    `json:"distance,string"`
}

package chell

import (
	"fmt"
	"time"

	"github.com/bingbaba/util/httptool"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

const (
	FMT_CLL_URL_PARAMS = "lineId=%s&lineName=%s&lineNo=%s&s=h5&v=3.3.9&userId=browser_%d&h5Id=browser_%d&sign=1&cityId=%s"
	URL_CLL_BUS_URL    = "http://web.chelaile.net.cn/api/bus/line!lineDetail.action"
	FMT_CLL_URL_SEARCH = `http://web.chelaile.net.cn/api/basesearch/client/clientSearch.action?key=%s&count=3&s=h5&v=3.3.9&userId=browser_%d&h5Id=browser_%d&sign=1&cityId=%s`
)

func search(cityid, keyword string) (csls []CllLineSearchLine, err error) {
	curtime := time.Now().UnixNano() / 1000000
	reqUrl := fmt.Sprintf(
		FMT_CLL_URL_SEARCH,
		keyword,
		curtime, curtime,
		cityid,
	)

	httpreq, err := getCllHttpRequest(reqUrl)
	if err != nil {
		return csls, err
	}

	cllresp := &CllLineSearchResp{}
	err = httptool.HttpDoJsonr(httpreq, cllresp)
	if err != nil {
		return csls, err
	}

	if cllresp.ErrMsg != "" || len(cllresp.Data.Lines) == 0 {
		err = fmt.Errorf("search %s line failed:%s", keyword, cllresp.ErrMsg)
		return csls, err
	}

	lines := cllresp.Data.Lines
	csls = make([]CllLineSearchLine, 0, len(lines)*2)
	for _, line := range lines {
		csls = append(csls, line)
	}

	return csls, nil
}

func loadBusline(cityid, lineno string) (*rtbus.BusLine, error) {
	csls, err := search(cityid, lineno)
	if err != nil {
		return nil, err
	}

	var bdi *rtbus.BusDirInfo
	lineid := csls[0].LineId
	bdi, err = getNewestCllBusDirInfo(cityid, lineid, lineno)
	if err != nil {
		return nil, err
	}

	//BusLine
	bl := &rtbus.BusLine{
		LineNum:  lineno,
		LineName: lineno,
		Directions: map[string]*rtbus.BusDirInfo{
			bdi.GetDirName(): bdi,
		},
	}

	//other line
	//fmt.Printf("%+v\n", cdd.Otherlines)
	for _, olineid := range bdi.OtherDirIDs {
		obdi, err := getNewestCllBusDirInfo(cityid, olineid, lineno)
		if err != nil {
			return nil, err
		}
		bl.Directions[obdi.GetDirName()] = obdi
	}

	return bl, nil
}

func getNewestCllBusDirInfo(cityid, lineid, lineno string) (*rtbus.BusDirInfo, error) {
	curtime := time.Now().UnixNano() / 1000000
	reqUrl := URL_CLL_BUS_URL +
		"?" +
		fmt.Sprintf(
			FMT_CLL_URL_PARAMS,
			lineid, lineno, lineno,
			curtime, curtime,
			cityid,
		)
	//fmt.Println(reqUrl) //debug

	httpreq, err := getCllHttpRequest(reqUrl)
	if err != nil {
		return nil, err
	}

	var cllresp *CllLineDirResp = &CllLineDirResp{}
	err = httptool.HttpDoJsonr(httpreq, cllresp)
	if err != nil {
		return nil, err
	}

	//fmt.Println(ToJsonString(cllresp))
	cdd := cllresp.Data
	bdi := cdd.getBusDirInfo()
	if bdi == nil {
		err = fmt.Errorf("can't get %s-%s-%s bus info!", cityid, lineid, lineno)
		return nil, err
	}
	bdi.ID = lineid

	for _, oline := range cdd.Otherlines {
		bdi.OtherDirIDs = append(bdi.OtherDirIDs, oline.LineId)
	}

	return bdi, nil
}

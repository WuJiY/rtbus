package chell

import (
	"errors"
	"fmt"
	"time"

	"github.com/bingbaba/util/httptool"
	"github.com/xuebing1110/location"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

const (
	URL_CLL_CITYS_FMT = `http://web.chelaile.net.cn/cdatasource/citylist?type=allRealtimeCity&s=h5&v=3.3.9&userId=browser_%d`
)

type AllCityResp struct {
	Status string `json:"status"`
	Data   struct {
		AllRealtimeCity []*rtbus.CityInfo `json:"allRealtimeCity"`
	} `json:"data"`
}

func GetAllCitys() ([]*rtbus.CityInfo, error) {
	reqUrl := fmt.Sprintf(URL_CLL_CITYS_FMT, time.Now().UnixNano()/1000000)
	httreq, err := getCllHttpRequest(reqUrl)
	if err != nil {
		return nil, err
	}

	cllresp := &AllCityResp{}
	err = httptool.HttpDoJson(httreq, cllresp)
	if err != nil {
		return nil, err
	}

	if cllresp.Status != "OK" {
		return nil, errors.New(cllresp.Status)
	}

	for _, city := range cllresp.Data.AllRealtimeCity {
		cityName := city.Name
		city.Code = location.GetCitycode(location.MustParseCity(cityName))
	}

	return cllresp.Data.AllRealtimeCity, nil
}

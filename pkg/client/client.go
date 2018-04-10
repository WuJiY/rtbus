package client

import (
	"github.com/xuebing1110/location"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

type RTBus struct {
	apis map[string]rtbus.CityRTBusApi
}

func newRTBus() *RTBus {
	return &RTBus{
		apis: make(map[string]rtbus.CityRTBusApi),
	}
}

func (rtb *RTBus) Register(cba rtbus.CityRTBusApi) bool {
	city := cba.City().Code
	_, found := rtb.apis[city]
	if found {
		return false
	}

	rtb.apis[city] = cba
	return true
}

func (rtb *RTBus) MustRegister(cba rtbus.CityRTBusApi) {
	city := cba.City().Code
	delete(rtb.apis, city)
	rtb.apis[city] = cba
}

func (rtb *RTBus) GetBusLine(city, lineno string) (bl *rtbus.BusLine, err error) {
	cba, found := rtb.getCityRTBus(city)
	if !found {
		return nil, ERROR_NOTSUPPORT
	}

	return cba.GetBusLine(lineno)
}

func (rtb *RTBus) GetBusLineDir(city, lineno, dirname string) (bdi *rtbus.BusDirInfo, err error) {
	cba, found := rtb.getCityRTBus(city)
	if !found {
		return nil, ERROR_NOTSUPPORT
	}

	return cba.GetBusLineDir(lineno, dirname)
}

func (rtb *RTBus) GetRunningBus(city, lineno, dirname string) (rbus []*rtbus.RunningBus, err error) {
	cba, found := rtb.getCityRTBus(city)
	if !found {
		return nil, ERROR_NOTSUPPORT
	}

	return cba.GetRunningBus(lineno, dirname)
}

func (rtb *RTBus) getCityRTBus(city string) (cba rtbus.CityRTBusApi, found bool) {
	city_code := location.GetCitycode(location.MustParseCity(city))
	if city_code == "" {
		city_code = city
	}
	cba, found = rtb.apis[city_code]
	return
}

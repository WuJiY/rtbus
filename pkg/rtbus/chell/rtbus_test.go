package chell

import (
	"testing"
	//"time"

	//"github.com/xuebing1110/rtbus/pkg/rtbus"
	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

func TestRTBusApi(t *testing.T) {
	// qingdao
	citys, err := GetAllCitys()
	if err != nil {
		t.Fatal(err)
	}
	var qingdao_city *rtbus.CityInfo
	for _, city := range citys {
		if city.Code == "0532" {
			qingdao_city = city
			break
		}
	}

	var api rtbus.CityRTBusApi
	api = NewChellRTBusApi(qingdao_city)

	lineno, linedir := testLine()
	bl, err := api.GetBusLine(lineno)
	if err != nil {
		t.Fatal(err)
	}
	if len(bl.Directions) != 2 {
		t.Fatalf("should be two direction for line %s", lineno)
	}

	_, err = api.GetBusLineDir(lineno, linedir)
	if err != nil {
		t.Fatal(err)
	}

	rbuses, err := api.GetRunningBus(lineno, linedir)
	if err != nil {
		t.Fatal(err)
	}
	if len(rbuses) == 0 {
		t.Fatalf("can't get any running buses in line %s", lineno)
	}
}

func testLine() (string, string) {
	return "318", "市政府-虎山军体中心"
}

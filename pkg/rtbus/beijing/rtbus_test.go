package beijing

import (
	"testing"
	//"time"

	"github.com/xuebing1110/rtbus/pkg/rtbus"
)

func TestRTBusApi(t *testing.T) {
	var api rtbus.CityRTBusApi
	var err error
	api, err = NewBJRTBusApi()
	if err != nil {
		t.Fatal(err)
	}

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
	return "675", "通州李庄-左家庄"
}

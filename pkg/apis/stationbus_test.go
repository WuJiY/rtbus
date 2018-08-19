package apis

import (
	"fmt"
	"testing"
)

func TestListLocalStationBuses(t *testing.T) {
	sbs, err := ListLocalStationBuses("120.44052,36.18169", 1, 3)
	if err != nil {
		t.Fatal(err)
	}

	if len(sbs) == 0 {
		t.Fatalf("can't get any station bus")
	}

	//for _, sb := range sbs {
	//	fmt.Printf("%+v\n", sb)
	//}

	if sbs[0].StationName != "宜川路合水路" &&
		sbs[0].StationName != "广水路金川路" &&
		sbs[0].StationName != "绿城百合花园" {
		t.Fatalf("search around poi failed!")
	}

	if sbs[0].SupportLineCount == 0 {
		t.Fatalf("not support %s rtbus", sbs[0].StationName)
	}
}

func TestGetStationBusesByStation(t *testing.T) {
	sb, err := GetStationBusesByStation("0532", "宜川路合水路", 3)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", sb)
	if sb.StationName != "宜川路合水路" {
		t.Fatalf("search around poi failed!")
	}

	if sb.SupportLineCount == 0 {
		t.Fatalf("not support %s rtbus", sb.StationName)
	}
}

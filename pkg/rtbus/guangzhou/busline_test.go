package guangzhou

import (
	"fmt"
	"testing"
)

func TestGetBusLineDirByLineid(t *testing.T) {
	bdi, err := getBusLineDirByLineid("8", 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range bdi.Stations {
		fmt.Printf("%+v\n", s)
	}

	for _, s := range bdi.RunningBuses {
		fmt.Printf("%+v\n", s)
	}
}

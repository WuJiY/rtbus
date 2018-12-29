package shanghai

import (
	//"fmt"
	"fmt"
	"testing"
)

func TestGetRunningBus(t *testing.T) {
	bdi, err := getBusLineDir("933", "")
	if err != nil {
		t.Fatal(err)
	}

	rbs, err := getRunningBus(bdi, 25)
	if err != nil {
		t.Fatal(err)
	}

	if len(rbs) == 0 {
		t.Fatalf("get running bus fault")
	}

	for _, rb := range rbs {
		fmt.Printf("%+v\n", rb)
	}
}

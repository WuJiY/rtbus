package shanghai

import (
	"fmt"
	"testing"
)

func TestSearch(t *testing.T) {
	api, err := NewRTBusApi()
	if err != nil {
		t.Fatal(err)
	}

	bdis, err := api.Search("933")
	if err != nil {
		t.Fatal(err)
	}

	if len(bdis) == 0 {
		t.Fatalf("get 933 line fault")
	}

	//for _, bdi := range bdis {
	//	fmt.Printf("%+v\n", bdi)
	//}
}

func TestGetLine(t *testing.T) {
	api, err := NewRTBusApi()
	if err != nil {
		t.Fatal(err)
	}

	bl, err := api.GetBusLine("933", true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", bl)
	for key, bdi := range bl.Directions {
		fmt.Printf("%s : %+v\n", key, bdi)
	}
}

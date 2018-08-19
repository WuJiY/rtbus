package apis

import (
	"fmt"
	"testing"
)

func TestGetBusLineDir(t *testing.T) {
	ld, err := GetBusLineDir("0532", "643", "0", true)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%+v\n", ld)
}

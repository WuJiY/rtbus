package client

import "fmt"

var (
	ERROR_NOTFOUND   = fmt.Errorf("not found")
	ERROR_NOTSUPPORT = fmt.Errorf("not support")
)

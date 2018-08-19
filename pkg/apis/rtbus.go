package apis

import (
	"github.com/xuebing1110/rtbus/pkg/client"
)

var (
	RTBusClient *client.RTBus
)

func init() {
	RTBusClient = client.DefaultRTBus
}

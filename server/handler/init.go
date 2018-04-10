package handler

import (
	"github.com/bingbaba/util/logs"
	"github.com/xuebing1110/rtbus/pkg/client"
)

var (
	RTBusClient *client.RTBus
	logger      *logs.Blogger
)

func init() {
	logger = logs.GetBlogger()
	RTBusClient = client.DefaultRTBus
}

type Response struct {
	ErrNo  int         `json:"errno"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data,omitempty"`
}

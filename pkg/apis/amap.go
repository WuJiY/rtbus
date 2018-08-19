package apis

import (
	"os"

	"github.com/xuebing1110/location/amap"
	"github.com/xuebing1110/rtbus/pkg/httputil"
)

var (
	amapClient *amap.Client
	AMAP_KEY   = `b3abf03fa1e83992727f0625a918fe73`
)

func init() {
	key := os.Getenv("AMAP_KEY")
	if key != "" {
		AMAP_KEY = key
	}
	amapClient = amap.NewClient(AMAP_KEY)
	amapClient.HttpClient = httputil.DEFAULT_HTTP_CLIENT
}

func SetAmapKey(key string) {
	amapClient = amap.NewClient(AMAP_KEY)
	amapClient.HttpClient = httputil.DEFAULT_HTTP_CLIENT
}

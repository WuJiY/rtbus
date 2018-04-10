package chell

import (
	"net/http"
)

const (
	URL_CLL_REFER = "http://web.chelaile.net.cn/ch5/index.html"
)

func getCllHttpRequest(req_url string) (httpreq *http.Request, err error) {
	httpreq, err = http.NewRequest("GET", req_url, nil)
	if err != nil {
		return
	}

	httpreq.Header.Add("Accept", "application/json, text/plain, */*")
	httpreq.Header.Add("Referer", URL_CLL_REFER)
	httpreq.Header.Add("User-Agent", `Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/57.0.2987.98 Mobile Safari/537.36`)
	return
}

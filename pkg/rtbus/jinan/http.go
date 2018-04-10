package jinan

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/xuebing1110/rtbus/pkg/httputil"
)

const (
	HOST           = "60.216.101.229"
	HEADER_VERSION = "android-insigma.waybook.jinan-2342"
)

func doRequest(path string, v interface{}) (err error) {
	req_url := "http://" + HOST + path

	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Host", HOST)
	req.Header.Set("version", HEADER_VERSION)

	var resp *http.Response
	resp, err = httputil.HttpDo(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	//read all
	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	return json.Unmarshal(data, v)
}

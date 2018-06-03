package guangzhou

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/xuebing1110/rtbus/pkg/httputil"
)

const (
	HOST           = "rycxapi.gci-china.com"
	HEADER_VERSION = "android-insigma.waybook.jinan-2342"
	USER_AGENT     = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15F79 MicroMessenger/6.6.5 NetType/WIFI Language/zh_CN"
	REFERER        = "https://servicewechat.com/wxe027f0adf505a625/12/page-frame.html"
)

func doRequest(path string, v interface{}) (err error) {
	req_url := "https://" + HOST + path

	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Host", HOST)
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Referer", REFERER)

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

func checkResponse(resp *Response) error {
	if resp.ErrCode != 0 || (resp.Message != "Success" && resp.Message != "OK") {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

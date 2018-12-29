package shanghai

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"encoding/xml"
	"github.com/xuebing1110/rtbus/pkg/httputil"
	"golang.org/x/text/encoding/simplifiedchinese"
)

const (
	HOST       = "61.129.57.81:8181"
	USER_AGENT = "Mozilla/5.0 (iPhone; CPU iPhone OS 11_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15F79 MicroMessenger/6.6.5 NetType/WIFI Language/zh_CN"
	REFERER    = "http://61.129.57.81:8181/BusEstimate.aspx"
)

func doRequest(path string, v interface{}) (err error) {
	req_url := "http://" + HOST + path

	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Host", HOST)
	req.Header.Set("User-Agent", USER_AGENT)
	req.Header.Set("Referer", REFERER)

	var resp *http.Response
	resp, err = httputil.HttpDo(req)
	if err != nil {
		return
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var body_utf8 []byte
	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	//if strings.Index(ct, "charset=gb2312") >= 0 {
	//	body_utf8, err = simplifiedchinese.HZGB2312.NewDecoder().Bytes(body)
	//	if err != nil {
	//		return
	//	}
	//} else
	if strings.Index(ct, "charset=gb") >= 0 {
		body_utf8, err = simplifiedchinese.GBK.NewDecoder().Bytes(body)
		if err != nil {
			return
		}
	} else {
		body_utf8 = body
	}

	//fmt.Printf("%s: %s\n", path, body_utf8)
	if strings.Index(ct, "text/html") >= 0 ||
		strings.Index(ct, "text/xml") >= 0 {

		if body_utf8[0] == '<' {
			return xml.Unmarshal(body_utf8, v)
		} else {
			return json.Unmarshal(body_utf8, v)
		}

	} else { // json
		return json.Unmarshal(body_utf8, v)
	}

}

func checkResponse(resp *Response) error {
	if resp.ErrCode != 1 || (resp.Message != "成功") {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

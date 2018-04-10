package beijing

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/xuebing1110/rtbus/pkg/httputil"
)

var (
	btic_headers = map[string]string{
		"HEADER_KEY_SECRET": "bjjw_jtcx",
		"IMSI":              "360120018215321",
		"UA":                "MATE8",
		"PLATFORM":          "android",
		"CUSTOM":            "aibang",
		"CID":               "67a88ec31de7a589a2344cc5d0469074",
		"IMEI":              "89031020265872",
		"NETWORK":           "gprs",
		"PKG_SOURCE":        "1",
		"CTYPE":             "json",
		"VID":               "5",
		"SOURCE":            "1",
		"UID":               "",
		"SID":               "",
		"PID":               "5",
		"Host":              "transapp.btic.org.cn:8512",
		"User-Agent":        "okhttp/3.3.1",
	}
)

func bticRequest(path string, params *url.Values, v interface{}) (err error) {
	req_url := "http://" + btic_headers["Host"] + path + "?" + params.Encode()
	// LOGGER.Info("url: %s", req_url)
	req, err := http.NewRequest(http.MethodGet, req_url, nil)
	if err != nil {
		return
	}

	// header
	for key, value := range btic_headers {
		req.Header.Set(key, value)
	}
	cur_time := fmt.Sprintf("%d", time.Now().Unix())
	req.Header.Set("TIME", cur_time)

	// token
	token := generateToken(cur_time, path)
	req.Header.Set("ABTOKEN", token)

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

func generateToken(cur_time, path string) string {
	body := fmt.Sprintf("%s%s%s", btic_headers["HEADER_KEY_SECRET"]+btic_headers["PLATFORM"]+btic_headers["CID"], cur_time, path)

	// LOGGER.Info("content: %s", body)
	sha1_data := fmt.Sprintf("%x", sha1.Sum([]byte(body)))
	token := fmt.Sprintf("%x", md5.Sum([]byte(sha1_data)))
	// LOGGER.Info("token: %s", token)
	return token
}

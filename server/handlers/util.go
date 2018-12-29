package handlers

import (
	gocxt "context"
	"encoding/json"
	"time"

	//"github.com/kataras/iris/context"
	"github.com/bingbaba/util/logs"
	"github.com/gin-gonic/gin"
)

var (
	REQUEST_TIMEOUT = 5 * time.Second
	logger          = new(logs.StdoutLogger)
)

type Response struct {
	Message  string      `json:"code"`
	Detail   string      `json:"detail"`
	Data     interface{} `json:"data,omitempty"`
	ScrollId string      `json:"scrollId,omitempty"`
}

func NewResponse(data []byte) (*Response, error) {
	r := new(Response)
	err := json.Unmarshal(data, r)
	return r, err
}

func sendResponse(ctx *gin.Context, data interface{}) {
	resp := &Response{
		Message: "OK",
		Detail:  "",
		Data:    data,
	}

	sendJson(ctx, resp)
}

func sendResponseWithScrollId(ctx *gin.Context, data interface{}, scrollid string) {
	resp := &Response{
		Message:  "OK",
		Detail:   "",
		Data:     data,
		ScrollId: scrollid,
	}

	sendJson(ctx, resp)
}

func sendJson(ctx *gin.Context, v interface{}) {
	pretty := ctx.Param("pretty")
	if pretty == "" || pretty == "0" || pretty == "false" || pretty == "False" {
		ctx.JSON(200, v)
	} else {
		ctx.JSON(200, v)
		//ctx.JSON(v, context.JSON{Indent: "    "})
	}
}

func sendJsonWithCode(ctx *gin.Context, code int, v interface{}) {
	pretty := ctx.Query("pretty")
	if pretty == "" || pretty == "0" || pretty == "false" || pretty == "False" {
		ctx.JSON(code, v)
	} else {
		ctx.IndentedJSON(code, v)
	}
}

func sendBadResponse(ctx *gin.Context, code int, msg, detail string) {
	sendBadResponseWithData(ctx, code, msg, detail, nil)
}

func sendBadResponseWithData(ctx *gin.Context, code int, msg, detail string, data interface{}) {
	resp := &Response{
		Message: msg,
		Detail:  detail,
		Data:    data,
	}
	ctx.Set("message", msg)
	ctx.Set("detail", detail)
	sendJsonWithCode(ctx, code, resp)
}

func getTimeoutContext(ctx *gin.Context) (gocxt.Context, gocxt.CancelFunc) {
	return gocxt.WithTimeout(gocxt.Background(), REQUEST_TIMEOUT)
}

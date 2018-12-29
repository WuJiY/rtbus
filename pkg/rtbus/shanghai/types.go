package shanghai

type Response struct {
	ErrCode int         `json:"code,string"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

package beijing

type BticAllLineResp struct {
	ErrCode     string `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	UpdateNum   string `json:"updateNum"`
	LineNum     string `json:"lineNum"`
	DataVersion string `json:"dataversion"`

	Lines struct {
		Line []*BticBasicLine `json:"line"`
	} `json:"lines"`
}

type BticBasicLine struct {
	ID       string `json:"id"`
	LineName string `json:"linename"`
	Classify string `json:"classify"`
	Status   string `json:"status"`
	Version  string `json:"version"`
}

type BticLineResp struct {
	ErrCode string `json:"errcode"`
	ErrMsg  string `json:"errmsg"`

	BusLine []*BticLine `json:"busline"`
}

type BticLine struct {
	ID         string  `json:"lineid"`
	ShortName  string  `json:"shotname"`
	LineName   string  `json:"linename"`
	Distince   float64 `json:"distince,string"`
	Ticket     string  `json:"ticket"`
	TotalPrice float64 `json:"totalPrice,string"`
	Time       string  `json:"time"`
	Type       string  `json:"type"`
	Coord      string  `json:"coord"`
	Status     int     `json:"status,string"`
	Version    int     `json:"version,string"`

	Stations struct {
		Station []*BticLineStation `json:"station"`
	} `json:"stations"`
}

type BticLineStation struct {
	Name string `json:"name"`
	No   string `json:"no"`
	Lon  string `json:"lon"`
	Lat  string `json:"lat"`
}

type BticLineRTResp struct {
	Root *BticLineRT `json:"root"`
}

type BticLineRT struct {
	Status  int    `json:"status,string"`
	Message string `json:"message"`
	Encrypt int    `json:"encrypt,string"`
	Num     int    `json:"num,string"`
	LineID  int    `json:"lid,string"`

	Data struct {
		Bus []*BticBus `json:"bus"`
	} `json:"data"`
}

type BticBus struct {
	GT                  string `json:"gt"`
	ID                  string `json:"id"`
	T                   string `json:"t"`
	NextStationName     string `json:"ns"`
	NextStationNum      string `json:"nsn"`
	NextStationDistance string `json:"nsd"`
	NextStationArrTime  string `json:"nst"`
	StationDistance     string `json:"sd"`
	StationArrTime      string `json:"st"`
	Lat                 string `json:"x"`
	Lon                 string `json:"y"`
}

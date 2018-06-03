package guangzhou

type Response struct {
	ErrCode int         `json:"retCode"`
	Message string      `json:"retMsg"`
	Result  interface{} `json:"retData"`
}

type SearchResult struct {
	Routes   []Line          `json:"route"`
	Stations []SearchStation `json:"station"`
}

type Line struct {
	Name    string `json:"n"`
	StartSn string `json:"start"`
	EndSn   string `json:"end"`
	RouteId string `json:"i"`
}

type SearchStation struct {
	Id    string `json:"1006702"`
	Name  string `json:"n"`
	Count int    `json:"c"`
}

type BusLineResult struct {
	BusLine        *Busline         `json:"rb"`
	RunningBusInfo []RunningBusInfo `json:"runb"`
}

type Busline struct {
	OrgName   string    `json:"organName"`
	FirstTime string    `json:"ft"`
	LastTime  string    `json:"lt"`
	RouteName string    `json:"rt"`
	Stations  []Station `json:"l"`
}

type RunningBusInfo struct {
	Bl  []RunningBus `json:"bl"`
	Bbl []RunningBus `json:"bbl"`
}

type RunningBus struct {
	SId string `json:"si"`
	Id  string `json:"i"`
	Sub string `json:"sub"`

	No  string  `json:"no"`
	Lat float64 `json:"la,omitempty,string"`
	Lon float64 `json:"lo,omitempty,string"`
	En  string  `json:"en"`
	Ec  string  `json:"ec"`
}

type Station struct {
	Name string  `json:"n"`
	Lat  float64 `json:"lat,omitempty,string"`
	Lon  float64 `json:"lon,omitempty,string"`
	SId  string  `json:"si"`
	Id   string  `json:"i"`

	IsSubway int          `json:"sw,string"`
	Subways  []SubwayInfo `json:"sinfo"`
	IsBrt    int          `json:"brt,string"`
	Order    string       `json:"order"`
}

type SubwayInfo struct {
	Color string `json:"color"`
	Name  string `json:"name"`
}

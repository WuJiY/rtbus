package jinan

type SearchResponse struct {
	Status ResponseStatus `json:"status"`
	Result struct {
		Result []*SearchResponseResult `json:"result"`
	} `json:"result"`
}

type ResponseStatus struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type SearchResponseResult struct {
	ID               string `json:"id"`
	LineName         string `json:"lineName"`
	StartStationName string `json:"startStationName"`
	EndStationName   string `json:"endStationName"`
	UpdateTime       string `json:"updateTime"`
}

type BusLineResponse struct {
	Status ResponseStatus        `json:"status"`
	Result BusLineResponseResult `json:"result"`
}

type BusLineResponseResult struct {
	ID               string           `json:"id"`
	Area             int              `json:"area"`
	LineName         string           `json:"lineName"`
	StartStationName string           `json:"startStationName"`
	EndStationName   string           `json:"endStationName"`
	TicketPrice      string           `json:"ticketPrice"`
	OperationTime    string           `json:"operationTime"`
	LinePoints       string           `json:"linePoints"`
	UpdateTime       string           `json:"updateTime"`
	State            string           `json:"state"`
	Stations         []BuslineStation `json:"stations"`
}

type BuslineStation struct {
	ID          string   `json:"id"`
	Area        int      `json:"area"`
	StationName string   `json:"stationName"`
	Lng         float64  `json:"lng"`
	Lat         float64  `json:"lat"`
	Buslines    string   `json:"buslines"`
	State       string   `json:"state"`
	UpdateTime  string   `json:"updateTime"`
	Distance    float64  `json:"distance"`
	busLineList []string `json:"busLineList"`
}

type RunningBusResponse struct {
	Status ResponseStatus              `json:"status"`
	Result []*RunningBusResponseResult `json:"result"`
}
type RunningBusResponseResult struct {
	BusID           string  `json:"busId"`
	Lng             float64 `json:"lng"`
	Lat             float64 `json:"lat"`
	Velocity        float64 `json:"velocity"`
	IsArrvLft       string  `json:"isArrvLft"`
	StationSeqNum   int     `json:"stationSeqNum"`
	Status          string  `json:"status"`
	BuslineId       string  `json:"buslineId"`
	ActTime         string  `json:"actTime"`
	CardId          string  `json:"cardId"`
	OrgName         string  `json:"orgName"`
	AverageVelocity float64 `json:"averageVelocity"`
	Coordinate      int     `json:"coordinate"`
}

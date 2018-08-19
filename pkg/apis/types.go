package apis

type Bus struct {
	Order    int     `json:"o"`
	Status   float64 `json:"s"`
	Location string  `json:"l"`
	Distance int     `json:"d"`
}

type StationBus struct {
	StationName      string        `json:"sn"`
	LineCount        int           `json:"lc"`
	SupportLineCount int           `json:"slc"`
	Location         string        `json:"l"`
	Lines            []BaseLineDir `json:"lines"`
}

type BaseLineDir struct {
	No         string `json:"no"`
	Dir        string `json:"dir"`
	AnotherDir string `json:"adir"`
	StartSn    string `json:"ssn"`
	EndSn      string `json:"esn"`
	NextSn     string `json:"nsn,omitempty"`
	Order      int    `json:"o,omitempty"`
	Buses      []Bus  `json:"buses"`
}

type LineDirection struct {
	*BaseLineDir

	ID        string `json:"id"`
	Price     string `json:"price"`
	FirstTime string `json:"ftime"`
	LastTime  string `json:"ltime"`

	Stations []Station `json:"ss"`
}

type Station struct {
	Order    int    `json:"o"`
	Sn       string `json:"sn"`
	Location string `json:"l"`
}

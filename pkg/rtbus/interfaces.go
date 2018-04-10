package rtbus

type CityRTBusApi interface {
	City() *CityInfo
	GetBusLine(lineno string) (bl *BusLine, err error)
	GetBusLineDir(lineno, dirname string) (bdi *BusDirInfo, err error)
	GetRunningBus(lineno, dirname string) (rbus []*RunningBus, err error)
}

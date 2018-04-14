package rtbus

type CityRTBusApi interface {
	City() *CityInfo
	Search(keyword string) ([]*BusDirInfo, error)
	GetBusLine(lineno string, with_running_bus bool) (bl *BusLine, err error)
	GetBusLineDir(lineno, dirname string) (bdi *BusDirInfo, err error)
	GetRunningBus(lineno, dirname string) (rbus []*RunningBus, err error)
}

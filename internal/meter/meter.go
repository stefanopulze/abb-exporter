package meter

import "time"

type Meter interface {
	Name() string
	QueryInfo() (*Info, error)
	QueryUsageStatus() (*UsageStatus, error)
	QueryTotalActiveImport() (float64, error)
	Tags() *map[string]string
}

type Info struct {
	Serial               uint32   `json:"serial"`
	FirmwareVersion      string   `json:"firmwareVersion"`
	ModbusMappingVersion Firmware `json:"modbusMappingVersion"`
	TypeDesignation      string   `json:"typeDesignation"`
}

type Firmware struct {
	Major uint16 `json:"major"`
	Minor uint16 `json:"minor"`
}

type UsageStatus struct {
	Timestamp   time.Time
	Voltage     float64 `json:"voltage"`
	Current     float64 `json:"current"`
	ActivePower float64 `json:"activePower"`
	Frequency   float64 `json:"frequency"`
}

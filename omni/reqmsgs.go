package omni

//go:generate stringer -type=SystemTrouble
//go:generate stringer -type=SystemFeature
//go:generate stringer -type=TempFormat
//go:generate stringer -type=TimeFormat
//go:generate stringer -type=DateFormat

type SystemTrouble uint8

const (
	Freeze SystemTrouble = 1 + iota
	BatteryLow
	ACPower
	PhoneLine
	DigitalCommunicator
	Fuse
	Freeze2
	BatteryLow2 // Listed twice in spec
)

type SystemFeature uint8

const (
	NuVoConcerto SystemFeature = 1 + iota
	NuVoEssentia
	NuVoGrandConcerto
	Russound
	HAIHiFi
	XanTech
	SpeakerCraft
	Proficient
	DSCSecurity
)

type TempFormat uint8

const (
	Fahrenheit TempFormat = 1 + iota
	Celsius
)

type TimeFormat uint8

const (
	Twelve TimeFormat = 1 + iota
	TwentyFour
)

type DateFormat uint8

const (
	MMDD DateFormat = 1 + iota
	DDMM
)

type SystemInfo struct {
	ModelNumber      uint8
	MajorVersion     uint8
	MinorVersion     uint8
	Revision         uint8
	LocalPhoneNumber [25]byte
}

type SystemStatus struct {
	DateValid   uint8
	Year        uint8
	Month       uint8
	Day         uint8
	DayOfWeek   uint8
	Hour        uint8
	Minute      uint8
	Second      uint8
	Daylight    uint8
	SunriseHour uint8
	SunriseMin  uint8
	SunsetHour  uint8
	SunsetMin   uint8
	Battery     uint8
}

type SystemTroubles struct {
	Troubles []SystemTrouble
}

type SystemFeatures struct {
	Features []SystemFeature
}
type SystemFormats struct {
	TempFormat TempFormat
	TimeFormat TimeFormat
	DateFormat DateFormat
}

type ObjectTypeCapacities struct {
	CapacityType ObjectType
	CapacityMSB  uint8
	CapacityLSB  uint8
}

type ObjectProperties struct {
	ObjectType uint8
}

type ObjectStatus struct {
}

type AudioSourceStatus struct {
}

type ZoneReadyStatus struct {
}

type ConnectedSecuritySystemStatus struct {
}

func omniTempToF(otemp uint8) float64 {
	celsius := -40.0 + (float64(otemp) / 2.0)
	return celsius*1.8 + 32
}

package omni

//go:generate stringer -type=SystemTrouble
//go:generate stringer -type=SystemFeature
//go:generate stringer -type=TempFormat
//go:generate stringer -type=TimeFormat
//go:generate stringer -type=DateFormat
//go:generate stringer -type=ObjectType

type SystemTrouble uint8

const (
	Freeze SystemTrouble = 1 + iota
	BatteryLow
	ACPower
	PhoneLine
	DigitalCommunicator
	Fuse
	Freeze2 // These are listed twice in spec
	BatteryLow2
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

type ObjectType uint8

const (
	Zone ObjectType = 1 + iota
	Unit
	Button
	Code
	Area
	Thermostat
	Message
	AuxilarySensor
	AudioSource
	AudioZone
	ExpansionEnclosure
	Console
	UserSetting
	AccessControlReader
	AccessControlLock
)

type SystemInfo struct {
	ModelNumber      uint8
	MajorVersion     uint8
	MinorVerison     uint8
	Revesion         uint8
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
}

type ObjectStatus struct {
}

type AudioSourceStatus struct {
}

type ZoneReadyStatus struct {
}

type ConnectedSecuritySystemStatus struct {
}

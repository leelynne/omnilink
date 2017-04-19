package omni

//go:generate stringer -type=SystemTrouble
//go:generate stringer -type=SystemFeature
//go:generate stringer -type=TempFormat
//go:generate stringer -type=TimeFormat
//go:generate stringer -type=DateFormat

type SystemTrouble uint8
type SystemFeature uint8
type TempFormat uint8
type TimeFormat uint8
type DateFormat uint8

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
const (
	Fahrenheit TempFormat = 1 + iota
	Celsius
)
const (
	Twelve TimeFormat = 1 + iota
	TwentyFour
)

const (
	MMDD DateFormat = 1 + iota
	DDMM
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
}

func test() {

}

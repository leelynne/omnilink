package omni

//go:generate stringer -type=ObjectType

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

	StatusSizeThermostat = 9
)

var StatusSizes = map[ObjectType]int{
	Zone:       4,
	Unit:       5,
	Thermostat: 9,
}

// Object and Property messages for each Object type. Structs match the byte layout specified in the protocol

type ThermostatProperties struct {
	ObjectType         uint8
	NumberMSB          uint8
	NumberLSB          uint8
	Communicating      uint8
	Temperature        uint8
	HeatSetPoint       uint8
	CoolSetPoint       uint8
	SystemMode         uint8
	FanMode            uint8
	HoldStatus         uint8
	Type               uint8
	Name               [12]byte
	Humidty            uint8
	HumidifySetPoint   uint8
	DehumidifySetPoint uint8
	OutdoorTemperature uint8
	ActionStatus       uint8
}

type ThermostatStatus struct {
	NumberMSB    uint8
	NumberLSB    uint8
	Status       uint8
	CurrentTemp  uint8
	HeatSetPoint uint8
	CoolSetPoint uint8
	SystemMode   uint8
	FanMode      uint8
	HoldStatus   uint8
}

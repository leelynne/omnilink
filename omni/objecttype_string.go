// Code generated by "stringer -type=ObjectType"; DO NOT EDIT.

package omni

import "fmt"

const _ObjectType_name = "ZoneUnitButtonCodeAreaThermostatMessageAuxilarySensorAudioSourceAudioZoneExpansionEnclosureConsoleUserSettingAccessControlReaderAccessControlLock"

var _ObjectType_index = [...]uint8{0, 4, 8, 14, 18, 22, 32, 39, 53, 64, 73, 91, 98, 109, 128, 145}

func (i ObjectType) String() string {
	i -= 1
	if i >= ObjectType(len(_ObjectType_index)-1) {
		return fmt.Sprintf("ObjectType(%d)", i+1)
	}
	return _ObjectType_name[_ObjectType_index[i]:_ObjectType_index[i+1]]
}

// Code generated by "stringer -type=SystemTrouble"; DO NOT EDIT.

package omni

import "fmt"

const _SystemTrouble_name = "FreezeBatteryLowACPowerPhoneLineDigitalCommunicatorFuseFreeze2BatteryLow2"

var _SystemTrouble_index = [...]uint8{0, 6, 16, 23, 32, 51, 55, 62, 73}

func (i SystemTrouble) String() string {
	i -= 1
	if i >= SystemTrouble(len(_SystemTrouble_index)-1) {
		return fmt.Sprintf("SystemTrouble(%d)", i+1)
	}
	return _SystemTrouble_name[_SystemTrouble_index[i]:_SystemTrouble_index[i+1]]
}

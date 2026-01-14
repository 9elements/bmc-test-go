// Package ipmimct holds all information an structures for the IPMI OEM Extension of MCT.
package ipmimct

import (
	"fmt"
	"strings"
)

// The constants describing values to send via ipmi twitter or mct oem extension.
const (
	FanPWMDuty           uint8 = 0x05
	ManufactureMode      uint8 = 0x06
	FloorDuty            uint8 = 0x07
	SetFruField          uint8 = 0x0B
	GetFruField          uint8 = 0x0C
	SetService           uint8 = 0x0D
	GetService           uint8 = 0x0E
	GetFirmwareString    uint8 = 0x10
	ConfigECCLeakyBucket uint8 = 0x1A
	ClearCMOS            uint8 = 0x3A
	GPIOStatus           uint8 = 0x41
	PnmGetReading        uint8 = 0xE2

	RandomDelayACRestorePowerOn uint8 = 0x18
	GetPostCode                 uint8 = 0x10
	RelinkLan                   uint8 = 0x12
	SetAMDSMBusOwner            uint8 = 0x24
)

func tyanCommand(cmd uint8, data []uint8) string {
	tyanManufacturerID := []byte{0xFD, 0x19, 0x00}

	const mctModuleCode uint8 = 0x2e

	ret := make([]uint8, 0)
	ret = append(ret, mctModuleCode)
	ret = append(ret, cmd)
	ret = append(ret, tyanManufacturerID...)

	return cmdToString(append(ret, data...))
}

func twitterCommand(cmd uint8, data []uint8) string {
	const twitterModuleCode uint8 = 0x30

	ret := make([]uint8, 0)
	ret = append(ret, twitterModuleCode)
	ret = append(ret, cmd)

	return cmdToString(append(ret, data...))
}

func cmdToString(cmd []uint8) string {
	var ret strings.Builder
	for _, item := range cmd {
		fmt.Fprintf(&ret, "0x%x ", item)
	}

	return ret.String()
}

// MCTFanDuty creates the byte sequence to query fan duty values via IPMI.
func MCTFanDuty(data []uint8) string {
	return tyanCommand(FanPWMDuty, data)
}

// MCTManufactureMode creates the byte sequence to query manufacture mode status via IPMI.
func MCTManufactureMode(data []uint8) string {
	return tyanCommand(ManufactureMode, data)
}

// MCTFloorDuty creates the byte sequence to query floor duty values via IPMI.
func MCTFloorDuty(data []uint8) string {
	return tyanCommand(FloorDuty, data)
}

// MCTGetFruField creates the byte sequence to query fru field values via IPMI.
func MCTGetFruField(data []uint8) string {
	return tyanCommand(GetFruField, data)
}

// MCTGetFirmwareString creates the byte sequence to query firmware version string via IPMI.
func MCTGetFirmwareString() string {
	return tyanCommand(GetFirmwareString, []byte{})
}

// MCTConfigECCLeakyBucket creates the byte sequence to query configuration of ECC Leaky Bucket values via IPMI.
func MCTConfigECCLeakyBucket(data []uint8) string {
	return tyanCommand(ConfigECCLeakyBucket, data)
}

// MCTGPIOStatus creates the byte sequence to query GPIO status values via IPMI.
func MCTGPIOStatus(data []uint8) string {
	return tyanCommand(GPIOStatus, data)
}

// TwitterSetService creates the byte sequence to set service values via IPMI Twitter OEM Extension.
func TwitterSetService(data []uint8) string {
	return twitterCommand(SetService, data)
}

// TwitterGetService creates the byte sequence to get service values via IPMI Twitter OEM Extension.
func TwitterGetService(data []uint8) string {
	return twitterCommand(GetService, data)
}

// TwitterPnmGetReading creates the byte sequence to get Pnm readings via IPMI Twitter OEM Extension.
func TwitterPnmGetReading(data []uint8) string {
	return twitterCommand(PnmGetReading, data)
}

// TwitterRandomDelayACRestorePowerOn creates the byte sequence to set
// random Delay AC Restore Power On via IPMI Twitter OEM Extension.
func TwitterRandomDelayACRestorePowerOn(data []uint8) string {
	return twitterCommand(RandomDelayACRestorePowerOn, data)
}

// TwitterGetPostCodes creates the byte sequence to query Post codes via IPMI Twitter OEM Extension.
func TwitterGetPostCodes() string {
	return twitterCommand(GetPostCode, []uint8{})
}

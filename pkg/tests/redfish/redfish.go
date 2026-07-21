// Package redfish groups all redfish related tests
package redfish

import "github.com/9elements/bmc-test-go/pkg/framework"

// GetTests returns all SMBIOS tests
//
//nolint:funlen
func GetTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "Redfish Chassis Test",
			ShortName: "redfish-00",
			Function:  testRedfishChassis,
		},
		{
			Name:      "Redfish System Test",
			ShortName: "redfish-01",
			Function:  testRedfishSystem,
		},

		{
			Name:      "Redfish Firmware Version",
			ShortName: "redfish-02",
			Function:  testRedfishGetFirmwareVersion,
		},

		{
			Name:      "Redfish Validate PSU Information",
			ShortName: "redfish-03",
			Function:  testRedfishCheckPowerSupplyInformation,
		},

		{
			Name:      "Redfish PSU Voltage Sensor Names",
			ShortName: "redfish-04",
			Function:  testRedfishVoltagesSensorNames,
		},

		{
			Name:      "Redfish PSU Voltage Sensor Values",
			ShortName: "redfish-05",
			Function:  testRedfishVoltageSensorValueThesholds,
		},

		{
			Name:      "Redfish Temperatur Sensor Names",
			ShortName: "redfish-06",
			Function:  testRedfishTemperatureSensorNames,
		},

		{
			Name:      "Redfish Temperatur Sensor Values",
			ShortName: "redfish-07",
			Function:  testRedfishTemperaturSensorValueThesholds,
		},

		{
			Name:      "Redfish Fan Sensor Names",
			ShortName: "redfish-08",
			Function:  testRedfishThermalFanNames,
		},

		{
			Name:      "Redfish Fan Sensor Values",
			ShortName: "redfish-09",
			Function:  testRedfishThermalFanValues,
		},

		{
			Name:      "Redfish Memory",
			ShortName: "redfish-10",
			Function:  testRedfishMemory,
		},

		{
			Name:      "Redfish Processor",
			ShortName: "redfish-11",
			Function:  testRedfishProcessor,
		},

		{
			Name:      "Redfish Firmware Downgrade Update",
			ShortName: "redfish-12",
			Function:  testRedfishFirmwareDowngradeUpdate,
		},
	}
}

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
			Status:    framework.StatusImplemented,
			Function:  testRedfishChassis,
		},
		{
			Name:      "Redfish System Test",
			ShortName: "redfish-01",
			Status:    framework.StatusImplemented,
			Function:  testRedfishSystem,
		},

		{
			Name:      "Redfish Firmware Version",
			ShortName: "redfish-02",
			Status:    framework.StatusImplemented,
			Function:  testRedfishGetFirmwareVersion,
		},

		{
			Name:      "Redfish Validate PSU Information",
			ShortName: "redfish-03",
			Status:    framework.StatusImplemented,
			Function:  testRedfishCheckPowerSupplyInformation,
		},

		{
			Name:      "Redfish PSU Voltage Sensor Names",
			ShortName: "redfish-04",
			Status:    framework.StatusImplemented,
			Function:  testRedfishVoltagesSensorNames,
		},

		{
			Name:      "Redfish PSU Voltage Sensor Values",
			ShortName: "redfish-05",
			Status:    framework.StatusImplemented,
			Function:  testRedfishVoltageSensorValueThesholds,
		},

		{
			Name:      "Redfish Temperatur Sensor Names",
			ShortName: "redfish-06",
			Status:    framework.StatusImplemented,
			Function:  testRedfishTemperatureSensorNames,
		},

		{
			Name:      "Redfish Temperatur Sensor Values",
			ShortName: "redfish-07",
			Status:    framework.StatusImplemented,
			Function:  testRedfishTemperaturSensorValueThesholds,
		},

		{
			Name:      "Redfish Fan Sensor Names",
			ShortName: "redfish-08",
			Status:    framework.StatusImplemented,
			Function:  testRedfishThermalFanNames,
		},

		{
			Name:      "Redfish Fan Sensor Values",
			ShortName: "redfish-09",
			Status:    framework.StatusImplemented,
			Function:  testRedfishThermalFanValues,
		},

		{
			Name:      "Redfish Memory",
			ShortName: "redfish-10",
			Status:    framework.StatusImplemented,
			Function:  testRedfishMemory,
		},

		{
			Name:      "Redfish Processor",
			ShortName: "redfish-11",
			Status:    framework.StatusImplemented,
			Function:  testRedfishProcessor,
		},

		{
			Name:      "Redfish Firmware Downgrade Update",
			ShortName: "redfish-12",
			Status:    framework.StatusImplemented,
			Function:  testRedfishFirmwareDowngradeUpdate,
		},
	}
}

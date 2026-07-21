// Package smbios implements tests for SMBIOS properties.
package smbios

import (
	"github.com/9elements/bmc-test-go/pkg/framework"
)

// GetTests returns all smbios tests as slice of smbios.Test.
func GetTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "Test SMBIOS Baseboard System Vendor",
			ShortName: "smbios-00",
			Function:  testSMBIOSBaseboardSystemVendor,
		},

		{
			Name:      "Test SMBIOS Baseboard Product Name",
			ShortName: "smbios-01",
			Function:  testSMBIOSBaseboardProductName,
		},

		{
			Name:      "Test SMBIOS Baseboard Product Version",
			ShortName: "smbios-02",
			Function:  testSMBIOSBaseboardProductVersion,
		},

		{
			Name:      "Test SMBIOS Baseboard Product Serial",
			ShortName: "smbios-03",
			Function:  testSMBIOSBaseboardProductSerial,
		},

		{
			Name:      "Test SMBIOS Baseboard Board Vendor",
			ShortName: "smbios-04",
			Function:  testSMBIOSBaseboardBoardVendor,
		},

		{
			Name:      "Test SMBIOS Baseboard Board Name",
			ShortName: "smbios-05",
			Function:  testSMBIOSBaseboardBoardName,
		},

		{
			Name:      "Test SMBIOS Baseboard Board Serial",
			ShortName: "smbios-06",
			Function:  testSMBIOSBaseboardBoardSerial,
		},
	}
}

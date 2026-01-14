// Package misc includes tests of various properties not fitting into other packages (yet)
package misc

import (
	"github.com/9elements/bmc-test-go/pkg/framework"
)

// GetTests returns all misc tests.
func GetTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "Interface USB network device",
			ShortName: "iface-00",
			Status:    framework.StatusImplemented,
			Function:  testUSBNetwork,
		},
		{
			Name:      "Test PCI Device Aspeed",
			ShortName: "pci-00",
			Status:    framework.StatusImplemented,
			Function:  testPCIDeviceAspeed,
		},
		{
			Name:      "SD Card Service",
			ShortName: "misc-02",
			Status:    framework.StatusImplemented,
			Function:  testSDCardDriver,
		},
		{
			Name:      "SD Card Device presence",
			ShortName: "misc-03",
			Status:    framework.StatusImplemented,
			Function:  testSDCardDevicePresent,
		},
		{
			Name:      "SD Card File System Operation",
			ShortName: "misc-04",
			Status:    framework.StatusImplemented,
			Function:  testSDCardFilesystemOps,
		},
	}
}

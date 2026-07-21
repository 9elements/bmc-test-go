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
			Function:  testUSBNetwork,
		},
		{
			Name:      "Test PCI Device Aspeed",
			ShortName: "pci-00",
			Function:  testPCIDeviceAspeed,
		},
		{
			Name:      "SD Card Service",
			ShortName: "misc-02",
			Function:  testSDCardDriver,
		},
		{
			Name:      "SD Card Device presence",
			ShortName: "misc-03",
			Function:  testSDCardDevicePresent,
		},
		{
			Name:      "SD Card File System Operation",
			ShortName: "misc-04",
			Function:  testSDCardFilesystemOps,
		},
	}
}

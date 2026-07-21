// Package bmclinux groups all tests for the bmc specific linux system.
package bmclinux

import (
	"github.com/9elements/bmc-test-go/pkg/framework"
)

// GetBMCLinuxTests returns all BMC linux tests.
func GetBMCLinuxTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "BMC Linux Startup health Test",
			ShortName: "bmc-linux-00",
			Function:  testBMCLinuxStartupHealth,
		},
		{
			Name:      "BMC Linux Systemctl failed Test",
			ShortName: "bmc-linux-01",
			Function:  testBMCLinuxSystemctlFailed,
		},
		{
			Name:      "BMC Linux Systemctl job list Test",
			ShortName: "bmc-linux-02",
			Function:  testBMCLinuxSystemctlJobList,
		},
	}
}

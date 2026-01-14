// Package tests provides the functions to get all tests, test suites or single tests.
package tests

import (
	"log"

	"github.com/9elements/bmc-test-go/pkg/framework"
	bmclinux "github.com/9elements/bmc-test-go/pkg/tests/bmc_linux"
	"github.com/9elements/bmc-test-go/pkg/tests/ipmi"
	"github.com/9elements/bmc-test-go/pkg/tests/misc"
	"github.com/9elements/bmc-test-go/pkg/tests/redfish"
	"github.com/9elements/bmc-test-go/pkg/tests/smbios"
)

// AllTests collects all tests from all test packages and returns them as a slice of framework.Test.
func AllTests() []*framework.Test {
	size := len(ipmi.GetTests()) + len(ipmi.GetIPMIMCTTests()) + len(ipmi.GetIPMITwitterTests()) +
		len(redfish.GetTests()) + len(misc.GetTests()) + len(smbios.GetTests()) + len(bmclinux.GetBMCLinuxTests())

	ret := make([]*framework.Test, 0, size)

	ret = append(ret, ipmi.GetTests()...)

	ret = append(ret, ipmi.GetIPMIMCTTests()...)

	ret = append(ret, ipmi.GetIPMITwitterTests()...)

	ret = append(ret, redfish.GetTests()...)

	ret = append(ret, misc.GetTests()...)

	ret = append(ret, smbios.GetTests()...)

	ret = append(ret, bmclinux.GetBMCLinuxTests()...)

	return ret
}

// GetSuite collects the tests of the given suite name.
func GetSuite(suite string) []*framework.Test {
	suites := map[string]func() []*framework.Test{
		"ipmi":         ipmi.GetTests,
		"ipmi-mct":     ipmi.GetIPMIMCTTests,
		"ipmi-twitter": ipmi.GetIPMITwitterTests,
		"redfish":      redfish.GetTests,
		"misc":         misc.GetTests,
		"smbios":       smbios.GetTests,
		"bmc-linux":    bmclinux.GetBMCLinuxTests,
	}

	ret, ok := suites[suite]
	if !ok {
		log.Printf("unknown suite name: %s", suite)

		return nil
	}

	return ret()
}

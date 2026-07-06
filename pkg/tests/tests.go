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
	ret := []*framework.Test{}

	for _, f := range Suites {
		ret = append(ret, f()...)
	}

	return ret
}

// Suites is the map of all available test suites
var Suites = map[string]func() []*framework.Test{
	"ipmi":         ipmi.GetTests,
	"ipmi-mct":     ipmi.GetIPMIMCTTests,
	"ipmi-twitter": ipmi.GetIPMITwitterTests,
	"redfish":      redfish.GetTests,
	"misc":         misc.GetTests,
	"smbios":       smbios.GetTests,
	"bmc-linux":    bmclinux.GetBMCLinuxTests,
}

// GetSuite collects the tests of the given suite name.
func GetSuite(suite string) []*framework.Test {
	ret, ok := Suites[suite]
	if !ok {
		log.Printf("unknown suite name: %s", suite)

		return nil
	}

	return ret()
}

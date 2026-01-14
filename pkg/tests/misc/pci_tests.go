package misc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

var errPCIBridgeNotFound = errors.New("expected PCI bridge from ASPEED Technology Inc. not found")

func testPCIDeviceAspeed(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("lspci", " | grep ASPEED", "|| true")
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), "PCI bridge: ASPEED Technology, Inc") {
			return false, fmt.Errorf("%w: on device: %s in string: %s", errPCIBridgeNotFound,
				host.Name, string(resp)), nil
		}
	}

	return true, nil, nil
}

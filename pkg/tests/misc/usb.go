package misc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

var errNoUSBNetwork = errors.New("no usb network device found")

func testUSBNetwork(dev *testdevice.Device, _ *configuration.Config) (bool, error, error) {
	resp, err := dev.ExecuteCommandLine("ip", "link", "show | grep usb || true")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "usb") {
		return false, errNoUSBNetwork, nil
	}

	return true, nil, nil
}

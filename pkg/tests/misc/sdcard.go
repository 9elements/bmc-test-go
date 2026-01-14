package misc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

var errServiceFailed = errors.New("service failed")

func testSDCardDriver(dev *testdevice.Device, _ *configuration.Config) (bool, error, error) {
	resp, err := dev.ExecuteCommandLine("systemctl", "list-unit-files", "--type=service", "|", "grep",
		"journal-to-sdcard.service", "|| true")
	if err != nil {
		return false, nil, fmt.Errorf("systemctl. Output: %s \n with error: %w", string(resp), err)
	}

	if strings.Contains(string(resp), "failed") {
		return false, fmt.Errorf("%w journal-to-sdcard.service: %s", errServiceFailed, string(resp)), nil
	}

	return true, nil, nil
}

var errSDCardDeviceMissing = errors.New("/dev/mmcblk0 not found")

func testSDCardDevicePresent(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ls", "-la", "/dev", "|", "grep", "mmcblk0", "|| true")
	if err != nil {
		return false, nil, fmt.Errorf("ls /dev failed. Output: %s \n with error: %w", string(resp), err)
	}

	if !strings.Contains(string(resp), "mmcblk0") {
		return false, errSDCardDeviceMissing, nil
	}

	return true, nil, nil
}

func testSDCardFilesystemOps(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("mkdir", "sdcard")
	if err != nil {
		return false, nil, fmt.Errorf("creating dir 'sdcard' failed. Output: %s \n with error: %w", string(resp), err)
	}

	resp, err = dev.ExecuteCommandLine("mound", "/dev/mmcblk0", "sdcard")
	if err != nil {
		return false, nil, fmt.Errorf("mounting sdcard failed. Output: %s \n with error: %w", string(resp), err)
	}

	return true, nil, nil
}

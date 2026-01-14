package bmclinux

import (
	"errors"
	"fmt"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

var (
	errConntentReturned = errors.New("content returned")
	errSystemctlFailed  = errors.New("systemctl --failed")
)

func testBMCLinuxStartupHealth(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("dmesg", "|", "grep", "fail", "||", "true")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if string(resp) != "" {
		return false, fmt.Errorf("dmesg %w: %s", errConntentReturned, string(resp)), nil
	}

	return true, nil, nil
}

func testBMCLinuxSystemctlFailed(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("systemctl", "--failed")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if strings.Contains(string(resp), "failed") {
		return false, fmt.Errorf("%w:\n%s", errSystemctlFailed, string(resp)), nil
	}

	return true, nil, nil
}

func testBMCLinuxSystemctlJobList(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("systemctl", "list-jobs")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "No jobs running.") {
		return false, fmt.Errorf("%w:\n%s", errSystemctlFailed, string(resp)), nil
	}

	return true, nil, nil
}

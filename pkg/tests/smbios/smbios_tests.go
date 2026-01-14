package smbios

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

const smbiosBasePath = "/sys/class/dmi/id/"

var errPropertyMismatch = errors.New("property mismatch")

func testSMBIOSBaseboardSystemVendor(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "sys_vendor")
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].ProductManufacturer) {
			return false, fmt.Errorf("%w: sys_vendor string: %s not equal board_manufacturer: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].ProductManufacturer), nil
		}
	}

	return true, nil, nil
}

func testSMBIOSBaseboardProductName(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "product_name")
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].ProductName) {
			return false, fmt.Errorf("%w: product_name string: %s not equal product_name: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].ProductName), nil
		}
	}

	return true, nil, nil
}

func testSMBIOSBaseboardProductSerial(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "product_serial")
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].ProductSerial) {
			return false, fmt.Errorf("%w: product_name string: %s not equal product_name: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].ProductSerial), nil
		}
	}

	return true, nil, nil
}

func testSMBIOSBaseboardProductVersion(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "product_version")
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].ProductVersion) {
			return false, fmt.Errorf("%w: product_name string: %s not equal product_name: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].ProductVersion), nil
		}
	}

	return true, nil, nil
}

func testSMBIOSBaseboardBoardVendor(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "board_vendor")
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].BoardManufacturer) {
			return false, fmt.Errorf("%w: sys_vendor string: %s not equal board_manufacturer: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].BoardManufacturer), nil
		}
	}

	return true, nil, nil
}

func testSMBIOSBaseboardBoardName(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "board_name")
	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].BoardProduct) {
			return false, fmt.Errorf("%w: product_name string: %s not equal product_name: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].BoardProduct), nil
		}
	}

	return true, nil, nil
}

func testSMBIOSBaseboardBoardSerial(dev *testdevice.Device, cfg *configuration.Config) (bool, error, error) {
	infoPath := path.Join(smbiosBasePath, "board_serial")

	if len(dev.Hosts) < 1 {
		return false, nil, nil
	}

	for _, host := range dev.Hosts {
		resp, err := host.ExecuteCommandLine("cat", infoPath)
		if err != nil {
			return false, nil, fmt.Errorf("command execution returned error: %w", err)
		}

		if !strings.Contains(string(resp), cfg.Fru["baseboard"].BoardSerial) {
			return false, fmt.Errorf("%w: product_name string: %s not equal product_name: %s on host: %s", errPropertyMismatch,
				string(resp), cfg.Fru["baseboard"].BoardSerial, host.Name), nil
		}
	}

	return true, nil, nil
}

// Package ipmi includes all tests if IPMI api and OEM Extensions.
package ipmi

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/framework"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

// GetTests returns all defined IPMI tests.
//
//nolint:funlen
func GetTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "IPMI Sensor Voltage",
			ShortName: "ipmi-00",
			Status:    framework.StatusImplemented,
			Function:  testIPMISensorsVoltage,
		},

		{
			Name:      "IPMI Sensor Temperature",
			ShortName: "ipmi-01",
			Status:    framework.StatusImplemented,
			Function:  testIPMISensorTemperature,
		},

		{
			Name:      "IPMI Sensor Fans",
			ShortName: "ipmi-02",
			Status:    framework.StatusImplemented,
			Function:  testIPMISensorFans,
		},

		{
			Name:      "IPMI Power Status",
			ShortName: "ipmi-03",
			Status:    framework.StatusImplemented,
			Function:  testIPMIPowerStatus,
		},
		{
			Name:      "IPMI Chassis Uptime",
			ShortName: "ipmi-04",
			Status:    framework.StatusImplemented,
			Function:  testIPMIChassisUptime,
		},
		{
			Name:      "IPMI Fru",
			ShortName: "ipmi-05",
			Status:    framework.StatusImplemented,
			Function:  testIPMIFru,
		},
		{
			Name:      "IPMI DCMI Power Reading",
			ShortName: "ipmi-06",
			Status:    framework.StatusImplemented,
			Function:  testIPMIDCMIPowerReading,
		},
		{
			Name:      "IPMI Watchdog Configuration",
			ShortName: "ipmi-07",
			Status:    framework.StatusImplemented,
			Function:  testIPMIWatchdogConfiguration,
		},
		{
			Name:      "IPMI System GUID",
			ShortName: "ipmi-08",
			Status:    framework.StatusImplemented,
			Function:  testIPMISystemGUID,
		},
		{
			Name:      "IPMI DCMI Power reading",
			ShortName: "ipmi-9",
			Status:    framework.StatusImplemented,
			Function:  nil,
		},
	}
}

// For precompilation. Used in testIPMIChassisUptime.
var pattern = regexp.MustCompile(`POH Counter  :\s+\d+\s*days,\s+\d+\s*hours`)

func stringToIpmiSensorName(sFull string) string {
	// Since the names come from dbus, they do contain underscores.
	// We found the sensor if we find the name with or without underscores.
	sUnderscore := strings.ReplaceAll(sFull, " ", "_")

	// They have this weird stuff where they shorten sensor names with abbreviations
	// in the ipmi stack, because of string length limitation.
	// So we have to match on that aswell.
	sAbbreviated := strings.ReplaceAll(strings.ReplaceAll(sUnderscore, "Output", "Out"), "Input", "In")

	const limit int = 16

	if len(sAbbreviated) <= limit {
		return sAbbreviated
	}

	// at most 16 chars are allowed
	return sAbbreviated[:16]
}

var (
	errLineParts     = errors.New("line part count mismatch")
	errInvalidSensor = errors.New("sensor not connected or not ok")
)

func sensorIterationHelper(resp []byte, reqs []string) (bool, error, error) {
	lines := strings.Split(string(resp), "\n")

	for _, req := range reqs {
		var found bool

		var selectedLine string

		for _, line := range lines {
			if strings.Contains(line, req) {
				found = true
				selectedLine = line
			}
		}

		if !found {
			return false, nil, nil
		}

		parts := strings.Split(selectedLine, "|")

		const exactParts int = 3

		if len(parts) != exactParts {
			return false, fmt.Errorf("%w: %s", errLineParts, selectedLine), nil
		}

		status := strings.TrimSpace(parts[2])

		if status != "ok" && status != "nc" {
			return false, fmt.Errorf("%w: %s '%s'", errInvalidSensor, req, selectedLine), nil
		}
	}

	return true, nil, nil
}

func testIPMISensorsVoltage(dev *testdevice.Device, testCfg *configuration.Config) (
	bool, error, error,
) {
	reqs := make([]string, 0, len(testCfg.Voltage))

	for _, item := range testCfg.Voltage {
		reqs = append(reqs, stringToIpmiSensorName(item))
	}

	resp, err := dev.ExecuteCommandLine("ipmitool sdr")
	if err != nil {
		log.Printf("ipmitool sdr failed: %s", string(resp))

		return false, fmt.Errorf("ipmitool sdr failed: %w", err), nil
	}

	return sensorIterationHelper(resp, reqs)
}

func testIPMISensorTemperature(dev *testdevice.Device, testCfg *configuration.Config) (
	bool, error, error,
) {
	reqs := make([]string, 0, len(testCfg.Temperatures))

	for _, item := range testCfg.Temperatures {
		reqs = append(reqs, stringToIpmiSensorName(item))
	}

	resp, err := dev.ExecuteCommandLine("ipmitool sdr")
	if err != nil {
		log.Printf("ipmitool sdr failed: %s", string(resp))

		return false, fmt.Errorf("ipmitool sdr failed: %w", err), nil
	}

	return sensorIterationHelper(resp, reqs)
}

func testIPMISensorFans(dev *testdevice.Device, testCfg *configuration.Config) (
	bool, error, error,
) {
	reqs := make([]string, 0, len(testCfg.Fans))

	for _, item := range testCfg.Fans {
		reqs = append(reqs, stringToIpmiSensorName(item))
	}

	resp, err := dev.ExecuteCommandLine("ipmitool sdr")
	if err != nil {
		log.Printf("ipmitool sdr failed: %s", string(resp))

		return false, fmt.Errorf("ipmitool sdr failed: %w", err), nil
	}

	return sensorIterationHelper(resp, reqs)
}

var errResImplausible = errors.New("result implausible")

func testIPMIPowerStatus(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "power", "status")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if string(resp) != "Chassis Power is on\n" {
		return false, fmt.Errorf("%w: testIPMIPowerStatus: %s", errResImplausible, string(resp)), nil
	}

	return true, nil, nil
}

var errRegesMismatch = errors.New("regex did not match")

func testIPMIChassisUptime(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "chassis", "poh")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	matches := pattern.Match(resp)

	if !matches {
		return false, fmt.Errorf("%w: pattern %s did not match string '%s'", errRegesMismatch, pattern, string(resp)), nil
	}

	return true, nil, nil
}

var errMissingDevice = errors.New("missing device")

func testIPMIFru(dev *testdevice.Device, testCfg *configuration.Config) (
	bool, error, error,
) {
	var testErr error

	resp, err := dev.ExecuteCommandLine("ipmitool", "fru", "|| true")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if strings.Contains(string(resp), "Device not present") {
		return false, fmt.Errorf("%w: %s", errMissingDevice, string(resp)), nil
	}

	for _, fru := range testCfg.Fru {
		foundAll := true

		for _, str := range []string{fru.ProductManufacturer, fru.ProductName, fru.ProductSerial} {
			if str == "" {
				continue
			}

			regex, err := regexp.Compile(str)
			if err != nil {
				return false, fmt.Errorf("%w: regex %s did not compile", err, regex), nil
			}

			if !regex.Match(resp) {
				foundAll = false
				testErr = fmt.Errorf("%w: FRU field: '%s' not present in fru contents: '%s'", testErr, str, string(resp))
			}
		}

		log.Printf("ipmitool FRU present:%s %s %s\n", fru.ProductManufacturer, fru.ProductName, fru.ProductSerial)

		if !foundAll {
			return false, testErr, nil
		}
	}

	return true, nil, nil
}

var errOutputMismatch = errors.New("output does not match")

func testIPMIDCMIPowerReading(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "dcmi", "power", "reading", "|| true")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "dcmi") && !strings.Contains(string(resp), "power") &&
		!strings.Contains(string(resp), "reading") {
		return false, fmt.Errorf("dcmi power reading %w:  %s", errOutputMismatch, string(resp)), nil
	}

	return true, nil, nil
}

var (
	errInitialCountdown      = errors.New("initial countdown insufficient")
	errMessageMissingContent = errors.New("message missing expected content")
)

func testIPMIWatchdogConfiguration(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "mc", "watchdog", "get")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	lines := strings.Split(string(resp), "\n")

	expected := []string{
		"Watchdog Timer Use:     SMS/OS",
		"Watchdog Timer Is:      Started/Running",
		"Watchdog Timer Action:  Power Cycle",
	}

	for _, exp := range expected {
		if !strings.Contains(string(resp), exp) {
			return false, fmt.Errorf("%w:\nmessage: %s expected: %s", errMessageMissingContent, string(resp), exp), nil
		}
	}

	for _, line := range lines {
		if !strings.Contains(line, "Initial Countdown") {
			continue
		}

		re := regexp.MustCompile(`\d+`)

		match := re.FindString(line)

		initialCountdown, err := strconv.Atoi(match)
		if err != nil {
			return false, nil, fmt.Errorf("string conversion failed: %w", err)
		}

		const countdownLimit = 900

		if initialCountdown >= countdownLimit {
			return true, nil, nil
		}

		return false, fmt.Errorf("%w: %d sec", errInitialCountdown, initialCountdown), nil
	}

	return true, nil, nil
}

var (
	errNoGUIDreturned = errors.New("no GUID returned")
	errGUIDMismatch   = errors.New("GUID mismatch")
)

func testIPMISystemGUID(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "mc", "guid", "||  true")
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	log.Println(string(resp))

	if strings.Contains(string(resp), "IPMI command failed: Unspecified error") {
		return false, errNoGUIDreturned, nil
	}

	for _, system := range cfg.Systems {
		if !strings.Contains(string(resp), system.GUID) {
			return false, fmt.Errorf("%w: expected %q, got %q", errGUIDMismatch, system.GUID,
				strings.TrimSpace(string(resp))), nil
		}
	}

	log.Printf("System GUIDs are valid\n")

	return true, nil, nil
}

package ipmi

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/framework"
	ipmimct "github.com/9elements/bmc-test-go/pkg/ipmi_mct"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

// GetIPMIMCTTests returns all IPMI Twitter OEM Extension tests.
//
//nolint:funlen
func GetIPMIMCTTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "MCP PWM Duty Set",
			ShortName: "ipmi-mct-00",
			Status:    framework.StatusImplemented,
			Function:  testMCTPWMDutySet,
		},

		{
			Name:      "MCP PWM Duty Get",
			ShortName: "ipmi-mct-01",
			Status:    framework.StatusImplemented,
			Function:  testMCTPWMDutyGet,
		},

		{
			Name:      "MCP Manufacture Mode Get",
			ShortName: "ipmi-mct-02",
			Status:    framework.StatusImplemented,
			Function:  testMCTManufacturModeGet,
		},

		{
			Name:      "MCP Manufacture Mode Set",
			ShortName: "ipmi-mct-03",
			Status:    framework.StatusImplemented,
			Function:  testMCTManufacturModeSet,
		},

		{
			Name:      "MCP Floor Duty Get",
			ShortName: "ipmi-mct-04",
			Status:    framework.StatusImplemented,
			Function:  testMCTFloorDutyGet,
		},

		{
			Name:      "MCP Floor Duty Set",
			ShortName: "ipmi-mct-05",
			Status:    framework.StatusImplemented,
			Function:  testMCTFloorDutySet,
		},

		{
			Name:      "MCP Get Fru Field",
			ShortName: "ipmi-mct-06",
			Status:    framework.StatusImplemented,
			Function:  testMCTGetFruField,
		},

		{
			Name:      "MCP Set Fru Field",
			ShortName: "ipmi-mct-07",
			Status:    framework.StatusImplemented,
			Function:  testMCTSetFruField,
		},

		{
			Name:      "MCP Get Firmware String",
			ShortName: "ipmi-mct-08",
			Status:    framework.StatusImplemented,
			Function:  testGetFirmwareString,
		},

		{
			Name:      "MCP Config ECC Leaky bucket set",
			ShortName: "ipmi-mct-09",
			Status:    framework.StatusImplemented,
			Function:  testECCLeakyBucketSet,
		},

		{
			Name:      "MCP Config ECC Leaky bucket get",
			ShortName: "ipmi-mct-10",
			Status:    framework.StatusImplemented,
			Function:  testECCLeakyBucketGet,
		},

		{
			Name:      "MCP GPIO Status",
			ShortName: "ipmi-mct-11",
			Status:    framework.StatusImplemented,
			Function:  testGPIOStatus,
		},
	}
}

var errNoTyanIdent = errors.New("return value does not contain Tyan identification (0xfd 0x19 0x00)")

func testMCTPWMDutySet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	setDutyCycle := []uint8{0x00, 0x32}

	ret, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTFanDuty(setDutyCycle))
	if err != nil {
		log.Print(err)
		log.Print(string(ret))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(ret), "fd 19 00") {
		return false, errNoTyanIdent, nil
	}

	setDutyCycle = []uint8{0x00, 0x64}

	_, err = dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTFanDuty(setDutyCycle))
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	return true, nil, nil
}

var errPWMDutyValue = errors.New("return value does not contain value 64, but it should")

func testMCTPWMDutyGet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	getDutyCycle := []uint8{0x00, 0xFE}

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTFanDuty(getDutyCycle))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00") {
		return false, errNoTyanIdent, nil
	}

	if !strings.Contains(string(resp), "64") {
		return false, errPWMDutyValue, nil
	}

	return true, nil, nil
}

func testMCTManufacturModeGet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	value := []uint8{0xFF} // get fan ctrl status

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTManufactureMode(value))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00") {
		return false, errNoTyanIdent, nil
	}

	return true, nil, nil
}

//nolint:dupl
func testMCTManufacturModeSet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	value := []uint8{0x00} // disable fan ctrl

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTManufactureMode(value))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00\n") {
		return false, errNoTyanIdent, nil
	}

	value = []uint8{0xFF} // get fan ctrl status

	resp, err = dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTManufactureMode(value))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00 00\n") {
		return false, errNoTyanIdent, nil
	}

	return true, nil, nil
}

func testMCTFloorDutyGet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	value := []uint8{0xff} // get value

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTFloorDuty(value))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00") {
		return false, errNoTyanIdent, nil
	}

	return true, nil, nil
}

//nolint:dupl
func testMCTFloorDutySet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	value := []uint8{0x2A} // set value to 0x2a

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTFloorDuty(value))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00\n") {
		return false, errNoTyanIdent, nil
	}

	value = []uint8{0xFF} // get value

	resp, err = dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTFloorDuty(value))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), "fd 19 00 2a\n") {
		return false, errNoTyanIdent, nil
	}

	return true, nil, nil
}

var errManufacturerUnexpected = errors.New("unexpected Product Manufacturer")

func testMCTGetFruField(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	var area uint8 = 0x03

	const shift uint8 = 4

	exp := cfg.Fru["baseboard"]
	req := []uint8{0x00, (area << shift)} // fruId, ((area << 4) | field)

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTGetFruField(req))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	var convProManu strings.Builder
	for _, b := range []byte(exp.ProductManufacturer) {
		fmt.Fprintf(&convProManu, "%x ", b)
	}

	// We cut the last character in the converted string,
	// because it is a blank and the response string has a newline (\n)
	convProManuLen := convProManu.Len() - 1

	if !strings.Contains(string(resp[10:]), convProManu.String()[:convProManuLen]) {
		return false, fmt.Errorf("%w: %s, but exp: %s", errManufacturerUnexpected,
			string(resp[10:]), convProManu.String()[:convProManuLen]), nil
	}

	return true, nil, nil
}

func testMCTSetFruField(_ *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	return false, framework.ErrNotImplemented, nil
}

var errUnexpectedFirmwareString = errors.New("unexpected firmware string")

func testGetFirmwareString(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTGetFirmwareString())
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	var convFirmwareString strings.Builder

	for _, b := range []byte(cfg.Name) {
		fmt.Fprintf(&convFirmwareString, "%x ", b)
	}

	length := convFirmwareString.Len()

	if !strings.Contains(string(resp[10:]), convFirmwareString.String()[:length]) {
		return false, fmt.Errorf("%w: %s, want: %s", errUnexpectedFirmwareString,
			string(resp[10:]), convFirmwareString.String()[:length]), nil
	}

	return true, nil, nil
}

func testECCLeakyBucketSet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	values := []uint8{0x08, 0x03} // set T1=0x08, T2=0x03

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTConfigECCLeakyBucket(values))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), " fd 19 00\n") {
		return false, errNoTyanIdent, nil
	}

	return true, nil, nil
}

func testECCLeakyBucketGet(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	values := []uint8{0x08, 0x03} // set T1=0x08, T2=0x03

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTConfigECCLeakyBucket(values))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), " fd 19 00\n") {
		return false, errNoTyanIdent, nil
	}

	values = []uint8{}

	resp, err = dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.MCTConfigECCLeakyBucket(values))
	if err != nil {
		log.Print(err)
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), " fd 19 00 08 03\n") {
		return false, fmt.Errorf("%w and expected values (08 03)", errNoTyanIdent), nil
	}

	return true, nil, nil
}

func testGPIOStatus(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	gpioIndex := []uint8{0x00}

	resp, err := dev.ExecuteCommandLine("ipmitool raw", ipmimct.MCTGPIOStatus(gpioIndex))
	if err != nil {
		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), " fd 19 00 01 00\n") || !strings.Contains(string(resp), " fd 19 00 01 01") {
		return false, fmt.Errorf("%w: return value does not contain expected values", errNoTyanIdent), nil
	}

	return true, nil, nil
}

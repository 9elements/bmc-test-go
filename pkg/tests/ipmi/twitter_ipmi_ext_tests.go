package ipmi

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/framework"
	ipmimct "github.com/9elements/bmc-test-go/pkg/ipmi_mct"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

// GetIPMITwitterTests returns all IPMI Twitter OEM Extension tests.
func GetIPMITwitterTests() []*framework.Test {
	return []*framework.Test{
		{
			Name:      "Twitter IPMI Ext Set Service",
			ShortName: "ipmi-twitter-00",
			Function:  testTwitterIPMIExtSetService,
		},
		{
			Name:      "Twitter IPMI Ext Get Service",
			ShortName: "ipmi-twitter-01",
			Function:  testTwitterIPMIExtGetService,
		},
		{
			Name:      "Twitter IPMI Ext Clear CMOS",
			ShortName: "ipmi-twitter-02",
			Function:  testTwitterIPMIExtClearCMOS,
		},

		{
			Name:      "Twitter IPMI Ext Pnm Get Reading",
			ShortName: "ipmi-twitter-03",
			Function:  testTwitterIPMIExtPnmGetReading,
		},

		{
			Name:      "Twitter IPMI Ext Random Delay AC Restore Power On",
			ShortName: "ipmi-twitter-04",
			Function:  testTwitterIPMIExtRandomDelayACRestorePowerOn,
		},

		{
			Name:      "Twitter IPMI Ext Get Post Codes",
			ShortName: "ipmi-twitter-05",
			Function:  testTwitterIPMIExtGetPostCodes,
		},
	}
}

var (
	errInvalidResponse      = errors.New("invalid response")
	errResultMissingContent = errors.New("result does not contain expected value")
)

func testTwitterIPMIExtSetService(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	ret, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.TwitterSetService([]byte{0x01}))
	if err != nil {
		log.Print(string(ret))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if string(ret) != "\n" {
		return false, fmt.Errorf("%w: %s", errInvalidResponse, string(ret)), nil
	}

	return true, nil, nil
}

func testTwitterIPMIExtGetService(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	ret, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.TwitterGetService([]byte{}))
	if err != nil {
		log.Print(string(ret))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if string(ret) != " 01\n" && string(ret) != " 00\n" {
		return false, fmt.Errorf("%w: %s", errInvalidResponse, string(ret)), nil
	}

	return true, nil, nil
}

func testTwitterIPMIExtClearCMOS(_ *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	// TODO: implement
	return false, fmt.Errorf("not implemented"), nil
}

func testTwitterIPMIExtPnmGetReading(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	vals := []uint8{0x0, 0x0, 0x0}

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.TwitterPnmGetReading(vals))
	if err != nil {
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), " 00 5b 00\n") {
		return false, fmt.Errorf("%w: (00 5b 00)", errResultMissingContent), nil
	}

	return true, nil, nil
}

func testTwitterIPMIExtRandomDelayACRestorePowerOn(dev *testdevice.Device, _ *configuration.Config,
) (
	bool, error, error,
) {
	const delayLSB uint8 = 0x30

	const delayMSB uint8 = 0x54

	delay := []uint8{0x82, delayLSB, delayMSB} // no idea what the 0x82 is

	resp, err := dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.TwitterRandomDelayACRestorePowerOn(delay))
	if err != nil {
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if !strings.Contains(string(resp), " 82 30 54\n") {
		return false, fmt.Errorf("%w: (82 30 54)", errResultMissingContent), nil
	}

	return true, nil, nil
}

var (
	errNoPostCodes      = errors.New("no post codes returned")
	errTooManyPostCodes = errors.New("too many post codes")
)

func testTwitterIPMIExtGetPostCodes(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	resp, err := dev.ExecuteCommandLine("ipmitool", "power", "on")
	if err != nil {
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	const timeSleep time.Duration = 40 * time.Second

	time.Sleep(time.Duration(timeSleep.Seconds())) // Wait 40 seconds for post codes to show up

	resp, err = dev.ExecuteCommandLine("ipmitool", "raw", ipmimct.TwitterGetPostCodes())
	if err != nil {
		log.Print(string(resp))

		return false, nil, fmt.Errorf("command execution returned error: %w", err)
	}

	if len(resp) == 0 {
		return false, errNoPostCodes, nil
	}

	const limitCode int = 20

	codes := strings.Split(string(resp), " ")
	if len(codes) > limitCode {
		return false, fmt.Errorf("%w: %d", errTooManyPostCodes, len(codes)), nil
	}

	return true, nil, nil
}

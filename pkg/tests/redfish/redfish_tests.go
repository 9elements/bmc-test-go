package redfish

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
	"github.com/stmcginnis/gofish/schemas"
)

var (
	errNoChassis = errors.New("no chassis returned")
	errNoSystem  = errors.New("no system returned")
)

func testRedfishChassis(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	chassis, err := dev.RedfishService().Chassis()
	if err != nil {
		return false, nil, fmt.Errorf("redfish chassis query failed: %w", err)
	}

	if len(chassis) < 1 {
		return false, errNoChassis, nil
	}

	return true, nil, nil
}

func testRedfishSystem(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	systems, err := dev.RedfishService().Systems()
	if err != nil {
		return false, nil, fmt.Errorf("redfish system query failed: %w", err)
	}

	if len(systems) < 1 {
		return false, errNoSystem, nil
	}

	return true, nil, nil
}

var (
	errFWVersionFormat = errors.New("fw version wrong format")
	errCommitIDEmpty   = errors.New("commit id empty")
)

func testRedfishGetFirmwareVersion(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	updateService, err := dev.RedfishService().UpdateService()
	if err != nil {
		return false, nil, fmt.Errorf("redfish update service query failed: %w", err)
	}

	fwInventories, err := updateService.FirmwareInventory()
	if err != nil {
		return false, nil, fmt.Errorf("update service firmware inventories query failed: %w", err)
	}

	fw := fwInventories[0]

	// e.g. v59.99-0-g1cb1416454-s8030nf
	parts := strings.Split(fw.Version, "-")

	const lowerLimit = 3

	if len(parts) < lowerLimit {
		return false, fmt.Errorf("%w: '%s' ", errFWVersionFormat, fw.Version), nil
	}

	commitID := parts[2][1:]

	if commitID == "" {
		return false, errCommitIDEmpty, nil
	}

	return true, nil, nil
}

var errNoPowerInformation = errors.New("no power information for chassis available")

func getPower(dev *testdevice.Device) ([]*schemas.Power, error) {
	chassis, err := dev.RedfishService().Chassis()
	if err != nil {
		return nil, fmt.Errorf("redfish chassis query failed: %w", err)
	}

	if len(chassis) < 1 {
		return nil, errNoChassis
	}

	power := make([]*schemas.Power, 0)

	for _, c := range chassis {
		p, err := c.Power()
		if err != nil {
			return nil, fmt.Errorf("call to c.Power failed: %w", err)
		}

		power = append(power, p)
	}

	if len(power) < 1 {
		return nil, errNoPowerInformation
	}

	return power, nil
}

func getPowerSupplies(dev *testdevice.Device) ([]*schemas.PowerSupply, error) {
	chassis, err := dev.RedfishService().Chassis()
	if err != nil {
		return nil, fmt.Errorf("redfish chassis query failed: %w", err)
	}

	if len(chassis) < 1 {
		return nil, errNoChassis
	}

	power := make([]*schemas.PowerSupply, 0)

	for _, c := range chassis {
		pss, err := c.PowerSubsystem()
		if err != nil {
			return nil, fmt.Errorf("call to c.PowerSubsystem failed: %w", err)
		}

		if pss == nil {
			continue
		}

		psu, err := pss.PowerSupplies()
		if err != nil {
			return nil, fmt.Errorf("failed to pss.PowerSupplies: %w", err)
		}

		power = append(power, psu...)
	}

	if len(power) < 1 {
		return nil, errNoPowerInformation
	}

	return power, nil
}

var (
	errItemCountMismatch = errors.New("item count mismatch")
	errNameMismatch      = errors.New("value of item mismatch")
)

//nolint:cyclop
func testRedfishCheckPowerSupplyInformation(dev *testdevice.Device, cfg *configuration.Config,
) (
	bool, error, error,
) {
	psus, err := getPowerSupplies(dev)
	if err != nil {
		return false, nil, err
	}

	for _, psu := range psus {
		if len(psus) != cfg.PSUCount {
			return false, fmt.Errorf("%w: psu. Have %d, want: %d", errItemCountMismatch, len(psus),
				cfg.PSUCount), nil
		}

		for _, expPSU := range cfg.PSUs {
			if psu.Name != expPSU.Name {
				continue
			}

			if psu.Name != expPSU.Name {
				return false, fmt.Errorf("%w: power supplied name: %s does not match config name: %s",
					errNameMismatch, psu.Name, expPSU.Name), nil
			}

			if psu.Model != expPSU.Model {
				return false, fmt.Errorf("%w: power supplied model: %s does not match config model: %s",
					errNameMismatch, psu.Model, expPSU.Model), nil
			}

			if psu.PartNumber != expPSU.Partnumber {
				return false, fmt.Errorf("%w: power supplied partnumber: %s does not match config partnumber: %s",
					errNameMismatch, psu.PartNumber, expPSU.Partnumber), nil
			}

			if psu.Manufacturer != expPSU.Manufacturer {
				return false, fmt.Errorf("%w: power supplied manufacturer: %s does not match config manufacturer: %s",
					errNameMismatch, psu.Manufacturer, expPSU.Manufacturer), nil
			}

			if psu.SerialNumber != expPSU.SerialNumber {
				return false, fmt.Errorf("%w: power supplied serial number: %s does not match config serial number: %s",
					errNameMismatch, psu.SerialNumber, expPSU.SerialNumber), nil
			}

			if *psu.EfficiencyPercent != expPSU.EfficiencyPercent {
				return false, fmt.Errorf("%w: power supplied efficiency percent: %f does not match config efficiency percent: %f",
					errNameMismatch, *psu.EfficiencyPercent, expPSU.EfficiencyPercent), nil
			}
		}
	}

	return true, nil, nil
}

var (
	errNumPowerSupplyMismatch = errors.New("power supply count mismatch")
	errSensorNotFound         = errors.New("sensor not found")
)

func testRedfishVoltagesSensorNames(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	power, err := getPower(dev)
	if err != nil {
		return false, nil, fmt.Errorf("getPower failed: %w", err)
	}

	if len(power) != cfg.PSUCount {
		return false, fmt.Errorf("%w: Have: %d, want: %d", errNumPowerSupplyMismatch, len(power), cfg.PSUCount), nil
	}

	voltage := make([]schemas.Voltage, 0)

	for _, p := range power {
		if p != nil {
			voltage = append(voltage, p.Voltages...)
		}
	}

	for _, cfgVoltage := range cfg.Voltage {
		found := false

		for _, rVolts := range voltage {
			if cfgVoltage == rVolts.Name {
				found = true
			}
		}

		if !found {
			return false, fmt.Errorf("%w: voltage - with name: %s not found", errSensorNotFound, cfgVoltage), nil
		}
	}

	return true, nil, nil
}

var (
	errNoInformation  = errors.New("no information available")
	errSensorNotValid = errors.New("sensor reading not reasonable")
)

func testRedfishVoltageSensorValueThesholds(dev *testdevice.Device, _ *configuration.Config) (bool, error, error) {
	power, err := getPower(dev)
	if err != nil {
		return false, nil, err
	}

	voltage := make([]schemas.Voltage, 0)

	for _, p := range power {
		if p != nil {
			voltage = append(voltage, p.Voltages...)
		}
	}

	if len(voltage) < 1 {
		return false, fmt.Errorf("%w: voltage", errNoInformation), nil
	}

	for _, rVolts := range voltage {
		lc := *rVolts.LowerThresholdCritical
		uc := *rVolts.UpperThresholdCritical
		readingVolts := *rVolts.ReadingVolts

		gap := uc - lc

		const f float32 = 1.2

		believable := (lc-gap*f) <= readingVolts && readingVolts <= (uc+gap*f)

		if !believable {
			return false, fmt.Errorf("%w: %s ", errSensorNotValid, rVolts.Name), nil
		}
	}

	return true, nil, nil
}

func getTemperatures(dev *testdevice.Device) ([]schemas.Temperature, error) {
	chassis, err := dev.RedfishService().Chassis()
	if err != nil {
		return nil, fmt.Errorf("redfish chassis query failed: %w", err)
	}

	if len(chassis) < 1 {
		return nil, errNoChassis
	}

	var temps []schemas.Temperature

	for _, c := range chassis {
		thermal, err := c.Thermal()
		if err != nil {
			return nil, fmt.Errorf("call to c.Thermal() failed: %w", err)
		}

		if thermal != nil {
			temps = append(temps, thermal.Temperatures...)
		}
	}

	return temps, nil
}

var errNoTempsReturned = errors.New("no temperatures returned")

func testRedfishTemperatureSensorNames(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	temps, err := getTemperatures(dev)
	if err != nil {
		return false, nil, fmt.Errorf("getTemperatures failed: %w", err)
	}

	if len(temps) < 1 {
		return false, errNoTempsReturned, nil
	}

	for _, exp := range cfg.Temperatures {
		found := false

		for _, t := range temps {
			if t.Name == exp {
				found = true
			}
		}

		if !found {
			return false, fmt.Errorf("%w: temperatur - with name: %s not found", errSensorNotFound, exp), nil
		}
	}

	return true, nil, nil
}

func testRedfishTemperaturSensorValueThesholds(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	temps, err := getTemperatures(dev)
	if err != nil {
		return false, nil, err
	}

	if len(temps) < 1 {
		return false, errNoTempsReturned, nil
	}

	for _, t := range temps {
		lc := *t.LowerThresholdCritical
		uc := *t.UpperThresholdCritical
		readingTemps := *t.ReadingCelsius

		gap := uc - lc

		const f float64 = 1.2

		believable := (lc-gap*f) <= readingTemps && readingTemps <= (uc+gap*f)

		if !believable {
			return false, fmt.Errorf("%w: temperatur sensor %s", errSensorNotValid, t.Name), nil
		}
	}

	return true, nil, nil
}

func getFans(dev *testdevice.Device) ([]schemas.ThermalFan, error) {
	chassis, err := dev.RedfishService().Chassis()
	if err != nil {
		return nil, fmt.Errorf("redfish chassis query failed: %w", err)
	}

	if len(chassis) < 1 {
		return nil, errNoChassis
	}

	ret := make([]schemas.ThermalFan, 0)

	for _, c := range chassis {
		thermals, err := c.Thermal()
		if err != nil {
			return nil, fmt.Errorf("call to c.Thermal() failed: %w", err)
		}

		if thermals != nil {
			ret = append(ret, thermals.Fans...)
		}
	}

	return ret, nil
}

func testRedfishThermalFanNames(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	fans, err := getFans(dev)
	if err != nil {
		return false, nil, err
	}

	for _, exp := range cfg.Fans {
		found := false

		for _, f := range fans {
			if exp == f.Name {
				found = true
			}
		}

		if !found {
			return false, fmt.Errorf("%w: fan - with name: %s not found", errSensorNotFound, exp), nil
		}
	}

	return true, nil, nil
}

func testRedfishThermalFanValues(dev *testdevice.Device, _ *configuration.Config) (
	bool, error, error,
) {
	fans, err := getFans(dev)
	if err != nil {
		return false, nil, err
	}

	for _, fan := range fans {
		if fan.ReadingUnits == "RPM" {
			lc := float64(*fan.LowerThresholdCritical)
			uc := float64(*fan.UpperThresholdCritical)
			reading := float64(*fan.Reading)

			gap := uc - lc

			const f float64 = 1.2

			believable := (lc-gap*f) <= reading && reading <= (uc+gap*f)

			if !believable {
				return false, fmt.Errorf("%w: fan sensor %s ", errSensorNotValid, fan.Name), nil
			}
		}
	}

	return true, nil, nil
}

var (
	errDimmCountMismatch    = errors.New("dimm count does not match")
	errDimmPropertyMismatch = errors.New("dimm property mismatch")
)

//nolint:cyclop
func testRedfishMemory(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	sys, err := dev.RedfishService().Systems()
	if err != nil {
		return false, nil, fmt.Errorf("call to RedfishService().Systems() failed: %w", err)
	}

	if len(sys) < 1 {
		return false, errNoSystem, nil
	}

	mem := make([]*schemas.Memory, 0)

	for _, s := range sys {
		m, err := s.Memory()
		if err != nil {
			return false, nil, fmt.Errorf("call to s.Memory failed: %w", err)
		}

		mem = append(mem, m...)
	}

	if len(mem) != cfg.Dimms.Count {
		return false, fmt.Errorf("%w. Have: %d, want: %d", errDimmCountMismatch, len(mem), cfg.Dimms.Count), nil
	}

	for _, m := range mem {
		if m.Manufacturer != cfg.Dimms.Manufacturer {
			return false, fmt.Errorf("%w: dimm manufacturer: Have: %s, want: %s", errDimmPropertyMismatch,
				m.Manufacturer, cfg.Dimms.Manufacturer), nil
		}

		if m.PartNumber != cfg.Dimms.PartNumber {
			return false, fmt.Errorf("%w: dimm part number: have: %s, want: %s", errDimmPropertyMismatch,
				m.PartNumber, cfg.Dimms.PartNumber), nil
		}

		if m.SerialNumber != cfg.Dimms.SerialNumber {
			return false, fmt.Errorf("%w: dimm serial number: have %s, want: %s", errDimmPropertyMismatch,
				m.SerialNumber, cfg.Dimms.SerialNumber), nil
		}

		if *m.CapacityMiB != cfg.Dimms.Capacity {
			return false, fmt.Errorf("%w: dimm capacity: have: %d, want: %d", errDimmPropertyMismatch,
				m.CapacityMiB, cfg.Dimms.Capacity), nil
		}

		if *m.OperatingSpeedMhz != cfg.Dimms.OperatingSpeedMHz {
			return false, fmt.Errorf("%w: dimm operating speed: have: %d, want: %d", errDimmPropertyMismatch,
				m.OperatingSpeedMhz, cfg.Dimms.OperatingSpeedMHz), nil
		}
	}

	return true, nil, nil
}

var (
	errNoProcessors              = errors.New("no processors returned")
	errProcessorCountMismatch    = errors.New("cpu count mismatch")
	errProcessorPropertyMismatch = errors.New("cpu property mismatch")
)

//nolint:cyclop
func testRedfishProcessor(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	sys, err := dev.RedfishService().Systems()
	if err != nil {
		return false, nil, fmt.Errorf("call to RedfishService.Systems() failed: %w", err)
	}

	if len(sys) < 1 {
		return false, errNoSystem, nil
	}

	cpus := make([]*schemas.Processor, 0)

	for _, s := range sys {
		cpu, err := s.Processors()
		if err != nil {
			return false, nil, fmt.Errorf("call to s.Processor() failed: %w", err)
		}

		cpus = append(cpus, cpu...)
	}

	if len(cpus) != len(cfg.Cpus) {
		return false, errNoProcessors, nil
	}

	if cfg.CPUCount != len(cpus) {
		return false, fmt.Errorf("%w: have: %d, want: %d", errProcessorCountMismatch, len(cpus), cfg.CPUCount), nil
	}

	for n, c := range cpus {
		if c.Manufacturer != cfg.Cpus[n].Manufacturer {
			return false, fmt.Errorf("%w: cpu manufacturer. Have: %s, want: %s", errProcessorPropertyMismatch,
				c.Manufacturer, cfg.Cpus[n].Manufacturer), nil
		}

		if string(c.InstructionSet) != cfg.Cpus[n].InstructionSet {
			return false, fmt.Errorf("%w: cpu instruction set. Have: %s, want: %s", errProcessorPropertyMismatch,
				string(c.InstructionSet), cfg.Cpus[n].InstructionSet), nil
		}

		if string(c.ProcessorType) != cfg.Cpus[n].Type {
			return false, fmt.Errorf("%w: cpu type. Have: %s, want: %s", errProcessorPropertyMismatch,
				string(c.ProcessorType), cfg.Cpus[n].Type), nil
		}

		if string(c.ProcessorArchitecture) != cfg.Cpus[n].InstructionSet {
			return false, fmt.Errorf("%w: cpu architecture. Have: %s, want: %s", errProcessorPropertyMismatch,
				string(c.ProcessorArchitecture), cfg.Cpus[n].InstructionSet), nil
		}

		if c.Model != cfg.Cpus[n].Model {
			return false, fmt.Errorf("%w: cpu model. Have: %s, want: %s", errProcessorPropertyMismatch, c.Model,
				cfg.Cpus[n].Model), nil
		}
	}

	return true, nil, nil
}

func getFirmwareVersion(dev *testdevice.Device) (string, error) {
	service, err := dev.RedfishService().UpdateService()
	if err != nil {
		return "", fmt.Errorf("call to RedfishService().UpdateService() failed: %w", err)
	}

	inv, err := service.FirmwareInventory()
	if err != nil {
		return "", fmt.Errorf("update service firmware inventories query failed: %w", err)
	}

	fw := inv[0]
	parts := strings.Split(fw.Version, "-")

	log.Printf("Debug: parts: %v\n", parts)

	if len(parts) < 3 {
		return "", fmt.Errorf("fw version did not split as expected: %s", fw.Version)
	}

	return parts[2][1:], nil
}

func makeTransport() *http.Transport {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		DialContext: (&net.Dialer{
			Timeout: time.Second, // Timeout for establishing a new connection
		}).DialContext,
	}

	return transport
}

var errHTTPResponseFailure = errors.New("response indicates failure")

func postRequestData(cfg *configuration.Config, imgpath string, timeout int64) (string, error) {
	image, err := os.ReadFile(path.Clean(imgpath))
	if err != nil {
		return "", fmt.Errorf("reading file failed: %w", err)
	}

	url := fmt.Sprintf("https://%s/redfish/v1/UpdateService/update", cfg.BMCHost)

	var buf bytes.Buffer

	_, err = buf.Read(image)
	if err != nil {
		return "", fmt.Errorf("failed to put image into buffer: %w", err)
	}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url,
		&buf)
	if err != nil {
		return "", fmt.Errorf("request generation failed: %w", err)
	}

	request.SetBasicAuth(cfg.BMCUser, cfg.BMCPassword)
	request.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{
		Transport: makeTransport(),
		Timeout:   time.Second * time.Duration(timeout),
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("client.Do failed: %w", err)
	}
	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf(":%w: url=%s, status=%s body=%s", errHTTPResponseFailure, url, response.Status, string(body))
	}

	var result struct {
		OdataID string `json:"@odata.id"`
	}

	err = json.Unmarshal(body, &result)
	if err == nil && result.OdataID != "" {
		return result.OdataID, nil
	}

	return "", nil
}

func resetBMC(cfg *configuration.Config) error {
	url := fmt.Sprintf("https://%s/redfish/v1/Managers/bmc/Actions/Manager.Reset", cfg.BMCHost)
	body := bytes.NewBufferString(`{"ResetType":"GracefulRestart"}`)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("failed to make request with context: %w", err)
	}

	req.SetBasicAuth(cfg.BMCUser, cfg.BMCPassword)
	req.Header.Set("Content-Type", "application/json")

	const defaultTimeoutVal = 30

	client := &http.Client{
		Transport: makeTransport(),
		Timeout:   time.Second * defaultTimeoutVal,
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("BMC reset failed:%w : status=%s body=%s", errHTTPResponseFailure, resp.Status, string(b))
	}

	return nil
}

var errFlashingFailed = errors.New("flashing failed")

// waitForTask polls the task until it reaches a terminal state.
// Returns (autoRebooted=true) if the BMC rebooted itself during the flash,
// meaning no manual Manager.Reset is needed.
//
//nolint:cyclop,funlen
func waitForTask(cfg *configuration.Config, taskURI string, timeout time.Duration) (bool, error) {
	const defaultTimeout = 10

	client := &http.Client{
		Transport: makeTransport(),
		Timeout:   time.Second * defaultTimeout,
	}

	url := fmt.Sprintf("https://%s%s", cfg.BMCHost, taskURI)

	deadline := time.Now().Add(timeout)
	wasUnreachable := false

	for time.Now().Before(deadline) {
		time.Sleep(defaultTimeout * time.Second)

		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
		if err != nil {
			return false, fmt.Errorf("creating request failed: %w", err)
		}

		req.SetBasicAuth(cfg.BMCUser, cfg.BMCPassword)

		resp, err := client.Do(req)
		if err != nil {
			wasUnreachable = true

			continue
		}

		rawBody, _ := io.ReadAll(resp.Body)

		err = resp.Body.Close()
		if err != nil {
			return false, fmt.Errorf("call to resp.Body.Close() failed: %w", err)
		}

		var task struct {
			TaskState  string `json:"TaskState"`
			TaskStatus string `json:"TaskStatus"`
			Messages   []struct {
				Message string `json:"Message"`
			} `json:"Messages"`
		}

		err = json.Unmarshal(rawBody, &task)
		if err != nil {
			continue
		}

		// After a reboot the task table is cleared — empty state means BMC rebooted during flash
		if task.TaskState == "" && wasUnreachable {
			return true, nil
		}

		switch task.TaskState {
		case "Completed":
			return false, nil

		case "Killed", "Exception", "Cancelled":
			msgs := make([]string, len(task.Messages))

			for i, m := range task.Messages {
				msgs[i] = m.Message
			}

			return false, fmt.Errorf("%w flash task failed (state=%s status=%s): %s", errFlashingFailed,
				task.TaskState, task.TaskStatus, strings.Join(msgs, "; "))
		}
	}

	return false, fmt.Errorf("%w: task did not complete within %s minutes", errFlashingFailed, timeout)
}

var errBMCrebootFailed = errors.New("bmc reboot failed")

func waitForUpdateReset(cfg *configuration.Config, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	const defaultTimeout = 10

	const defaultTimeoutHalf = 5

	client := http.Client{
		Transport: makeTransport(),
		Timeout:   time.Second * defaultTimeout,
	}

	url := fmt.Sprintf("https://%s/redfish/v1", cfg.BMCHost)

	poll := func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return fmt.Errorf("creating new request failed: %w", err)
		}

		req.SetBasicAuth(cfg.BMCUser, cfg.BMCPassword)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("call to client.Do() failed: %w", err)
		}

		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("call to resp.Body.Close() failed: %w", err)
		}

		return nil
	}

	// wait for BMC to go down
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: BMC did not go down within %s minutes", errBMCrebootFailed, timeout)

		case <-time.After(defaultTimeoutHalf * time.Second):
		}

		err := poll()
		if err != nil {
			break
		}
	}

	// wait for BMC to come back up
	return waitForBMCUp(ctx, poll)
}

func waitForBMCUpSimple(cfg *configuration.Config, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	const defaultTimeout = 10

	bmcClient := http.Client{Transport: makeTransport(), Timeout: time.Second * defaultTimeout}
	bmcURL := fmt.Sprintf("https://%s/redfish/v1", cfg.BMCHost)

	return waitForBMCUp(ctx, func() error {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, bmcURL, nil)
		if err != nil {
			return fmt.Errorf("request generation failed: %w", err)
		}

		req.SetBasicAuth(cfg.BMCUser, cfg.BMCPassword)

		resp, err := bmcClient.Do(req)
		if err != nil {
			return fmt.Errorf("bmcClient.Do failed: %w", err)
		}

		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("call to resp.Body.Close() failed: %w", err)
		}

		return nil
	})
}

func waitForBMCUp(ctx context.Context, poll func() error) error {
	const defaultTimeoutHalf = 5

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w: BMC did not come back within the timeout", errBMCrebootFailed)

		case <-time.After(defaultTimeoutHalf * time.Second):
		}

		err := poll()
		if err != nil {
			continue
		}

		return nil
	}
}

// flashAndWait uploads an image, waits for the flash to complete and the BMC
// to come back up, then reconnects the Redfish client.
func flashAndWait(dev *testdevice.Device, cfg *configuration.Config, imgpath string) error {
	const dataRequestTimeout = 120

	const waitForBMCUp = 10 * time.Minute

	const waitForTaskDelay = 5 * time.Minute

	taskURI, err := postRequestData(cfg, imgpath, dataRequestTimeout)
	if err != nil {
		return err
	}

	if taskURI == "" {
		return nil
	}

	autoRebooted, err := waitForTask(cfg, taskURI, waitForTaskDelay)
	if err != nil {
		return err
	}

	if autoRebooted {
		err := waitForBMCUpSimple(cfg, waitForBMCUp)
		if err != nil {
			return err
		}
	} else {
		err := resetBMC(cfg)
		if err != nil {
			return err
		}

		err = waitForUpdateReset(cfg, waitForBMCUp)
		if err != nil {
			return err
		}
	}

	log.Printf("flashAndWait: BMC back online")

	err = dev.Reconnect(cfg)
	if err != nil {
		return fmt.Errorf("reconnect after flash failed: %w", err)
	}

	return nil
}

var (
	errDowngradeFailed = errors.New("downgrade failed")
	errRestoreFailed   = errors.New("restore failed")
)

func testRedfishFirmwareDowngradeUpdate(dev *testdevice.Device, cfg *configuration.Config) (
	bool, error, error,
) {
	commitIDactual, err := getFirmwareVersion(dev)
	if err != nil {
		return false, nil, err
	}

	log.Printf("firmware downgrade: current commit ID: %s", commitIDactual)

	// flash golden (older) image
	log.Printf("firmware downgrade: flashing golden image: %s", cfg.GoldenImage)

	err = flashAndWait(dev, cfg, cfg.GoldenImage)
	if err != nil {
		return false, nil, err
	}

	oldcommitID, err := getFirmwareVersion(dev)
	if err != nil {
		return false, nil, err
	}

	log.Printf("firmware downgrade: commit ID after downgrade: %s", oldcommitID)

	downgradeOK := oldcommitID != commitIDactual

	// always restore CI image, regardless of whether downgrade succeeded
	log.Printf("firmware downgrade: restoring CI image: %s", cfg.Binary)

	err = flashAndWait(dev, cfg, cfg.Binary)
	if err != nil {
		return false, nil, err
	}

	renewCommitID, err := getFirmwareVersion(dev)
	if err != nil {
		return false, nil, err
	}

	log.Printf("firmware downgrade: commit ID after restore: %s", renewCommitID)

	if !downgradeOK {
		return false, fmt.Errorf("%w: downgrade did not change commit ID (before=%s, after=%s)",
			errDowngradeFailed, commitIDactual, oldcommitID), nil
	}

	if renewCommitID != commitIDactual {
		return false, fmt.Errorf("%w: restore did not return to original commit ID (want=%s, got=%s)",
			errRestoreFailed, commitIDactual, renewCommitID), nil
	}

	return true, nil, nil
}

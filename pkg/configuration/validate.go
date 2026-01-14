package configuration

import (
	"errors"
	"fmt"
	"log"
)

var (
	errFieldEmpty          = errors.New("field is empty")
	errNoHostsDefines      = errors.New("hosts field missing or malformed")
	errMissingEntries      = errors.New("missing entries for field")
	errVersionValue        = errors.New("version value must be 0")
	errNameFieldEmpty      = errors.New("name field must not be empty")
	errBMCField            = errors.New("bmc structure invalid")
	errHostField           = errors.New("host structures invalid")
	errTestDataField       = errors.New("test data strcture invalid")
	errMismatchSystemCount = errors.New("system count and amount of systems defined does not match")
	errMismatchCPUCount    = errors.New("cpu count and amount of cpu defined does not match")
	errMismatchPSUCount    = errors.New("psu count and amount of psu defined does not match")
)

func validateCfgBMC(bmc *BMC) error {
	var err error
	if bmc.BMCHost == "" {
		err = fmt.Errorf("%w: bmc field 'host'", errFieldEmpty)
	}

	if bmc.BMCUser == "" {
		err = fmt.Errorf("%w\n %w: 'bmc field 'host'", err, errFieldEmpty)
	}

	if bmc.BMCPassword == "" {
		err = fmt.Errorf("%w\n %w: bmc field 'password'", err, errFieldEmpty)
	}

	if bmc.BMCSSHKey == "" {
		err = fmt.Errorf("%w\n %w: bmc field 'sshKey'", err, errFieldEmpty)
	}

	if bmc.SSHPort == "" {
		err = fmt.Errorf("%w\n %w: bmc field 'sshPort'", err, errFieldEmpty)
	}

	return err
}

func validateCfgHosts(hosts []Host, num int) error {
	if len(hosts) < 1 {
		return errNoHostsDefines
	}

	var errRet error

	numHosts := 0
	for num, host := range hosts {
		numHosts++

		err := validateCfgHost(&host, num)
		if err != nil {
			errRet = fmt.Errorf("%w : %w", errBMCField, err)
		}
	}

	if num != numHosts {
		errRet = fmt.Errorf("%w: number of hosts and defined hosts dont match", errRet)
	}

	return errRet
}

func validateCfgHost(host *Host, num int) error {
	var err error

	if host.Name == "" {
		err = fmt.Errorf("%w\n %w: host%d field 'name'", err, errFieldEmpty, num)
	}

	if host.IP == "" {
		err = fmt.Errorf("%w\n %w: host%d field 'ip'", err, errFieldEmpty, num)
	}

	if host.Username == "" {
		err = fmt.Errorf("%w\n %w: host%d field 'username'", err, errFieldEmpty, num)
	}

	if host.Password == "" {
		err = fmt.Errorf("%w\n %w: host%d field 'password'", err, errFieldEmpty, num)
	}

	if host.SSHKey == "" {
		err = fmt.Errorf("%w\n %w: host%d field 'sshKey'", err, errFieldEmpty, num)
	}

	if host.SSHPort == "" {
		err = fmt.Errorf("%w\n %w: host%d field 'sshPort'", err, errFieldEmpty, num)
	}

	return err
}

func validateCfgTestData(data *TestData) error {
	var err error

	if data.SystemCount < len(data.Systems) {
		err = errMismatchSystemCount
	}

	if data.CPUCount != len(data.Cpus) {
		err = fmt.Errorf("%w\n %w", err, errMismatchCPUCount)
	}

	if data.PSUCount != len(data.PSUs) {
		err = fmt.Errorf("%w\n %w", err, errMismatchPSUCount)
	}

	return err
}

// Validate checks the config file for required fields to be filled.
func Validate(cfg *Config) error {
	var errRes error

	log.Printf("%s\n", "Validating config")

	if cfg.Version != 0 {
		return errVersionValue
	}

	if cfg.Name == "" {
		return errNameFieldEmpty
	}

	err := validateCfgBMC(&cfg.BMC)
	if err != nil {
		errRes = fmt.Errorf("%w: %w", errBMCField, err)
	}

	err = validateCfgHosts(cfg.Hosts, cfg.NumHosts)
	if err != nil {
		errRes = fmt.Errorf("%w : %w", errHostField, err)
	}

	err = validateCfgTestData(&cfg.TestData)
	if err != nil {
		errRes = fmt.Errorf("%w : %w", errTestDataField, err)
	}

	if len(cfg.Voltage) < 1 {
		return fmt.Errorf("testdata->voltage: %w", errMissingEntries)
	}

	if len(cfg.Temperatures) < 1 {
		return fmt.Errorf("testdata->temperature: %w", errMissingEntries)
	}

	if len(cfg.Fans) < 1 {
		return fmt.Errorf("testdata->fans: %w", errMissingEntries)
	}

	return errRes
}

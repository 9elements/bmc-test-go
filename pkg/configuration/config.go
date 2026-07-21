// Package configuration package holds the structures and functions to manage configuration yaml files.
package configuration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

var errNoPath = errors.New("environment variable empty")

func expandPath(p string) string {
	p = os.ExpandEnv(p)
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			p = filepath.Join(home, p[2:])
		}
	}

	return p
}

func load(path string) (*Config, error) {
	var cfg Config

	p := filepath.Clean(path)

	file, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration file: %w", err)
	}

	return &cfg, nil
}

func validateImagePath(path string, field string) (string, error) {
	var binPath string
	if path != "" {
		binPath = expandPath(path)
		if binPath == "" {
			return "", fmt.Errorf("%w: for field: %s", errNoPath, field)
		}
	}

	return binPath, nil
}

// LoadConfig loads a configuration file in yaml format from the given path string.
func LoadConfig(path string) (*Config, error) {
	cfg, err := load(path)
	if err != nil {
		return nil, err
	}

	cfg.Binary, err = validateImagePath(cfg.Binary, "Binary")
	if err != nil {
		return nil, err
	}

	cfg.GoldenImage, err = validateImagePath(cfg.GoldenImage, "GoldenImage")
	if err != nil {
		return nil, err
	}

	cfg.BMCSSHKey, err = validateImagePath(cfg.BMCSSHKey, "BMCSSHKey")
	if err != nil {
		return nil, err
	}

	for num, host := range cfg.Hosts {
		if host.SSHKey != "" {
			hostkey := expandPath(host.SSHKey)

			cfg.Hosts[num].SSHKey, err = validateImagePath(hostkey, "host SSHKey")
			if err != nil {
				return nil, err
			}
		}
	}

	return cfg, nil
}

// TestData aggregates the information required by the tests for validation against the system.
type TestData struct {
	Cpus         []CPU                `yaml:"cpu"`
	Systems      []System             `yaml:"systems"`
	SystemCount  int                  `yaml:"systemCount"`
	CPUCount     int                  `yaml:"cpuCount"`
	PSUCount     int                  `yaml:"psuCount"`
	PSUs         []PSU                `yaml:"psu"`
	Voltage      []string             `yaml:"voltage"`
	Temperatures []string             `yaml:"temperatures"`
	Fans         []string             `yaml:"fans"`
	Dimms        DimmInfo             `yaml:"dimms"`
	Fru          map[string]FruExpect `yaml:"fru"`
}

// System holds the information for system specific tests.
type System struct {
	GUID string `yaml:"guid"`
}

// PSU holds the information for one power supply unit.
type PSU struct {
	Name              string  `yaml:"name"`
	Manufacturer      string  `yaml:"manufacturer"`
	Model             string  `yaml:"model"`
	Partnumber        string  `yaml:"partNumber"`
	SerialNumber      string  `yaml:"serialNumber"`
	EfficiencyPercent float64 `yaml:"efficiencyPercent"`
}

// FruExpect holds the information of one field replaceable unit for validation.
type FruExpect struct {
	ProductManufacturer string `yaml:"productManufacturer"`
	ProductName         string `yaml:"productName"`
	ProductSerial       string `yaml:"productSerial"`
	ProductVersion      string `yaml:"productVersion"`
	BoardManufacturer   string `yaml:"boardManufacturer"`
	BoardSerial         string `yaml:"boardSerial"`
	BoardPartNumber     string `yaml:"boardPartnumber"`
	BoardProduct        string `yaml:"boardProduct"`
}

// DimmInfo holds the information of one dimm for validation.
type DimmInfo struct {
	Count             int    `yaml:"count"`
	PartNumber        string `yaml:"partNumber"`
	Manufacturer      string `yaml:"manufacturer"`
	SerialNumber      string `yaml:"serial"`
	Capacity          int    `yaml:"capacity"`
	OperatingSpeedMHz int    `yaml:"operatingSpeedMhz"`
}

// CPU holds the information of the systems CPU for validation.
type CPU struct {
	Manufacturer   string `yaml:"manufacturer"`
	Model          string `yaml:"model"`
	Architecture   string `yaml:"architecture"`
	InstructionSet string `yaml:"instructionSet"`
	Type           string `yaml:"type"`
}

// Host holds the information to connect to the specified host.
type Host struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSHKey   string `yaml:"sshKey"`
	SSHPort  int    `yaml:"sshPort"`
}

// BMC holds the information to connect to the specified BMC.
type BMC struct {
	BMCHost        string `yaml:"host"`
	BMCRedfishPort int    `yaml:"redfishPort"`
	BMCUser        string `yaml:"user"`
	BMCPassword    string `yaml:"password"`
	BMCSSHKey      string `yaml:"sshKey"`
	SSHPort        int    `yaml:"sshPort"`
}

// RedfishPort returns the configured redfish port, fallback to default
func (b BMC) RedfishPort() int {
	if b.BMCRedfishPort == 0 {
		return 443
	}

	return b.BMCRedfishPort
}

// Firmware holds the paths for the binary running on the BMC and an golden image for
// tests of update, downgrade and reboot tests.
type Firmware struct {
	Binary      string `yaml:"ciImage"`
	GoldenImage string `yaml:"goldenImage"`
}

// Config is the main configuration structure holding device information an validation values.
type Config struct {
	BMC      `yaml:"bmcConfig"`
	Firmware `yaml:"firmware"`
	TestData `yaml:"testdata"`

	Version  int      `yaml:"version"`
	Name     string   `yaml:"name"`
	NumHosts int      `yaml:"numHosts"`
	Hosts    []Host   `yaml:"hosts"`
	Tests    []string `yaml:"tests"`
}

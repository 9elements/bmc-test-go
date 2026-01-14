// Package testdevice describes a device with bmc and one or multiple hosts which will be tested.
package testdevice

import (
	"fmt"
	"log"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/reporting"
)

// Device holds the basic structure for the server under test.
// It can hold one BMC and multiple hosts.
type Device struct {
	*BMC
	*reporting.Reporter

	Hosts []*Host
}

// NewDevice creates a Device structure according to the configuration file and some commandline arguments.
func NewDevice(execEnv string, logType string, logFile string, cfg *configuration.Config) (*Device, error) {
	bmc, err := newBMC(execEnv, cfg)
	if err != nil {
		return nil, err
	}

	hosts := make([]*Host, 0)

	for _, hostCfg := range cfg.Hosts {
		if hostCfg.Password == "" && hostCfg.SSHKey == "" {
			continue
		}

		host, err := newHost(hostCfg)
		if err != nil {
			return nil, err
		}

		hosts = append(hosts, host)
	}

	rep, err := reporting.SetupReporter(logType, logFile)
	if err != nil {
		return nil, fmt.Errorf("failed to set up reporting: %w", err)
	}

	if len(hosts) < 1 {
		log.Printf("No host tests available, missing ssh password/keyfile")
	}

	return &Device{BMC: bmc, Reporter: rep, Hosts: hosts}, nil
}

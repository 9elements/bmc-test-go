package testdevice

import (
	"errors"
	"fmt"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/stmcginnis/gofish"
)

var errUnknownConenction = errors.New("unknown connection type")

// BMC holds the connection information to the bmc. It's local or remote.
// Also a redfish Service structure via gofish library is provided.
type BMC struct {
	con            BMCConnection
	redfishService *gofish.Service
}

func makeGofishConfig(config *configuration.Config) gofish.ClientConfig {
	return gofish.ClientConfig{
		Endpoint: "https://" + config.BMCHost,
		Username: config.BMCUser,
		Password: config.BMCPassword,
		Insecure: true,
	}
}

func newBMC(execEnv string, cfg *configuration.Config) (*BMC, error) {
	ret := &BMC{}

	switch execEnv {
	case "local":
		ret.con = &LocalConn{}
	case "remote":
		hostport := fmt.Sprintf("%s:%s", cfg.BMCHost, cfg.SSHPort)

		remoteConn, err := NewRemoteConn(hostport, cfg.BMCUser, cfg.BMCPassword)
		if err != nil {
			return nil, err
		}

		ret.con = remoteConn
	default:
		return nil, fmt.Errorf("%w: %s", errUnknownConenction, execEnv)
	}

	conn, err := gofish.Connect(makeGofishConfig(cfg))
	if err != nil {
		return nil, fmt.Errorf("gofish connect failed: %w", err)
	}

	ret.redfishService = conn.Service

	return ret, nil
}

// Close closes the connection to the BMC and delegates the call to the connection structure.
func (b *BMC) Close() error {
	err := b.con.Close()
	if err != nil {
		return fmt.Errorf("closing bmc connection failed: %w", err)
	}

	return nil
}

// ExecuteCommandLine sends a commandline call to the bmc.
func (b *BMC) ExecuteCommandLine(cmd string, args ...string) ([]byte, error) {
	resp, err := b.con.ExecuteCmdline(cmd, args...)
	if err != nil {
		return nil, fmt.Errorf("command returned an error: %w", err)
	}

	return resp, nil
}

// RedfishService returns the gofish.Service structure attached to the bmc.
func (b *BMC) RedfishService() *gofish.Service {
	return b.redfishService
}

// Reconnect reconnects to the bmc.
func (b *BMC) Reconnect(cfg *configuration.Config) error {
	conn, err := gofish.Connect(makeGofishConfig(cfg))
	if err != nil {
		return fmt.Errorf("reconnect failed: %w", err)
	}

	b.redfishService = conn.Service

	return nil
}

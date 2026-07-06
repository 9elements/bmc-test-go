package testdevice

import (
	"errors"
	"fmt"
	"net"
	"strconv"

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
	hostport := net.JoinHostPort(config.BMCHost, strconv.Itoa(config.RedfishPort()))

	return gofish.ClientConfig{
		Endpoint: "https://" + hostport,
		Username: config.BMCUser,
		Password: config.BMCPassword,
		Insecure: true,
	}
}

func newBMC(execEnv string, config *configuration.Config) (*BMC, error) {
	ret := &BMC{}

	switch execEnv {
	case "local":
		ret.con = &LocalConn{}
	case "remote":
		hostport := net.JoinHostPort(config.BMCHost, strconv.Itoa(config.SSHPort))

		remoteConn, err := NewRemoteConn(hostport, config.BMCUser, config.BMCPassword)
		if err != nil {
			return nil, err
		}

		ret.con = remoteConn
	default:
		return nil, fmt.Errorf("%w: %s", errUnknownConenction, execEnv)
	}

	conn, err := gofish.Connect(makeGofishConfig(config))
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
func (b *BMC) Reconnect(config *configuration.Config) error {
	conn, err := gofish.Connect(makeGofishConfig(config))
	if err != nil {
		return fmt.Errorf("reconnect failed: %w", err)
	}

	b.redfishService = conn.Service

	return nil
}

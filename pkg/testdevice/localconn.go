package testdevice

import (
	"context"
	"fmt"
	"os/exec"
)

// LocalConn describes a local connection, which means the tests for the bmc are executed on the
// locally and not over a remote connection.
type LocalConn struct{}

// ExecuteCmdline delegates the command and arguments to the device/shell locally.
func (l *LocalConn) ExecuteCmdline(cmd string, args ...string) ([]byte, error) {
	ctx := context.Background()
	command := exec.CommandContext(ctx, cmd, args...) // #nosec G204

	resp, err := command.Output()
	if err != nil {
		return resp, fmt.Errorf("command output received an error: %w", err)
	}

	return resp, nil
}

// Close does nothing, but is required to fulfil the BMCConnection interface.
func (l *LocalConn) Close() error {
	return nil
}

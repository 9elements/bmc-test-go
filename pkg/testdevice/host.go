package testdevice

import (
	"fmt"
	"log"
	"os"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"golang.org/x/crypto/ssh"
)

// Host holds connection information to a host under control of one BMC.
type Host struct {
	Name       string
	hostport   string
	sshConfig  *ssh.ClientConfig
	remoteConn *RemoteConn
}

func newHost(hostCfg configuration.Host) (*Host, error) {
	var sshConfig *ssh.ClientConfig

	if hostCfg.SSHKey != "" {
		keyBytes, err := os.ReadFile(hostCfg.SSHKey)
		if err != nil {
			return nil, fmt.Errorf("reading host ssh key file failed: %w", err)
		}

		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("parsing private ssh  host key failed: %w", err)
		}

		sshConfig = &ssh.ClientConfig{
			User:            hostCfg.Username,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // #nosec G106
		}
	} else {
		sshConfig = &ssh.ClientConfig{
			User:            hostCfg.Username,
			Auth:            []ssh.AuthMethod{ssh.Password(hostCfg.Password)},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(), // #nosec G106
		}
	}

	return &Host{
		Name:      hostCfg.Name,
		hostport:  fmt.Sprintf("%s:%s", hostCfg.IP, hostCfg.SSHPort),
		sshConfig: sshConfig,
	}, nil
}

// Close closes the connection to the host in case of ssh (remote) connection.
func (h *Host) Close() error {
	if h.remoteConn == nil {
		return nil
	}

	err := h.remoteConn.Close()
	if err != nil {
		return fmt.Errorf("closing connection failed: %w", err)
	}

	return nil
}

// ExecuteCommandLine sends a commandline call to the host.
func (h *Host) ExecuteCommandLine(cmd string, args ...string) ([]byte, error) {
	err := h.connect()
	if err != nil {
		return nil, err
	}

	defer func() {
		err := h.remoteConn.Close()
		if err != nil {
			log.Print(err)
		}

		h.remoteConn = nil
	}()

	return h.remoteConn.ExecuteCmdline(cmd, args...)
}

func (h *Host) connect() error {
	if h.remoteConn != nil {
		return nil
	}

	client, err := ssh.Dial("tcp", h.hostport, h.sshConfig)
	if err != nil {
		return fmt.Errorf("ssh dial failed: %w", err)
	}

	h.remoteConn = &RemoteConn{
		IP:     h.hostport,
		User:   h.sshConfig.User,
		client: client,
	}

	return nil
}

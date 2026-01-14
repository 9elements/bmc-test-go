package testdevice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"golang.org/x/crypto/ssh"
)

// RemoteConn holds all information for a remote bmc connection and the ssh client structure.
type RemoteConn struct {
	IP       string
	User     string
	Password string
	client   *ssh.Client
}

// NewRemoteConn creates a RemoteConn structure and configures accoding to parameters.
func NewRemoteConn(ip string, user string, password string) (*RemoteConn, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	}

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey() // #nosec G106

	client, err := ssh.Dial("tcp", ip, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial tcp connection: %w", err)
	}

	return &RemoteConn{
		IP:       ip,
		User:     user,
		Password: password,
		client:   client,
	}, nil
}

// Close closes the remote connection.
func (r *RemoteConn) Close() error {
	err := r.client.Close()
	if err != nil {
		return fmt.Errorf("closing remote connection failed: %w", err)
	}

	return nil
}

// ExecuteCmdline sends the cmd and args via an ssh session to the remote bmc.
func (r *RemoteConn) ExecuteCmdline(cmd string, args ...string) ([]byte, error) {
	session, err := r.client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("creating new ssh session failed: %w", err)
	}

	defer func() {
		err := session.Close()
		if err != nil && !errors.Is(err, io.EOF) {
			log.Print(err)
		}
	}()

	var cmdsArgs strings.Builder

	cmdsArgs.WriteString(cmd)

	for _, arg := range args {
		cmdsArgs.WriteString(" ")
		cmdsArgs.WriteString(arg)
	}

	resp, err := session.CombinedOutput(cmdsArgs.String())
	if err != nil {
		return resp, fmt.Errorf("command retured with errors: %w", err)
	}

	return resp, nil
}

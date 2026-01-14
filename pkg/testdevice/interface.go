package testdevice

// BMCConnection is the interface which describes a BMC connection and can be used for
// local or remote connection type.
type BMCConnection interface {
	Close() error
	ExecuteCmdline(cmd string, args ...string) ([]byte, error)
}

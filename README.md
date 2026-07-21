# bmc-test-go

`bmc-test-go` is a Go-based test framework for validating BMC and OpenBMC systems.
It provides a command-line test runner, YAML-based device configuration, Redfish
access through `gofish`, and local or SSH-based command execution on the BMC.

## Dependencies (Arch Linux)

```sh
sudo pacman -S go go-task
```

## Building

```sh
go-task build
```

## Example

```sh
./bin/bmc-test-go-amd64 suite -config contrib/example_config.yaml -suite redfish -exec-env local
```

## Features

- Run tests locally or against a remote BMC over SSH.
- Connect to the Redfish service over HTTPS.
- Group tests into suites such as IPMI, Redfish, SMBIOS, and BMC Linux checks.
- Store platform-specific expectations in a YAML configuration file.
- Emit either default text logs or structured JSON logs.

## Scope

`bmc-test-go` is intended for platform and integration testing of BMC behavior.
It is NOT a replacement for Redfish specification conformance tooling.

For Redfish-focused validation, use the DMTF tools:

- [Redfish Service Validator](https://github.com/DMTF/Redfish-Service-Validator):
  validates a Redfish service against Redfish CSDL schema. The tool is
  implementation-agnostic and performs `GET` requests to verify responses.
- [Redfish Interop Validator](https://github.com/DMTF/Redfish-Interop-Validator):
  validates a Redfish service against a Redfish interoperability profile.

Typical use of `bmc-test-go` is different: it can combine Redfish checks with
BMC shell commands, host-side checks, firmware update flows, platform-specific
expectations, and other system behavior that is outside pure Redfish schema or
profile compliance.

Consider for example checking for the expected asset tracking information or
sensor count.


## Configuration

The test runner expects a YAML configuration file.
Pass a different configuration file with the command-line configuration option.

Common environment variables:

- `IMAGE`: image under test
- `GOLDEN_IMAGE`: known-good image used by downgrade or recovery tests
- `BMC_SSH_KEY`: SSH key for the BMC root user
- `HOST_SSH_KEY`: SSH key for the host root user

When a configuration field references an environment variable and the expanded
value is empty, configuration loading fails.


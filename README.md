# bmc-test-go
bmc testing framework

## Environment Variables for config

The config parsing checks for the following environment variables:
- IMAGE : The base image to check
- GOLDEN_IMAGE : Known good image to downgrade to
- BMC_SSH_KEY : SSH Key for bmc root user
- HOST_SSH_KEY : SSH Key for the host root user

If the configuration file has the environment variable set, but the value is empty, the parser detects an error.
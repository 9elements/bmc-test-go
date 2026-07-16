#!/bin/bash

# example test suite execution with qemu catalina as the DUT

./bin/bmc-test-go-linux-amd64-v0.0.0 suite -config contrib/qemu-catalina.yaml -suite redfish -exec-env local

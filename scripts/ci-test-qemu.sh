#!/bin/bash

set -eu

./scripts/qemu-bmc.sh &

URL="https://127.0.0.1:4403/redfish/v1"

until curl \
	--silent \
	--insecure \
	--fail \
	--user root:0penBmc \
	"$URL" >/dev/null; do
	echo "Waiting for Redfish..."
	sleep 5
done

echo "sleep 40 seconds to wait for services to initialize"
sleep 40

./bin/bmc-test-go-linux-amd64-v0.0.0 suite \
 -config contrib/qemu-catalina.yaml \
 -suite redfish \
 -exec-env local

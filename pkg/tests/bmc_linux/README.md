
# BMC Linux tests

Tests for several linux system properties of the openbmc linux userland.

| No |  Name  | Status  | Precondition | Description|
|---|--------|---------|--------------|------------|
| 0 | BMCLinuxStartupHealthTest      | :white_check_mark: |   | Checks dmesg for some failure in linux kernel boot log |
| 1 | BMCLinuxSystemctlFailedTest    | :white_check_mark: |   | Checks systemctl for failed services |
| 2 | BMCLinuxSystemctlJobListTest   | :white_check_mark: |   | Checks systemctl for waiting services to be started |

# Redfish Tests

Tests for several Redfish Endpoints and validty of values

| No |  Name  | Status  | Precondition | Description|
|---|--------|---------|--------------|------------|
| 0 | RedfishChassisTest                  | :white_check_mark: |   | Checks if a chassis is returned |
| 1 | RedfishSystemTest                   | :white_check_mark: |   | Checks if a system is returned |
| 2 | RedfishFirmwareversionTest          | :white_check_mark: |   | Checks if a firmware version is returned |
| 3 | RedfishPSUInfoTest                  | :white_check_mark: |   | Checks if expected power supply information is returned (expected defined in config file) |
| 4 | RedfishVoltageSensorNameTest        | :white_check_mark: |   | Checks if expected voltage sensor names are returned  (expected defined in config file)|
| 5 | RedfishVoltageSensorValueTest       | :white_check_mark: |   | Checks if voltage sensor values are valid |
| 6 | RedfishTemperatureSensorNameTest    | :white_check_mark: |   | Checks if expected temperature sensor names are returned  (expected defined in config file) |
| 7 | RedfishTemperatureSensorValueTest   | :white_check_mark: |   | Checks if temperature sensor values are valid |
| 8 | RedfishFanSensorNameTest            | :white_check_mark: |   | Checks if expected thermal fan sensor names are returned (expected defined in config file)|
| 9 | RedfishFanSensorValueTest           | :white_check_mark: |   | Checks if thermal fan sensor values are valid |
| 10 | RedfishMemoryTest                  | :white_check_mark: |   | Checks expected memory/dimm information is returned (expected defined in config file)|
| 11 | RedfishProcessorTest               | :white_check_mark: |   | Checks if expected CPU information is returned (expected defined in config file)|

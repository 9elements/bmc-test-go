
# IPMI Tests

Tests for standard ipmitool commands

| No | Name                          | Status             | Description |
|---:|-------------------------------|--------------------|-------------|
|  0 | IPMISensorVoltageTest         | :white_check_mark: | Checks output of 'ipmitool sdr' for values from the config files testdata section |
|  1 | IPMISensorTemperaturTest      | :white_check_mark: | Checks output of 'ipmitool sdr' for values from the config files testdata section |
|  2 | IPMISensorFansTest            | :white_check_mark: | Checks output of 'ipmitool sdr' for values from the config files testdata section |
|  3 | IPMIPowerStatusTest           | :white_check_mark: | Checks output of 'ipmitool power status' against values from config files testdata section |
|  4 | IPMIChassisUptimeTest         | :white_check_mark: | Checks output of 'ipmitool chassis poh' against a precompiles regex pattern |
|  5 | IPMIFruTest                   | :white_check_mark: | Checks output of 'ipmitool fru' against values from the config files testdata-fru section |
|  6 | IPMIDCMIPowerReadingTest      | :x:                ||
|  7 | IPMIWatchdogConfigurationTest | :white_check_mark: | Checks output of 'ipmitool mc watchdog get' |
|  8 | IPMISystemGUIDTest            | :white_check_mark: | Checks output of 'ipmitool mc guid' against a predefines GUID against config file value |

# IPMI MCT Tests

Tests for ipmitool MCT OEM Extention commands

| No | Name                                  | Status             | Description |
|---:|---------------------------------------|--------------------|-------------|
|  0 | MCTIPMIPWMDutySetTest                 | :white_check_mark: | Checks output of PWM Duty Set |
|  1 | MCTIPMIPWMDutyGetTest                 | :white_check_mark: | Checks output of PWM Duty Get |
|  2 | MCTIPMIManufactureModeGetTest         | :white_check_mark: | Checks output of Get Manufacturer mode |
|  3 | MCTIPMIManufactureModeSetTest         | :white_check_mark: | Checks output of Set Manufacturer mode |
|  4 | MCTIPMIFloorDutyGetTest               | :white_check_mark: | Checks output of Get Floor Duty |
|  5 | MCTIPMIGetFruFieldTest                | :white_check_mark: | Checks output of Get Fru EEPROM fields |
|  6 | MCTIPMISetFruFieldTest                | :x:                | Checks output of Set Fru EEPROM fields |
|  7 | MCTIPMIGetFirmwareStringTest          | :white_check_mark: | Checks output of Get Firmware string |
|  8 | MCTIPMIConfigECCLeakyBucketSetTest    | :white_check_mark: | Checks output of Set Config ECC Leaky Bucket |
|  9 | MCTIPMIGPIOStatusTest                 | :white_check_mark: | Checks output of GPIO Status |

# IPMI Twitter Tests

Tests for ipmitool Twitter OEM Extention commands

| No | Name                                      | Status             | Description |
|---:|-------------------------------------------|--------------------|-------------|
|  0 | TwitterIPMIExtSetService                  | :white_check_mark: | Checks output of Set Service |
|  1 | TwitterIPMIExtGetService                  | :white_check_mark: | Checks output of Get Service |
|  2 | TwitterIPMIExtClearCMOS                   | :x:                | Checks output of Clear CMOS |
|  3 | TwitterIPMIExtPnmGetReading               | :white_check_mark: | Checks output of Get PNM Reading |
|  4 | TwitterIPMIExtRandomDelayACRestorePowerOn | :white_check_mark: | Checks output of Random Delay AC Restore Power On |
|  5 | TwitterIPMIExtGetPostCodes                | :white_check_mark: | Checks output of Get Post Codes |

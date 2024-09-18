// Package testcase implements the common structs used for defining
// testcases.
package testcase

// Config struct is used to pass arbitrary information into the tests
// such as IP address, expected FRU content, expected sensors
type Config struct {
	Host     string
	Username string
	Password string

	// ... to be extended
}

// Result struct is used for returning test results from running a test
type Result struct {
	// what was tested
	Name string

	// != nil in case something went wrong,
	// this should contain a descriptive string for the logs.
	// e.g. "while running command X, expected A but output was B'
	Error error

	// to store stdout/stderr from a command, if any
	Output string
}

// Test struct is used to define a test which can return one or more results
type Test struct {
	// e.g. "IPMI Sensor"
	Name string

	// e.g. "Testing the 'ipmitool sensor' output
	Description string

	Function func(tc Config) []Result

	Results []Result
}

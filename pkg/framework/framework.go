// Package framework provides the basic structure to implement test suites.
package framework

import (
	"errors"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
)

// Result implements a specific type for a test result for easy string conversion.
type Result int

// The following constants indicate the result of a test.
const (
	ResultNotRun Result = 0 + iota
	ResultDependencyFailure
	ResultInternalFailure
	ResultNotImplemented
	ResultFail
	ResultSuccess
)

func (r Result) String() string {
	return [...]string{"TESTNOTRUN", "DEPENDENCY_FAILED", "INTERNAL_ERROR", "Not Implemented", "FAIL", "PASS"}[r]
}

// Status implements a specific type for the implementation status of a test for easy string conversion.
type Status int

// The following constants indicate the implementation status of a test.
const (
	StatusImplemented Status = 0 + iota
	StatusNotImplemented
	StatusPartlyImplemented
)

// ErrNotImplemented serves as the default error returned by functions
// not being implemented yet.
var ErrNotImplemented = errors.New("not implemented")

func (s Status) String() string {
	return [...]string{"Implemented", "Not implemented", "Partly implemented"}[s]
}

// Test represents a one specific test.
type Test struct {
	Name      string
	ShortName string
	Status    Status
	Result    Result
	Function  func(*testdevice.Device, *configuration.Config) (bool, error, error)
	ErrorText string
}

// Run executes the test function and sets the error test and test result values.
func (t *Test) Run(dev *testdevice.Device, cfg *configuration.Config) bool {
	rc, testError, internalError := t.Function(dev, cfg)
	if testError != nil {
		t.ErrorText = testError.Error()
		t.Result = ResultFail

		return rc
	}

	if internalError != nil {
		t.ErrorText += internalError.Error()
		t.Result = ResultInternalFailure

		return rc
	}

	t.Result = ResultSuccess

	return rc
}

// RunTest is the function which manages the execution of all tests.
func RunTest(test *Test, dev *testdevice.Device, testCfg *configuration.Config) bool {
	if test.Status.String() == StatusNotImplemented.String() {
		return true
	}

	ret := test.Run(dev, testCfg)

	isError := test.ErrorText != "" || test.Result.String() == ResultFail.String()
	if isError {
		dev.Reporter.WithField("Name", test.Name).WithField("Result", test.Result.String()).
			WithField("Error", test.ErrorText).Info("Failed")
	} else {
		dev.Reporter.WithField("Name", test.Name).WithField("Result", test.Result.String()).Info("Success")
	}

	return ret
}

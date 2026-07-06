// Package main implements the core logic
package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/9elements/bmc-test-go/pkg/configuration"
	"github.com/9elements/bmc-test-go/pkg/framework"
	"github.com/9elements/bmc-test-go/pkg/testdevice"
	"github.com/9elements/bmc-test-go/pkg/tests"
)

var errNoTestsToExecute = errors.New("no tests to execute")

func runAll(dev *testdevice.Device, cfg *configuration.Config) error {
	tests := tests.AllTests()
	if len(tests) < 1 {
		return errNoTestsToExecute
	}

	for _, test := range tests {
		if !framework.RunTest(test, dev, cfg) {
			log.Printf("%s", test.ErrorText)

			continue
		}
	}

	return nil
}

func runSingle(testname string, dev *testdevice.Device, cfg *configuration.Config) error {
	tests := tests.AllTests()
	if len(tests) < 1 {
		return errNoTestsToExecute
	}

	for _, test := range tests {
		if test.ShortName == testname {
			if !framework.RunTest(test, dev, cfg) {
				log.Printf("%s", test.ErrorText)
			}
		}
	}

	return nil
}

func runSuite(suite string, dev *testdevice.Device, cfg *configuration.Config) error {
	tests := tests.GetSuite(suite)
	if len(tests) < 1 {
		return errNoTestsToExecute
	}

	logger := log.New(os.Stdout, "", log.Lmsgprefix)

	for _, test := range tests {
		if !framework.RunTest(test, dev, cfg) {
			logger.Printf("%s", test.ErrorText)
		}
	}

	return nil
}

func listTests(suite string) error {
	var t []*framework.Test
	if suite != "" {
		t = tests.GetSuite(suite)
	} else {
		t = tests.AllTests()
	}

	logger := log.New(os.Stdout, "", log.Lmsgprefix)

	logger.Printf("List tests %s", suite)
	logger.Printf("-------------------------------------------------------")
	logger.Printf("-------------------------------------------------------")
	logger.Printf("%-30s | %s\n", "Long name", "Short name")
	logger.Printf("-------------------------------------------------------")

	for _, test := range t {
		logger.Printf("%-30s | %s\n", test.Name, test.ShortName)
		logger.Printf("-------------------------------------------------------")
	}

	return nil
}

func runFromCfg(dev *testdevice.Device, cfg *configuration.Config) {
	alltests := tests.AllTests()

	for _, testName := range cfg.Tests {
		for _, test := range alltests {
			if testName == test.ShortName {
				framework.RunTest(test, dev, cfg)
			}
		}
	}
}

func listSuites() {
	logger := log.New(os.Stdout, "", log.Lmsgprefix)

	logger.Printf("List Suites")
	logger.Printf("-------------------------------------------------------")
	logger.Printf("-------------------------------------------------------")

	for name := range tests.Suites {
		logger.Println(name)
	}
}

func setupSuite(f *flags) (*testdevice.Device, *configuration.Config, error) {
	cfg, err := configuration.LoadConfig(f.configPath)
	if err != nil {
		return nil, nil, fmt.Errorf("loading configuration file failed: %w", err)
	}

	bmc, err := testdevice.NewDevice(f.execEnv, f.logType, f.logFile, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set up new device: %w", err)
	}

	return bmc, cfg, nil
}

func runCfgValidation(cfg *configuration.Config) error {
	err := configuration.Validate(cfg)
	if err != nil {
		return fmt.Errorf("validation of cfg failed: %w", err)
	}

	return nil
}

var (
	errUnknownCommand       = errors.New("unknown command")
	errInvalidArgumentCount = errors.New("invalid argument count")
)

func runAgainstDevice(f *flags) error {
	bmc, cfg, err := setupSuite(f)
	if err != nil {
		return fmt.Errorf("failed to set up device: %w", err)
	}

	defer func() {
		err := bmc.Close()
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			log.Print(err)
		}
	}()

	switch f.cmd {
	case runSingleCmd:
		return runSingle(f.testName, bmc, cfg)
	case runSuiteCmd:
		return runSuite(f.suite, bmc, cfg)
	case runAllCmd:
		return runAll(bmc, cfg)
	case runCfgCmd:
		runFromCfg(bmc, cfg)
	case cfgCheckCmd:
		return runCfgValidation(cfg)
	default:
		return fmt.Errorf("%w: %s", errUnknownCommand, f.cmd)
	}

	return nil
}

func run(args []string) error {
	minArgLen := 1
	if len(args) < minArgLen {
		return fmt.Errorf("%w: %d", errInvalidArgumentCount, len(args))
	}

	flags, err := parseFlags(args)
	if err != nil {
		return err
	}

	switch flags.cmd {
	case listTestsCmd:
		return listTests(flags.suite)
	case listSuitesCmd:
		listSuites()

		return nil
	}

	return runAgainstDevice(flags)
}

func main() {
	err := run(os.Args[1:])
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

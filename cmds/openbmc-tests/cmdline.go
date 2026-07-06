package main

import (
	"errors"
	"flag"
	"fmt"
)

const (
	runSingleCmd  = "single"
	runSuiteCmd   = "suite"
	runAllCmd     = "all"
	cfgCheckCmd   = "cfg-check"
	listTestsCmd  = "list-tests"
	listSuitesCmd = "list-suites"
	runCfgCmd     = "run-from-cfg"
)

type flags struct {
	cmd        string
	configPath string
	execEnv    string
	suite      string
	testName   string
	logType    string
	logFile    string
}

func defineCfgCheckFlagSet(f *flags) *flag.FlagSet {
	cfgCheckFS := flag.NewFlagSet(cfgCheckCmd, flag.ExitOnError)
	cfgCheckFS.StringVar(&f.configPath, "config", "contrib/example_config.yaml",
		"Path to device specific test case configuration file")

	return cfgCheckFS
}

func defineRunSingleCmdFlagSet(f *flags) *flag.FlagSet {
	runSingleFS := flag.NewFlagSet(runSingleCmd, flag.ExitOnError)
	runSingleFS.StringVar(&f.testName, "testname", "", "Run a single test with 'testname'. Use the short name. No default")
	runSingleFS.StringVar(&f.execEnv, "exec-env", "remote",
		"Run the test locally or remotely '<local, remote>'. Default: remote")
	runSingleFS.StringVar(&f.configPath, "config", "contrib/example_config.yaml",
		"Path to device specific test case configuration file")
	runSingleFS.StringVar(&f.logType, "log", "default", "Logging options <default, structured>")
	runSingleFS.StringVar(&f.logFile, "logfile", "", "Path to log file. Only with structured logging option.")

	return runSingleFS
}

func defineRunSuiteCmdFlagSet(f *flags) *flag.FlagSet {
	runSuiteFS := flag.NewFlagSet(runSuiteCmd, flag.ExitOnError)
	runSuiteFS.StringVar(&f.execEnv, "exec-env", "remote",
		"Run the test locally or remotely '<local, remote>'. Default: remote")
	runSuiteFS.StringVar(&f.suite, "suite", "", "Specify the suite. <redfish, ipmi, ifaces>")
	runSuiteFS.StringVar(&f.configPath, "config", "contrib/example_config.yaml",
		"Path to device specific test case configuration file")
	runSuiteFS.StringVar(&f.logType, "log", "default", "Logging options <default, structured>")
	runSuiteFS.StringVar(&f.logFile, "logfile", "", "Path to log file. Only with structured logging option.")

	return runSuiteFS
}

func defineRunAllCmdFlagSet(f *flags) *flag.FlagSet {
	runAllFS := flag.NewFlagSet(runAllCmd, flag.ExitOnError)
	runAllFS.StringVar(&f.execEnv, "exec-env", "remote", "Run the test locally or remotely '<local, remote>'.")
	runAllFS.StringVar(&f.configPath, "config", "contrib/example_config.yaml",
		"Path to device specific test case configuration file")
	runAllFS.StringVar(&f.logType, "log", "default", "Logging options <default, structured>")
	runAllFS.StringVar(&f.logFile, "logfile", "", "Path to log file. Only with structured logging option.")

	return runAllFS
}

func defineListTestCmdFlagSet(f *flags) *flag.FlagSet {
	listTestsFS := flag.NewFlagSet(listTestsCmd, flag.ExitOnError)
	listTestsFS.StringVar(&f.logType, "log", "default", "Logging options <default, structured>")
	listTestsFS.StringVar(&f.logFile, "logfile", "", "Path to log file. Only with structured logging option.")
	listTestsFS.StringVar(&f.suite, "suite", "", "Specify the suite. <redfish, ipmi, ifaces>")

	return listTestsFS
}

func defineListSuitesCmdsFlagSet() *flag.FlagSet {
	return flag.NewFlagSet(listSuitesCmd, flag.ExitOnError)
}

func defineRunFromCfgCmdsFlagSet(f *flags) *flag.FlagSet {
	runFromCfgFS := flag.NewFlagSet(runCfgCmd, flag.ExitOnError)
	runFromCfgFS.StringVar(&f.execEnv, "exec-env", "remote",
		"Run the test locally or remotely '<local, remote>'. Default: remote")
	runFromCfgFS.StringVar(&f.configPath, "config", "contrib/example_config.yaml",
		"Path to device specific test case configuration file")
	runFromCfgFS.StringVar(&f.logType, "log", "default", "Logging options <default, structured>")
	runFromCfgFS.StringVar(&f.logFile, "logfile", "", "Path to log file. Only with structured logging option.")

	return runFromCfgFS
}

func defineFlagSets(f *flags) map[string]*flag.FlagSet {
	flagsets := map[string]*flag.FlagSet{
		cfgCheckCmd:   defineCfgCheckFlagSet(f),
		runSingleCmd:  defineRunSingleCmdFlagSet(f),
		runAllCmd:     defineRunAllCmdFlagSet(f),
		runSuiteCmd:   defineRunSuiteCmdFlagSet(f),
		listTestsCmd:  defineListTestCmdFlagSet(f),
		listSuitesCmd: defineListSuitesCmdsFlagSet(),
		runCfgCmd:     defineRunFromCfgCmdsFlagSet(f),
	}

	return flagsets
}

var errUnknownCmd = errors.New("unknown command provided or help called")

func parseFlags(args []string) (*flags, error) {
	f := &flags{
		cmd: args[0],
	}

	flagsets := defineFlagSets(f)

	var err error

	fs, ok := flagsets[f.cmd]
	if !ok {
		for _, fs := range flagsets {
			fs.Usage()
		}

		return nil, errUnknownCmd
	}

	err = fs.Parse(args[1:])
	if err != nil {
		return nil, fmt.Errorf("failed to parse arguments: %w", err)
	}

	return f, nil
}

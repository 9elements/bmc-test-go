// Package reporting creates the reporting mechanism for tests.
package reporting

import (
	"fmt"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

// Reporter holds the Logger.
type Reporter struct {
	*logrus.Logger
}

// SetupReporter creates a logrus Logger according to commandline parameters.
func SetupReporter(repType string, logfile string) (*Reporter, error) {
	file := os.Stdout

	var fm os.FileMode = 0o600

	var err error

	if logfile != "" {
		file, err = os.OpenFile(path.Clean(logfile), os.O_CREATE|os.O_RDWR, fm)
		if err != nil {
			return nil, fmt.Errorf("open or creat file failed: %w", err)
		}
	}

	logrusser := logrus.New()
	logrusser.SetOutput(file)

	if repType == "structured" {
		logrusser.SetFormatter(&logrus.JSONFormatter{})
	}

	return &Reporter{logrusser}, nil
}

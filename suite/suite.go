// Package suite interfaces with the legacy suite.staffjoy.com system
package suite

import (
	"net/url"
	"os"

	"github.com/Sirupsen/logrus"
	"v2.staffjoy.com/environments"
)

const (
	apiKeyEnvVar = "SUITE_API_KEY"
	// ServiceName is how this package identifies itself in logs
	ServiceName = "suite"
)

// SuiteConfigs is a map of environment to the url for
// the suite
var SuiteConfigs = map[string]url.URL{
	"development": {
		Scheme: "http",
		Host:   "suite.local",
	},
	"staging": {
		Scheme: "https",
		Host:   "stage.staffjoy.com",
	},
	"production": {
		Scheme: "https",
		Host:   "suite.staffjoy.com",
	},
}

var (
	logger *logrus.Entry
	config environments.Config
	apiKey string
)

func init() {
	var err error
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine suite configuration")
	}
	logger = config.GetLogger(ServiceName)
	apiKey = os.Getenv(apiKeyEnvVar)
}

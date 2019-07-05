// Package environments determines code behavior across services
// such as debug mode.
//
// usage:
// config := environments.GetConfig(os.Getenv(environments.EnvVar))
package environments

import (
	"fmt"
	"os"

	sentry "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	intercom "gopkg.in/intercom/intercom-go.v2"
)

const (
	// DefaultEnv is the fallback environment
	DefaultEnv = "development"
	// EnvVar is the typical environment variable used to determine
	// the current environment, e.g. with      os.Getenv(EnvVar)
	EnvVar = "ENV"
	// SentryEnvVar is the environment variable name for fetching this secret from k8s
	SentryEnvVar = "SENTRY_DSN"
	// DeployEnvVar is set by Kubernetes during a new deployment so we can identify the code version
	DeployEnvVar = "DEPLOY"
	// GoogleCloudSecretPath is the file location for the mounted GCloud auth token
	GoogleCloudSecretPath = "/etc/secrets/gcloud.json"
)

// Config controls behavior for the environment across services
type Config struct {
	Name         string       // Name of the environment, which shows up in logs
	Debug        bool         // Controls security and log verbosity
	ExternalApex string       // Apex domain off of which services operate externally
	InternalApex string       // Apex domain off of which services operate internally
	LogLevel     logrus.Level // Verbosity of logging
	Scheme       string       // default URL scheme - http or https
}

var configs = map[string]Config{
	"development": {
		Name:         "development",
		Debug:        true,
		ExternalApex: "staffjoy-v2.local",
		InternalApex: "development",
		LogLevel:     logrus.DebugLevel,
		Scheme:       "http",
	},
	"test": {
		Name:         "test",
		Debug:        true,
		ExternalApex: "staffjoy-v2.local",
		InternalApex: "development",
		LogLevel:     logrus.ErrorLevel,
		Scheme:       "http",
	},
	"staging": {
		Name:         "staging",
		Debug:        false,
		ExternalApex: "staffjoystaging.com",
		InternalApex: "staging",
		LogLevel:     logrus.InfoLevel,
		Scheme:       "https",
	},
	"production": {
		Name:         "production",
		Debug:        false,
		ExternalApex: "staffjoy.com",
		InternalApex: "production",
		LogLevel:     logrus.InfoLevel,
		Scheme:       "https",
	},
}

// GetConfig returns the environment config
// Typically you should use os.Getenv("ENV") and pass that in here,
func GetConfig(envName string) (conf Config, err error) {
	if envName == "" {
		envName = DefaultEnv
	}

	conf, ok := configs[envName]
	if ok == false {
		err = fmt.Errorf("Unable to determine environment on boot (environment name %s not found)", envName)
	}
	return
}

// GetIntercomClient reutrns an intercom.io client
func (c *Config) GetIntercomClient() *intercom.Client {
	return intercom.NewClient(os.Getenv("INTERCOM_APP_ID"), os.Getenv("INTERCOM_API_KEY"))
}

// GetLogger returns a structured logger
func (c *Config) GetLogger(serviceName string) *logrus.Entry {
	// For a valid config, we should not have an err
	// Configure logger. This is Staffjoy standard format.
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(c.LogLevel)
	return logrus.WithFields(logrus.Fields{
		"env":     c.Name,
		"service": serviceName,
	})
}

// GetSentryDSN returns the secret API key for Sentry.io
func (c *Config) GetSentryDSN() string {
	return os.Getenv(SentryEnvVar)
}

// GetDeployVersion returns the current code version from Jenkins in non-dev environments
func (c *Config) GetDeployVersion() string {
	return os.Getenv(DeployEnvVar)
}

func (c *Config) getSentryClient(dsn string) (*sentry.Client, error) {
	client, err := sentry.NewWithTags(dsn, map[string]string{"env": c.Name})
	if err != nil {
		return nil, err
	}
	// Set the code version so that we can track it
	client.SetRelease(c.GetDeployVersion())
	return client, nil
}

// Abastracted for interface
func (c *Config) isDebug() bool {
	return c.Debug
}

// Package environments determines code behavior across services
// such as debug mode.
//
// usage:
// config := environments.GetConfig(os.Getenv(environments.EnvVar))
package environments

import (
	sentry "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

// SentryConfig is an interface for being able to retrieve error handler funcs
type SentryConfig interface {
	isDebug() bool
	GetSentryDSN() string
	GetDeployVersion() string
	getSentryClient(string) (*sentry.Client, error)
	GetLogger(string) *logrus.Entry
}

// SentryClient is a type for handling errors (without importing sentry directly)
type SentryClient interface {
	CapturePanicAndWait(func(), map[string]string, ...sentry.Interface) (interface{}, string)
	CaptureError(error, map[string]string, ...sentry.Interface) string
	SetHttpContext(*sentry.Http)
	SetUserContext(*sentry.User)
	ClearContext()
}

// ErrorClient returns a an error handler for sending to GetSentry.com
func ErrorClient(c SentryConfig) SentryClient {
	// Do not do anything in dev
	if c.isDebug() {
		return nil
	}
	logger := c.GetLogger("environments")

	dsn := c.GetSentryDSN()
	if dsn == "" {
		logger.Errorf("Unable to determine Sentry DSN")
		return nil
	}

	client, err := c.getSentryClient(dsn)
	if err != nil {
		logger.Errorf("Unable to open sentry client - %s", err)
		return nil
	}
	return client
}

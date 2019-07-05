package environments

import (
	"fmt"
	"os"
	"testing"

	sentry "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type FakeConfig struct {
	debug        bool
	dsn          string
	clientCalled bool
	deploy       string
}

func (f *FakeConfig) isDebug() bool {
	return f.debug
}

func (f *FakeConfig) GetSentryDSN() string {
	return f.dsn
}

func (f *FakeConfig) GetDeployVersion() string {
	return f.deploy
}

func (f *FakeConfig) getSentryClient(dsn string) (*sentry.Client, error) {
	f.clientCalled = true

	return &sentry.Client{Tags: map[string]string{"env": "test"}}, nil
}

func (f *FakeConfig) GetLogger(service string) *logrus.Entry {
	logrus.SetOutput(os.Stdout)
	return logrus.WithFields(logrus.Fields{
		"env": "test",
	})
}

func TestGetErrorClientDebugReturnsNil(t *testing.T) {
	assert := assert.New(t)
	conf := &FakeConfig{debug: true}

	errFunc := ErrorClient(conf)
	assert.Nil(errFunc)
}

func TestGetErrorClientSucceeds(t *testing.T) {
	assert := assert.New(t)
	testDsn := "dsn_not_dns"
	conf := &FakeConfig{debug: false, dsn: testDsn}
	fmt.Printf(conf.dsn)

	out := ErrorClient(conf)
	assert.NotNil(out)
	assert.True(conf.clientCalled)
}

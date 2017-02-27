package environments

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultWorks(t *testing.T) {
	assert := assert.New(t)
	explicitDefaultConfig, err := GetConfig(DefaultEnv)
	assert.NoError(err)
	implicitDefaultConfig, err := GetConfig("")
	assert.NoError(err)
	assert.Equal(explicitDefaultConfig, implicitDefaultConfig)
}

func TestGetEnvSucceeds(t *testing.T) {
	assert := assert.New(t)
	names := []string{"development", "staging", "production"}
	for _, name := range names {
		conf, err := GetConfig(name)
		assert.NoError(err)
		assert.NotNil(conf)
		assert.Equal(name, conf.Name)
	}
}

func TestUnknownEnvThrowsError(t *testing.T) {
	assert := assert.New(t)
	// invalid name
	_, err := GetConfig("FakeEnv")
	assert.Error(err)
}

func TestGetLoggerWorks(t *testing.T) {
	assert := assert.New(t)
	serviceName := "servicio"
	names := []string{"development", "staging", "production"}
	for _, name := range names {
		conf, err := GetConfig(name)
		assert.NoError(err)
		logger := conf.GetLogger(serviceName)
		assert.Equal(serviceName, logger.Data["service"])
		assert.Equal(name, logger.Data["env"])
	}
}

func TestGetPublicSentryDSN(t *testing.T) {
	assert := assert.New(t)

	secretDSN := "https://username:password@sentry.io/12345"
	expectedPublicDSN := "https://username@sentry.io/12345"
	actual, err := GetPublicSentryDSN(secretDSN)
	assert.NoError(err)
	assert.Equal(expectedPublicDSN, actual)
}

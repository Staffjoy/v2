package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"v2.staffjoy.com/environments"

	"github.com/stretchr/testify/assert"
)

type MobileConfigDecoder struct {
	Pattern string `json:"hideNavForURLsMatchingPattern"`
}

func TestConfigPathReturnsConfig(t *testing.T) {
	assert := assert.New(t)

	rec := httptest.NewRecorder()

	conf, err := environments.GetConfig("test")
	assert.NoError(err)
	r := NewRouter(conf, conf.GetLogger("test"))

	req, err := http.NewRequest(http.MethodGet, "https://www.staffjoy.com/mobile-config.json", nil)
	assert.NoError(err)
	r.ServeHTTP(rec, req)

	body, err := ioutil.ReadAll(rec.Body)
	assert.NoError(err)
	assert.Equal(rec.Code, http.StatusOK)
	assert.Equal("application/json", rec.Header().Get("Content-Type"))

	decoded := &MobileConfigDecoder{}
	err = json.Unmarshal(body, decoded)
	assert.NoError(err)
	assert.Equal(MobileConfigRegex, decoded.Pattern)
}

func TestMobileConfigRegex(t *testing.T) {
	assert := assert.New(t)

	testDomains := []struct {
		domain  string
		allowed bool
	}{
		// internal domains
		{"http://dev.staffjoy.com", true},
		{"https://stage.staffjoy.com", true},
		{"https://www.staffjoy.com", true},
		{"https://suite.staffjoy.com", true},
		// external domains
		{"https://help.staffjoy.com", false},
		{"https://blog.staffjoy.com", false},
		{"https://google.com", false},
		{"http://7bridg.es", false},
	}

	for _, test := range testDomains {
		match, err := regexp.MatchString(MobileConfigRegex, test.domain)
		assert.NoError(err)
		assert.Equal(test.allowed, match)

	}
}

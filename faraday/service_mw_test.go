package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"v2.staffjoy.com/environments"
)

func TestApexRedirectsToWWW(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		debug       bool
		reqURL      string
		expectedURL string
	}{
		{true, "http://staffjoy.rocks/foo/bar?hello=world", "http://www.staffjoy.rocks/foo/bar?hello=world"},
		{false, "http://staffjoy.rocks/foo/bar?hello=world", "https://www.staffjoy.rocks/foo/bar?hello=world"},
		{true, "http://staffjoy.rocks", "http://www.staffjoy.rocks"},
	}

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

	for _, method := range methods {
		for _, test := range tests {
			config := environments.Config{Debug: test.debug, ExternalApex: "staffjoy.rocks"}
			mw := NewServiceMiddleware(config, nil)
			req, err := http.NewRequest(method, test.reqURL, nil)
			assert.NoError(err)
			rec := httptest.NewRecorder()

			nextCalled := false

			next := func(res http.ResponseWriter, req *http.Request) {
				nextCalled = true
			}

			// Run it!
			mw.ServeHTTP(rec, req, next)

			// Did not continue
			assert.False(nextCalled)

			res := rec.Result()
			assert.Equal(res.StatusCode, 301)
			assert.Equal(test.expectedURL, res.Header.Get("Location"))
		}
	}
}

func TestHostToService(t *testing.T) {
	assert := assert.New(t)

	apex := "staffjoy.com"

	noErrTests := map[string]string{
		"foo.staffjoy.com":          "foo",
		"a.b.staFfjoy.com":          "a.b",
		"staffjoy.com.staffjoy.COM": "staffjoy.com",
		"feynman.staffjoy.com:80":   "feynman",
	}
	for host, expected := range noErrTests {
		actual, err := HostToService(host, apex)
		assert.Nil(err)
		assert.Equal(actual, expected)
	}

	// These tests should throw an error
	errTests := []string{
		"foo.staffjoy.co",
		"foo.staffjoy.com.co",
		"feynman",
		"staffjoy.com",
		"moc.yojffats.sdrawkcab",
		"feynmanstaffjoy.com",
		"",
	}
	for _, host := range errTests {
		_, err := HostToService(host, apex)
		assert.NotNil(err)
	}
}

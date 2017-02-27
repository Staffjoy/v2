package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"v2.staffjoy.com/environments"
)

func TestDebugMakesNoChanges(t *testing.T) {
	assert := assert.New(t)
	forbiddenHeaders := []string{
		"Strict-Transport-Security",
		"X-XSS-Protection",
		"X-Frame-Options",
	}
	nextCalled := false
	rec := httptest.NewRecorder()

	mw := NewSecurityMiddleware(environments.Config{Debug: true})

	next := func(res http.ResponseWriter, req *http.Request) {
		nextCalled = true
	}
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	for _, method := range methods {

		req, err := http.NewRequest(method, "https://www.staffjoy.com/I/am/going/somewhere/cool", nil)
		assert.NoError(err)

		mw.ServeHTTP(rec, req, next)

		assert.True(nextCalled)
		for _, header := range forbiddenHeaders {
			val := rec.Header().Get(header)
			assert.Empty(val)
		}
	}
}

func TestNotDebugMakesChanges(t *testing.T) {
	assert := assert.New(t)
	requiredHeaders := []string{
		"Strict-Transport-Security",
		"X-XSS-Protection",
		"X-Frame-Options",
	}
	nextCalled := false

	mw := NewSecurityMiddleware(environments.Config{Debug: false})

	next := func(res http.ResponseWriter, req *http.Request) {
		nextCalled = true
	}

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	for _, method := range methods {
		nextCalled = false
		rec := httptest.NewRecorder()

		req, err := http.NewRequest(method, "https://www.staffjoy.com/I/am/going/somewhere/cool", nil)
		assert.NoError(err)

		mw.ServeHTTP(rec, req, next)

		assert.True(nextCalled)
		res := rec.Result()
		for _, header := range requiredHeaders {
			val := res.Header.Get(header)
			assert.NotEmpty(val)
		}
	}
}

func TestInsecureRequestRedirectsToHTTPS(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		debug              bool
		expectedNextCalled bool
		reqURL             string
		expectedURL        string
	}{
		{true, true, "http://staffjoy.rocks/foo/bar?hello=world", "http://www.staffjoy.rocks/foo/bar?hello=world"},
		{false, false, "http://staffjoy.rocks/foo/bar?hello=world", "https://staffjoy.rocks/foo/bar?hello=world"},
		{false, true, "https://staffjoy.rocks", ""},
	}

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

	for _, method := range methods {
		for _, test := range tests {
			config := environments.Config{Debug: test.debug}
			mw := NewSecurityMiddleware(config)

			req, err := http.NewRequest(method, test.reqURL, nil)
			assert.NoError(err)

			rec := httptest.NewRecorder()

			nextCalled := false

			next := func(res http.ResponseWriter, req *http.Request) {
				nextCalled = true
				res.WriteHeader(http.StatusOK)
			}

			// Run it!
			mw.ServeHTTP(rec, req, next)

			// Did not continue
			assert.Equal(test.expectedNextCalled, nextCalled)

			res := rec.Result()

			// If blocked execution, then should redirect
			if test.expectedNextCalled == false {
				assert.Equal(res.StatusCode, 301)
				assert.Equal(test.expectedURL, res.Header.Get("Location"))
			} else {
				assert.Equal(res.StatusCode, http.StatusOK)
			}
		}
	}
}

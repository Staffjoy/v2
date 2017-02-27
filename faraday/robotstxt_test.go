package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"v2.staffjoy.com/faraday/services"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
	"v2.staffjoy.com/environments"
)

func TestNotRobotsGoesNext(t *testing.T) {
	assert := assert.New(t)

	nextCalled := false
	config := environments.Config{Debug: false}
	rec := httptest.NewRecorder()

	mw := NewRobotstxtMiddleware(config)

	next := func(res http.ResponseWriter, req *http.Request) {
		nextCalled = true
	}

	req, err := http.NewRequest("GET", "http://example.com/totally/not/a/robot", nil)
	assert.NoError(err)
	mw.ServeHTTP(rec, req, next)
	assert.True(nextCalled)
	assert.Empty(ioutil.ReadAll(rec.Body))
}

func TestHits(t *testing.T) {
	assert := assert.New(t)

	var authTests = []struct {
		security       int
		configName     string
		expectedResult string
	}{
		{services.Public, "production", robotstxtAllow},
		{services.Public, "staging", robotstxtDeny},
		{services.Public, "development", robotstxtDeny},
		{services.Authenticated, "production", robotstxtDeny},
		{services.Authenticated, "staging", robotstxtDeny},
		{services.Admin, "production", robotstxtDeny},
		{services.Admin, "staging", robotstxtDeny},
	}

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}

	for _, method := range methods {
		for _, tt := range authTests {

			nextCalled := false
			rec := httptest.NewRecorder()

			mw := NewRobotstxtMiddleware(environments.Config{Debug: false, Name: tt.configName})

			next := func(res http.ResponseWriter, req *http.Request) {
				nextCalled = true
			}

			req, err := http.NewRequest(method, "https://www.staffjoy.com/robots.txt", nil)
			assert.NoError(err)

			service := services.Service{Security: tt.security}

			// We mock the upstream ServicesMiddleware
			context.Set(req, requestedService, service)

			mw.ServeHTTP(rec, req, next)

			assert.False(nextCalled)
			body, err := ioutil.ReadAll(rec.Body)
			assert.NoError(err)
			assert.Equal(tt.expectedResult, string(body))
		}
	}
}

package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
)

func TestStaticPagesRenderValidHtml(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	r := NewRouter()

	for path, page := range staticPages {
		req, err := http.NewRequest("GET", path, nil)
		assert.NoError(err)
		r.ServeHTTP(recorder, req)
		assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
		assert.Equal(http.StatusOK, recorder.Code)
		assert.NotZero(recorder.Body.Len())
		bodyHTML, err := ioutil.ReadAll(recorder.Body)
		assert.NoError(err)
		if page.Title != "" {
			assert.Contains(string(bodyHTML), page.Title)
		}
		assert.Contains(string(bodyHTML), "Staffjoy")

		if page.Description != "" {
			assert.Contains(string(bodyHTML), page.Description)
		} else {
			assert.Contains(string(bodyHTML), defaultDescription)
		}
		// Check for valid html
		_, err = html.Parse(recorder.Body)
		assert.NoError(err)
	}
}

func TestUnknownPageReturnsNotFound(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	router := NewRouter()
	unknownRoute := "/this-does-not-exist"

	req, err := http.NewRequest("GET", unknownRoute, nil)
	assert.NoError(err)

	router.ServeHTTP(recorder, req)
	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusNotFound, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Oops")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// Sanity check
func TestHomepageExists(t *testing.T) {
	homepageRoute := "/"
	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	router := NewRouter()

	req, err := http.NewRequest("GET", homepageRoute, nil)
	assert.NoError(err)

	router.ServeHTTP(recorder, req)
	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusOK, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Staffjoy")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// Test that at least one asset loads
func TestAStaticAssetLoads(t *testing.T) {
	assets := []string{"/assets/css/main.css", "/assets/js/common.js", "/assets/images/logo.svg"}

	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	router := NewRouter()

	for _, asset := range assets {
		req, err := http.NewRequest("GET", asset, nil)
		assert.NoError(err)

		router.ServeHTTP(recorder, req)
		assert.Equal(http.StatusOK, recorder.Code)
		assert.NotZero(recorder.Body.Len())
	}
}

func TestLegacyPathRedirectsToSuite(t *testing.T) {
	inputToExpected := map[string]string{
		"/manager/organizations/1#locations/6":                                       "https://suite.staffjoy.com/manager/organizations/1#locations/6",
		"/myschedules/organizations/1/locations/6/roles/656/users/1#week/2016-08-15": "https://suite.staffjoy.com/myschedules/organizations/1/locations/6/roles/656/users/1#week/2016-08-15",
		"/euler/#schedule-monitoring":                                                "https://suite.staffjoy.com/euler/#schedule-monitoring",
		"/auth/change-password":                                                      "https://suite.staffjoy.com/auth/change-password",
		"/api/v2/internal/cron/":                                                     "https://suite.staffjoy.com/api/v2/internal/cron/",
	}

	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	router := NewRouter()

	for input, expected := range inputToExpected {
		// I don't think we need to worry about testing different methods because
		// the goal is preserving permalinks, not getting a post request to succeed
		req, err := http.NewRequest("GET", input, nil)
		assert.NoError(err)
		router.ServeHTTP(recorder, req)
		assert.Equal(http.StatusMovedPermanently, recorder.Code)
		assert.Equal(expected, recorder.Header().Get("Location"))
	}
}

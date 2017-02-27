package apidocs

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"v2.staffjoy.com/environments"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
)

func TestHomepageRendersValidHtml(t *testing.T) {
	assert := assert.New(t)
	rec := httptest.NewRecorder()

	config, err := environments.GetConfig("test")
	assert.NoError(err)
	p := page{logger: config.GetLogger("apidocstest")}
	mux, err := p.newMux()
	assert.NoError(err)

	req, err := http.NewRequest("GET", "/ui/", nil)
	assert.NoError(err)

	mux.ServeHTTP(rec, req)

	assert.Equal("text/html; charset=UTF-8", rec.Header().Get("Content-Type"))
	assert.Equal(http.StatusOK, rec.Code)
	assert.NotZero(rec.Body.Len())
	/*
		bodyHTML, err := ioutil.ReadAll(rec.Body)
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
	*/
	// Check for valid html
	_, err = html.Parse(rec.Body)
	assert.NoError(err)

}

// Test that at least one asset loads
func TestAStaticAssetLoads(t *testing.T) {
	assert := assert.New(t)

	assets := []string{"/ui/css/style.css", "/ui/js/swagger-ui.js"}
	rec := httptest.NewRecorder()
	config, err := environments.GetConfig("test")
	assert.NoError(err)
	p := page{logger: config.GetLogger("apidocstest")}
	mux, err := p.newMux()
	assert.NoError(err)

	for _, asset := range assets {
		req, err := http.NewRequest("GET", asset, nil)
		assert.NoError(err)

		mux.ServeHTTP(rec, req)
		assert.Equal(http.StatusOK, rec.Code)
		assert.NotZero(rec.Body.Len())
	}
}

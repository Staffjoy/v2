package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxyHeaders(t *testing.T) {
	assert := assert.New(t)
	testHeaders := map[string]string{
		"foo": "bar",
		"a":   "test",
		"1":   "2 3 4",
	}

	req, err := http.NewRequest("GET", "https://www.staffjoy.com/foo/bar/test.jpg?local=true", nil)
	assert.NoError(err)

	// Set up request
	for k, v := range testHeaders {
		req.Header.Set(k, v)
	}

	rec := httptest.NewRecorder()
	ProxyHeaders(req.Header, rec.Header())

	for k, v := range testHeaders {
		assert.Equal(v, rec.Header().Get(k))
	}
	// Check that the header lengths are same
	assert.Equal(len(testHeaders), len(rec.Header()))
}

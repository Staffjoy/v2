package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	assert := assert.New(t)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", HEALTHPATH, nil)
	assert.NoError(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(Handler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	assert.Equal(recorder.Code, http.StatusOK)

	// Check the response body is what we expect.
	assert.Equal(recorder.Body.String(), `{"hello":"world"}`)
}

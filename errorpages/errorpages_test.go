package errorpages

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/html"

	"github.com/stretchr/testify/assert"
)

type fakeTemplate struct{}

func (f *fakeTemplate) Execute(wr io.Writer, data interface{}) (err error) {
	return fmt.Errorf("There is an error")
}

// TestNotFound checks for a valid 404 page
func TestNotFound(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	NotFound(recorder)

	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusNotFound, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "doesn't exist")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// TestInternalServerError checks for a valid 500 page
func TestInternalServerError(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	InternalServerError(recorder)

	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusInternalServerError, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Server Error")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// TestTooManyRequests checks for a valid 429 page
func TestTooManyRequests(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	TooManyRequests(recorder)

	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusTooManyRequests, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Too Many Requests")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// TestForbidden checks for a valid 403 page
func TestForbidden(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	Forbidden(recorder)

	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusForbidden, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Forbidden")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// Test checks for a valid 403 page
func TestGatewayTimeout(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	GatewayTimeout(recorder)

	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusGatewayTimeout, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Timeout")

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

// TestFailedTemplateFallsBackToPlaintext checks that a plaintext message is returned if templating fails
func TestFailedTemplateFallsBackToPlaintext(t *testing.T) {
	assert := assert.New(t)

	p := &page{Title: "tester", HeaderCode: 314}
	tmpl := &fakeTemplate{}
	recorder := httptest.NewRecorder()

	p.writeResponse(recorder, tmpl)
	assert.Equal("text/plain", recorder.Header().Get("Content-Type"))
	assert.Equal(p.HeaderCode, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	body, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(body), p.Title)
}

func TestInternalErrorWithSentry(t *testing.T) {
	sentryDSN := "https://username@sentry.io/12345"
	errID := "foobar"
	assert := assert.New(t)
	recorder := httptest.NewRecorder()

	InternalServerErrorWithSentry(recorder, errID, sentryDSN)

	assert.Equal("text/html; charset=UTF-8", recorder.Header().Get("Content-Type"))
	assert.Equal(http.StatusInternalServerError, recorder.Code)
	assert.NotZero(recorder.Body.Len())

	bodyHTML, err := ioutil.ReadAll(recorder.Body)
	assert.NoError(err)
	assert.Contains(string(bodyHTML), "Error")
	assert.Contains(string(bodyHTML), errID)
	assert.Contains(string(bodyHTML), sentryDSN)

	// Check for valid html
	_, err = html.Parse(recorder.Body)
	assert.NoError(err)
}

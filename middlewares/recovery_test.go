package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/negroni"
	"v2.staffjoy.com/environments"
)

func TestRecovery(t *testing.T) {
	assert := assert.New(t)
	recorder := httptest.NewRecorder()
	serviceName := "faraday"
	conf := environments.Config{Debug: true}
	rec := NewRecovery(serviceName, conf, "")

	n := negroni.New()

	// replace log for testing
	n.Use(rec)

	// force a panic
	n.UseHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		panic("omg i'm panicing!")
	}))

	n.ServeHTTP(recorder, (*http.Request)(nil))

	assert.Equal(recorder.Code, http.StatusInternalServerError)
	assert.NotZero(recorder.Body.Len())
}

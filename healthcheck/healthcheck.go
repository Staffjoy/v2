// healthcheck is a library that provides a basic health check handler for Staffjoy applications.
// We generally host this endpoint at "/health" on port 80
//
// Usage:
// r.HandleFunc(healthcheck.HEALTHPATH healthcheck.Handler)

package healthcheck

import (
	"encoding/json"
	"net/http"
)

const (
	// HEALTHPATH is the standard healthcheck path in our app
	HEALTHPATH string = "/health"
)

// Handler returns a basic JSON
func Handler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	// We shouldn't have any errors
	msg, _ := json.Marshal(map[string]string{"hello": "world"})
	res.Write([]byte(msg))
	return
}

package main

import (
	"net/http"

	"v2.staffjoy.com/faraday/services"

	"github.com/gorilla/context"
	"v2.staffjoy.com/environments"
)

const (
	robotstxtPath  = "/robots.txt"
	robotstxtAllow = "User-agent: *\nDisallow:"   // Disallow nothing
	robotstxtDeny  = "User-agent: *\nDisallow: /" // Disallow everything
)

// RobotstxtMiddleware is a negroni middleware that determines whether search engines
// should access a service
type RobotstxtMiddleware struct {
	Config environments.Config
}

// NewRobotstxtMiddleware returns a new middleware for controlling search engines
func NewRobotstxtMiddleware(config environments.Config) *RobotstxtMiddleware {
	return &RobotstxtMiddleware{
		Config: config,
	}

}

func (svc *RobotstxtMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.URL.Path == robotstxtPath {
		// Tell search engine what to do!
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "text/plain")

		// Default to disallow
		service := context.Get(req, requestedService).(services.Service)
		var body string
		if (svc.Config.Name == "production") && (service.Security == services.Public) {
			body = robotstxtAllow
		} else {
			body = robotstxtDeny
		}
		res.Write([]byte(body))
		return
	}

	next(res, req)
}

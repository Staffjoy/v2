package main

import (
	"fmt"
	"net/http"
	"strings"

	"v2.staffjoy.com/faraday/services"

	"github.com/gorilla/context"
	"v2.staffjoy.com/environments"
)

const (
	// If we hit apex- where do we redirect users?
	defaultService = "www"
)

// ServiceMiddleware is a negroni middleware that matches requests against
// known Staffjoy services. If a match is found, Faraday proxies that app.
// If a match is not found, this middleware blocks the request
type ServiceMiddleware struct {
	Config   environments.Config
	Services services.ServiceDirectory
}

// NewServiceMiddleware returns a new middleware for services
func NewServiceMiddleware(config environments.Config, services services.ServiceDirectory) *ServiceMiddleware {
	return &ServiceMiddleware{
		Config:   config,
		Services: services,
	}

}

func (svc *ServiceMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// if you're hitting naked domain - go to www
	// e.g. staffjoy.com/foo?true=1 should redirect to www.staffjoy.com/foo?true=1
	if req.Host == svc.Config.ExternalApex {
		// It's hitting naked domain - redirect to www
		redirectURL := req.URL
		redirectURL.Host = defaultService + "." + svc.Config.ExternalApex
		if svc.Config.Debug {
			redirectURL.Scheme = "http"
		} else {
			redirectURL.Scheme = "https"
		}
		http.Redirect(res, req, redirectURL.String(), http.StatusMovedPermanently)
		return
	}

	subdomain, err := HostToService(req.Host, svc.Config.ExternalApex)
	if err != nil {
		res.WriteHeader(400)
		res.Write([]byte("Unsupported domain"))
		return
	}

	service, ok := svc.Services[subdomain]
	if !ok {
		res.WriteHeader(400)
		res.Write([]byte("Unsupported service"))
		return
	}

	context.Set(req, requestedService, service)

	next(res, req)
}

// HostToService finds the subdomain being accessed
func HostToService(host string, externalApex string) (string, error) {
	// First, prepend a "." to root domain
	externalApex = "." + strings.ToLower(externalApex)

	// For safety, set host to lowercase
	host = strings.ToLower(host)

	// Host may contain a port, so chop it off
	colonIndex := strings.Index(host, ":")
	if colonIndex > -1 {
		host = host[:colonIndex]
	}

	// Check that host ends with subdomain
	if len(host) <= len(externalApex) {
		return "", fmt.Errorf("Host is less than root domain length")
	}

	subdomainLength := len(host) - len(externalApex)
	if host[subdomainLength:] != externalApex {
		return "", fmt.Errorf("Does not contain root domain")
	}

	// Return subdomain
	return host[:subdomainLength], nil
}

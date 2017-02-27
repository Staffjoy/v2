package main

import (
	"net/http"
	"strings"

	"v2.staffjoy.com/environments"
)

// SecurityMiddleware is a negroni middleware that does nice things
// like HSTS and framebusting
type SecurityMiddleware struct {
	Config environments.Config
}

// NewSecurityMiddleware returns a new middleware for security
func NewSecurityMiddleware(config environments.Config) *SecurityMiddleware {
	return &SecurityMiddleware{
		Config: config,
	}

}

func (svc *SecurityMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// TODO - Determine how to force SSL. Depends on our load balancer config.

	if origin := req.Header.Get("Origin"); origin != "" {
		res.Header().Set("Access-Control-Allow-Origin", origin)
		res.Header().Set("Access-Control-Allow-Credentials", "true")
		res.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		res.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Cookie, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}

	if svc.Config.Debug == false {
		// Check if SSL
		isSSL := strings.EqualFold(req.URL.Scheme, "https") || req.TLS != nil
		if !isSSL {
			// Check if Cloudflare proxied it
			if req.Header.Get("X-Forwarded-Proto") == "https" {
				isSSL = true
			}
		}

		// If not SSL, then redirect.
		if !isSSL {
			url := req.URL
			url.Scheme = "https"
			url.Host = req.Host
			http.Redirect(res, req, url.String(), http.StatusMovedPermanently)
			return
		}

		// HSTS - force SSL
		res.Header().Add("Strict-Transport-Security", "max-age=315360000; includeSubDomains; preload")
		// No iFrames
		res.Header().Add("X-Frame-Options", "DENY")
		// Cross-site scripting protection
		res.Header().Add("X-XSS-Protection", "1; mode=block")
	}

	next(res, req)

}

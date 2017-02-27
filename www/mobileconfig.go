package main

import (
	"encoding/json"
	"net/http"
)

const (
	// MobileConfigPath  is a URL path that iPhone and Android apps check
	MobileConfigPath = "/mobileconfig.json"
	// MobileConfigRegex is a pattern for internal "apps"
	MobileConfigRegex = `^https?://(dev|stage|www)\.staffjoy\.com`
	// regexKey is the key in JSON used to find the MobileConfigRegex
	regexKey = "hideNavForURLsMatchingPattern"
)

// MobileConfigMiddleware is a negroni middleware that controls which URLs are
// treated as "external" on our mobile applications
type MobileConfigMiddleware struct{}

// NewMobileConfigMiddleware returns a new middleware for controlling search engines
func NewMobileConfigMiddleware() *MobileConfigMiddleware {
	return &MobileConfigMiddleware{}
}

func (svc *MobileConfigMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if req.URL.Path == MobileConfigPath {
		res.WriteHeader(http.StatusOK)
		res.Header().Set("Content-Type", "application/json")
		body, err := json.Marshal(map[string]string{regexKey: MobileConfigRegex})
		if err != nil {
			panic("Cannot encode mobile config")
		}
		res.Write(body)
		return
	}

	next(res, req)
}

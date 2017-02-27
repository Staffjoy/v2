package main

import (
	"encoding/json"
	"net/http"
)

const (
	// MobileConfigPath  is a URL path that iPhone and Android apps check
	MobileConfigPath = "/mobile-config.json"
	// MobileConfigRegex is a pattern for internal "apps"
	MobileConfigRegex = `^https?://(dev|stage|www|suite)\.staffjoy\.com`
	// regexKey is the key in JSON used to find the MobileConfigRegex
	regexKey = "hideNavForURLsMatchingPattern"
)

// MobileConfigHandler writes json for controlling what iPhone/Android
// applications consider an internal domain
func MobileConfigHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)
	res.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(map[string]string{regexKey: MobileConfigRegex})
	if err != nil {
		panic("Cannot encode mobile config")
	}
	res.Write(body)
}

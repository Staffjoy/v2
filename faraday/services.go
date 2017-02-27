package main

// Configuration for back-end services

const (
	// Public security means a user may be logged out or logged in
	Public = iota
	// Authenticated security means a user must be logged in
	Authenticated = iota
	// Admin security means a user must be both logged in and have sudo flag
	Admin = iota
)

// ServiceDirectory allows access to a backend service using its subdomain
type ServiceDirectory map[string]Service

// Service is an app on Staffjoy that runs on a subdomain
type Service struct {
	security      int    // Public, Authenticated, or Admin
	restrictDev   bool   // If True, service is suppressed in stage and prod
	backendDomain string // Backend service to query
}

// StaffjoyServices is a map of subdomains -> specs
// Sudomain is <string> + Env["rootDomain"]
// e.g. "login" service on prod is "login" + "staffjoy.com""
//
// KEEP THIS LIST IN ALPHABETICAL ORDER please
var StaffjoyServices = ServiceDirectory{
	"account": {
		security:      Authenticated,
		restrictDev:   false,
		backendDomain: "accountapi-service",
	},
	"app": {
		security:      Authenticated,
		restrictDev:   false,
		backendDomain: "app-service",
	},

	"code": {
		security:      Public,
		restrictDev:   false,
		backendDomain: "code-service",
	},
	"company": {
		security:      Authenticated,
		restrictDev:   false,
		backendDomain: "companyapi-service",
	},
	"faraday": {
		// Debug site for faraday
		security:      Public,
		restrictDev:   true,
		backendDomain: "httpbin.org",
	},
	"login": {
		security:      Public,
		restrictDev:   false,
		backendDomain: "login-service",
	},
	"myaccount": {
		security:      Authenticated,
		restrictDev:   false,
		backendDomain: "myaccount-service",
	},
	"superpowers": {
		security:      Authenticated,
		restrictDev:   true,
		backendDomain: "superpowers-service",
	},
	"signal": {
		security:      Admin,
		restrictDev:   false,
		backendDomain: "signal.staffjoy.com",
	},
	"waitlist": {
		security:      Public,
		restrictDev:   false,
		backendDomain: "waitlist-service",
	},
	"whoami": {
		security:      Public,
		restrictDev:   false,
		backendDomain: "whoami-service",
	},
	"www": {
		security:      Public,
		restrictDev:   false,
		backendDomain: "www-service",
	},
}

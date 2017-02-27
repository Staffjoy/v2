package services

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
	Security      int    // Public, Authenticated, or Admin
	RestrictDev   bool   // If True, service is suppressed in stage and prod
	BackendDomain string // Backend service to query
	NoCacheHTML   bool   // If True, injects a header for HTML responses telling the browser not to cache HTML

}

// StaffjoyServices is a map of subdomains -> specs
// Sudomain is <string> + Env["rootDomain"]
// e.g. "login" service on prod is "login" + "staffjoy.com""
//
// KEEP THIS LIST IN ALPHABETICAL ORDER please
var StaffjoyServices = ServiceDirectory{
	"account": {
		Security:      Authenticated,
		RestrictDev:   false,
		BackendDomain: "accountapi-service",
	},
	"app": {
		Security:      Authenticated,
		RestrictDev:   false,
		BackendDomain: "app-service",
		NoCacheHTML:   true,
	},
	"code": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "code-service",
	},
	"company": {
		Security:      Authenticated,
		RestrictDev:   false,
		BackendDomain: "companyapi-service",
	},
	"faraday": {
		// Debug site for faraday
		Security:      Public,
		RestrictDev:   true,
		BackendDomain: "httpbin.org",
	},
	"euler": {
		Security:      Admin,
		RestrictDev:   true,
		BackendDomain: "euler-service",
		NoCacheHTML:   true,
	},
	"ical": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "ical-service",
	},
	"login": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "login-service",
	},
	"myaccount": {
		Security:      Authenticated,
		RestrictDev:   false,
		BackendDomain: "myaccount-service",
		NoCacheHTML:   true,
	},
	"superpowers": {
		Security:      Authenticated,
		RestrictDev:   true,
		BackendDomain: "superpowers-service",
	},
	"signal": {
		Security:      Admin,
		RestrictDev:   false,
		BackendDomain: "signal.staffjoy.com",
	},
	"waitlist": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "waitlist-service",
	},
	"whoami": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "whoami-service",
	},
	"www": {
		Security:      Public,
		RestrictDev:   false,
		BackendDomain: "www-service",
	},
}

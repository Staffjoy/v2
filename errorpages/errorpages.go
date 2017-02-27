// Package errorpages standardizes common html error pages across our application.
package errorpages

import (
	"encoding/base64"
	"io"
	"net/http"
	"text/template"
)

// Split this out so we can mock in tests
type errorTemplate interface {
	Execute(io.Writer, interface{}) error
}

var (
	tmpl         *template.Template
	imageBase64  string
	notFoundPage = &page{
		Title:       "Oops! The page you were looking for doesn't exist.",
		Explanation: "You may have mistyped the address, or the page may have been moved.",
		HeaderCode:  http.StatusNotFound,
		LinkText:    "Click here to go back to Staffjoy",
		LinkHref:    "https://www.staffjoy.com",
	}
	internalServerErrorPage = &page{
		Title:       "Internal Server Error",
		Explanation: "Oops! Something broke. We're paging our engineers to look at it immediately.",
		HeaderCode:  http.StatusInternalServerError,
		LinkText:    "Click here to check out our system status page",
		LinkHref:    "https://status.staffjoy.com",
	}
	tooManyRequestsPage = &page{
		Title:       "Too Many Requests",
		Explanation: "Calm down - our system thinks that you are making too many requests.",
		HeaderCode:  http.StatusTooManyRequests,
		LinkText:    "Contact our support team for help",
		LinkHref:    "mailto:help@staffjoy.com",
	}
	forbiddenPage = &page{
		Title:       "Access Forbidden",
		Explanation: "Sorry, it looks like you do not have permission to access this page.",
		HeaderCode:  http.StatusForbidden,
		LinkText:    "Need help? Click here to contact our support team.",
		LinkHref:    "mailto:help@staffjoy.com",
	}
	timeoutPage = &page{
		Title:       "Timeout Error",
		Explanation: "Sorry, our servers seem to be slow. Please try again in a moment.",
		HeaderCode:  http.StatusGatewayTimeout,
		LinkText:    "Click here to check out our system status page",
		LinkHref:    "https://status.staffjoy.com",
	}
)

// Load the template
func init() {
	tmplData, err := Asset("assets/error.tmpl")
	if err != nil {
		panic("Unable to find error template in bindata")
	}

	tmpl, err = template.New("Error").Parse(string(tmplData))
	if err != nil {
		panic("Unable to parse error template")
	}

	imgFile, err := Asset("assets/staffjoy_coffee.png")
	if err != nil {
		panic("Unable to find error image in bindata")
	}
	imageBase64 = base64.StdEncoding.EncodeToString(imgFile)
}

type page struct {
	Title           string // Used in <title> and <h1>
	Explanation     string // Tell the user what's wrong
	HeaderCode      int    // http status code
	LinkText        string // Where do you want user to go?
	LinkHref        string // what's the link?
	SentryErrorID   string // What do we track the error as on the backend?
	SentryPublicDSN string // Config for app
	ImageBase64     string
}

func (p *page) writeResponse(res http.ResponseWriter, tmpl errorTemplate) {
	res.WriteHeader(p.HeaderCode)
	res.Header().Set("Content-Type", "text/html; charset=UTF-8")
	p.ImageBase64 = imageBase64

	err := tmpl.Execute(res, p)
	if err != nil {
		// Fall back to plaintextResponse
		p.writePlaintextResponse(res)
		return
	}
}

func (p *page) writePlaintextResponse(res http.ResponseWriter) {
	res.Header().Set("Content-Type", "text/plain")
	res.Write([]byte(p.Title))
}

// NotFound writes a 404 message
func NotFound(res http.ResponseWriter) {
	notFoundPage.writeResponse(res, tmpl)
}

// InternalServerError writes a 500 message
func InternalServerError(res http.ResponseWriter) {
	internalServerErrorPage.writeResponse(res, tmpl)
}

// InternalServerErrorWithSentry shows an error page and collects user feedback for debugging purposes
func InternalServerErrorWithSentry(res http.ResponseWriter, sentryErrorID, sentryPublicDSN string) {
	// Copy the internal error page
	customErrorPage := *internalServerErrorPage
	customErrorPage.SentryErrorID = sentryErrorID
	customErrorPage.SentryPublicDSN = sentryPublicDSN
	customErrorPage.writeResponse(res, tmpl)
}

// TooManyRequests writes a 429 message
func TooManyRequests(res http.ResponseWriter) {
	tooManyRequestsPage.writeResponse(res, tmpl)
}

// Forbidden writes a 403 message
func Forbidden(res http.ResponseWriter) {
	forbiddenPage.writeResponse(res, tmpl)
}

// GatewayTimeout writes a 504 message
func GatewayTimeout(res http.ResponseWriter) {
	timeoutPage.writeResponse(res, tmpl)
}

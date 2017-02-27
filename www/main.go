// www is the main marketing site for Staffjoy
//
// It uses the Go backend to process new account signups and provide
// security like CSRF
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
	"v2.staffjoy.com/healthcheck"
	"v2.staffjoy.com/middlewares"

	"github.com/russross/blackfriday"

	"github.com/urfave/negroni"
)

var (
	// CSRF is the cross site request forgery secret
	CSRF http.Handler

	logger       *logrus.Entry
	config       environments.Config
	tmpl         *template.Template
	signingToken = os.Getenv("SIGNING_SECRET")

	// Subfolders that are served directly
	assetPaths = []string{"assets/css", "assets/images", "assets/js", "assets/data", "assets/fonts", "assets/breaktime-cover"}

	// Paths that we are 301 redirecting to suite.staffjoy.com
	legacyPaths = []string{"/api/v2/", "/auth/", "/euler/", "/myschedules/", "/manager/"}

	// Register new marketing pages here
	// TODO - make a sitemap.xml based on this
	staticPages = map[string]*page{
		"/":                {Title: "Staffjoy - Online Scheduling Software", Description: "Staffjoy is a web application that helps small businesses create schedules online and automatically communicate them via text message with hourly workers.", TemplateName: "home.tmpl"},
		"/about/":          {Title: "About Staffjoy", Description: "Learn about the members of the Staffjoy team and the origin of the company.", TemplateName: "about.tmpl"},
		"/careers/":        {Title: "Staffjoy Careers", Description: "If you’re looking to improve the way small businesses schedule their hourly workers, you are invited to apply to join our team in San Francisco.", TemplateName: "careers.tmpl", CSSId: "careers"},
		"/pricing/":        {Title: "Staffjoy Pricing", Description: "Staffjoy’s software pricing is affordable for any size team. There is a monthly subscription based on the number of employees your company has.", TemplateName: "pricing.tmpl"},
		"/privacy-policy/": {Title: "Staffjoy Privacy Policy", Description: "Staffjoy’s Privacy Policy will walk you through through security protocols, data storage, and legal compliance that all clients need to know. ", TemplateName: "privacypolicy.tmpl"},
		"/sign-up/":        {Title: "Sign Up for Your 30 Day Free Staffjoy Trial", Description: "Sign up for a 30 day free trial of Staffjoy today to create your schedule online. We’ll distribute it to your team using automated text messages.", TemplateName: "signup.tmpl", CSSId: "sign-up"},
		"/early-access/":   {Title: "Early Access Signup", Description: "Get early access for Staffjoy", TemplateName: "early.tmpl", CSSId: "sign-up"},
		"/terms/":          {Title: "Staffjoy Terms and Conditions", Description: "Staffjoy’s Terms and Conditions point out the liability, disclaimers, exclusions, and more that all users of our website must agree to.", TemplateName: "terms.tmpl"},
	}
	confirmPage      = &page{Title: "Open your email and click on the confirmation link!", Description: "Check your email and click the link for next steps", TemplateName: "confirm.tmpl", CSSId: "confirm"}
	resetConfirmPage = &page{Title: "Please check your email for a reset link!", Description: "Check your email and click the link for next steps", TemplateName: "confirm.tmpl", CSSId: "confirm"}
	newCompanyPage   = &page{Title: "Create a new company", Description: "Get started with a new Staffjoy account", TemplateName: "new_company.tmpl", CSSId: "newCompany"}
	breaktimeSource  = make(map[string]string)
)

const (
	// ServiceName is how we refer to this app in logs
	ServiceName = "www"
	// All templates in this folder will be loaded
	templateFolder       = "assets/templates"
	breaktimeAssetFolder = "assets/breaktime-content"

	// For SEO / web crawlers
	defaultDescription = "Staffjoy is an application that helps businesses create and share schedules with hourly workers."
	// URL where we should redirect legacy links
	legacyURL = "https://suite.staffjoy.com"

	passwordResetPath = "/password-reset/"
	confirmTemplate   = "confirm.tmpl"
)

// Added in template
func hasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func init() {
	var err error
	// Set the ENV environment variable to control dev/stage/prod behavior
	config, err = environments.GetConfig(os.Getenv(environments.EnvVar))
	if err != nil {
		panic("Unable to determine configuration")
	}
	logger = config.GetLogger(ServiceName)

	// Load templates
	templateFilenames, err := AssetDir(templateFolder)
	if err != nil {
		logger.Panicf("Unable to load template files: %s", err)
	}
	for _, name := range templateFilenames {
		tmplData, err := Asset(fmt.Sprintf("%s/%s", templateFolder, name))
		if err != nil {
			logger.Panicf("Unable to locate specified asset - %s", err)
		}
		// Create template on first loop
		if tmpl == nil {
			tmpl, err = template.New(name).Funcs(template.FuncMap{"hasField": hasField}).Parse(string(tmplData))
		} else {
			tmpl, err = tmpl.New(name).Funcs(template.FuncMap{"hasField": hasField}).Parse(string(tmplData))
		}
		if err != nil {
			logger.Panicf("Unable to parse template - %s", err)
		}
	}

	breaktimeFilenames, err := AssetDir(breaktimeAssetFolder)
	if err != nil {
		logger.Panicf("Unable to load breaktime files: %s", err)
	}
	for _, name := range breaktimeFilenames {
		sourceData, err := Asset(fmt.Sprintf("%s/%s", breaktimeAssetFolder, name))
		if err != nil {
			logger.Panicf("Unable to locate specified asset - %s", err)
		}
		// Create template on first loop
		breaktimeSource[name] = string(blackfriday.MarkdownBasic(sourceData))
	}

	if len(signingToken) == 0 && !config.Debug {
		panic("no signing token")
	}

	logger.Debugf("Initialized www %s environment", config.Name)
}

// NewRouter builds the mux router for the site
// (abstracted for testing purposes)
func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc(healthcheck.HEALTHPATH, healthcheck.Handler)
	r.HandleFunc("/confirm/", signUpHandler).Methods(http.MethodPost)
	r.HandleFunc("/activate/{token}", activateHandler)
	r.HandleFunc("/reset/{token}", confirmResetHandler)
	r.HandleFunc("/login/", loginHandler)
	r.HandleFunc("/logout/", logoutHandler)
	r.HandleFunc("/new-company/", newCompanyHandler)
	r.HandleFunc("/breaktime/", breaktimeListHandler)
	r.HandleFunc("/breaktime/{slug}", breaktimeEpisodeHandler)
	r.HandleFunc(passwordResetPath, resetHandler)

	// Register static pages
	version := config.GetDeployVersion()
	for route, info := range staticPages {
		info.Version = version
		r.HandleFunc(route, info.Handler)
	}
	confirmPage.Version = version
	resetConfirmPage.Version = version
	newCompanyPage.Version = version

	// Register asset folders we want served externally
	for _, path := range assetPaths {
		urlPath := fmt.Sprintf("/%s/", path) // Wrap in slashes
		r.PathPrefix(urlPath).Handler(http.StripPrefix(urlPath, http.FileServer(
			&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: path})))
	}

	// redirect old routes to suite.staffjoy.com for legacy users
	for _, path := range legacyPaths {
		r.PathPrefix(path).HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			http.Redirect(res, req, fmt.Sprintf("%s%s", legacyURL, req.URL), http.StatusMovedPermanently)
		})
	}

	r.NotFoundHandler = http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		errorpages.NotFound(res)
	})

	return r
}

func main() {
	r := NewRouter()
	n := negroni.New()

	sentryPublicDSN, err := environments.GetPublicSentryDSN(config.GetSentryDSN())
	if err != nil {
		logger.Fatalf("Cannot get sentry info - %s", err)
	}
	n.Use(middlewares.NewRecovery(ServiceName, config, sentryPublicDSN))
	n.UseHandler(r)

	CSRF := csrf.Protect(
		[]byte(signingToken),
		csrf.Domain(config.ExternalApex),
		csrf.Secure(!config.Debug),
		csrf.Path("/"),
		csrf.CookieName("sjcsrf"),
		csrf.ErrorHandler(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			logger.Infof("failed CSRF - %s", csrf.FailureReason(req))
			errorpages.Forbidden(res)
		})),
		csrf.FieldName("csrf"),
	)

	s := &http.Server{
		Addr:           ":80",
		Handler:        CSRF(n),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Panicf("%s", s.ListenAndServe())
}

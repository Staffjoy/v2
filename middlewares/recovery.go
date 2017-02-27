package middlewares

import (
	"net/http"
	"runtime"

	"github.com/Sirupsen/logrus"
	sentry "github.com/getsentry/raven-go"
	"v2.staffjoy.com/environments"
	"v2.staffjoy.com/errorpages"
)

// Recovery is a Negroni middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	ServiceName     string // for logging
	Logger          *logrus.Entry
	ErrorClient     environments.SentryClient // Can use this to send prod errors to a system
	StackAll        bool
	StackSize       int
	SentryPublicDSN string
}

// NewRecovery returns a new instance of Recovery
func NewRecovery(serviceName string, conf environments.Config, sentryPublicDSN string) *Recovery {
	return &Recovery{
		ServiceName: serviceName,
		Logger: conf.GetLogger(serviceName).WithFields(logrus.Fields{
			"middleware": "Recovery",
		}),
		StackAll:        false,
		StackSize:       1024 * 8,
		ErrorClient:     environments.ErrorClient(&conf),
		SentryPublicDSN: sentryPublicDSN,
	}
}

func (mw *Recovery) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if mw.ErrorClient != nil {
		// Sending errors to sentry
		mw.ErrorClient.SetHttpContext(sentry.NewHttp(req))
		// TODO - set user context https://blog.sentry.io/2015/07/07/logging-go-errors.html
		err, errID := mw.ErrorClient.CapturePanicAndWait(func() {
			next(res, req)
		}, map[string]string{"service": mw.ServiceName},
		)
		mw.ErrorClient.ClearContext()
		if err != nil {
			mw.Logger.Warningf("Reported error id %s to sentry (err %v)", errID, err)
			errorpages.InternalServerErrorWithSentry(res, errID, mw.SentryPublicDSN)
		}
	} else { // Development
		defer func() {
			if err := recover(); err != nil {
				stack := make([]byte, mw.StackSize)
				stack = stack[:runtime.Stack(stack, mw.StackAll)]
				mw.Logger.WithFields(logrus.Fields{
					"stack": string(stack),
					"err":   err,
				}).Errorf("PANIC RECOVERED- %s", err)
				errorpages.InternalServerError(res)
			}
		}()
		next(res, req)
	}
}

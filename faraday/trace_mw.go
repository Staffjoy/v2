package main

import (
	"fmt"
	"net/http"

	"google.golang.org/api/option"

	"golang.org/x/net/context"

	"cloud.google.com/go/trace"
	"github.com/Sirupsen/logrus"
	"v2.staffjoy.com/environments"
)

const (
	testRate  = 1.0 // 100%
	testLimit = 100 // Max 100 req/s
)

var (
	projectID string
)

func init() {
	projectID = environments.GetGoogleCloudProject()
}

// TraceMiddleware is a negroni middleware that reports URL timing to
// Google Cloud
type TraceMiddleware struct {
	Logger      *logrus.Entry
	Config      environments.Config
	TraceClient *trace.Client
}

// NewTraceMiddleware returns a new middleware for traces
func NewTraceMiddleware(logger *logrus.Entry, config environments.Config) (*TraceMiddleware, error) {
	var traceClient *trace.Client
	if !config.Debug {
		var err error
		traceClient, err = trace.NewClient(context.Background(), projectID, option.WithServiceAccountFile(environments.GoogleCloudSecretPath))
		if err != nil {
			return nil, fmt.Errorf("Unable to initialize trace client - %v", err)
		}

		policy, err := trace.NewLimitedSampler(testRate, testLimit)
		if err != nil {
			return nil, fmt.Errorf("Unable to initialize trace client - %v", err)
		}
		traceClient.SetSamplingPolicy(policy)
	}

	logger.Infof("trace client init %v", traceClient)
	return &TraceMiddleware{
		Logger: logger.WithFields(logrus.Fields{
			"middleware": "TraceMiddleware",
		}),
		Config:      config,
		TraceClient: traceClient,
	}, nil

}

func (svc *TraceMiddleware) ServeHTTP(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	var span *trace.Span

	if !svc.Config.Debug {
		// Start trace
		span = svc.TraceClient.SpanFromRequest(req)
	}

	next(res, req)

	if !svc.Config.Debug {
		// Finish trace
		go func(span *trace.Span) {
			err := span.FinishWait()
			if err != nil {
				svc.Logger.Warningf("Could not finish trace - %v", err)
			}
		}(span)
	}
}

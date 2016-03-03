package jshapi

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/go-stdlogger"
)

// Logger can be overridden with your own logger to utilize any custom features
// it might have. Interface defined here: https://github.com/derekdowling/go-stdlogger/blob/master/logger.go
var Logger std.Logger = log.New(os.Stderr, "jshapi: ", log.LstdFlags)

// SendAndLog is a jsh wrapper function that first prepares a jsh.Sendable response,
// and then handles logging 5XX errors that it encounters in the process.
func SendAndLog(ctx context.Context, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {

	intentionalErr, isType := sendable.(jsh.ErrorType)
	if isType {
		// determine error status before taking any additional actions
		var status int

		list, isList := intentionalErr.(jsh.ErrorList)
		if isList {
			status = list[0].Status
		}

		err, isErr := intentionalErr.(*jsh.Error)
		if isErr {
			status = err.Status
		}

		if status >= 500 {
			Logger.Printf("Returning ISE: %s\n", intentionalErr.Error())
		}
	}

	sendErr := jsh.Send(w, r, sendable)
	if sendErr != nil && sendErr.Status >= 500 {
		Logger.Printf("Error sending response: %s\n", sendErr.Error())
	}
}

// LeveledLogger is a context-aware logger interface that differentiates between
// log levels (debug, info, warning, error, critical).
type LeveledLogger interface {
	Debugf(ctx context.Context, format string, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warningf(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Criticalf(ctx context.Context, format string, args ...interface{})
}

// StandardLogger is a leveled logger that calls Printf on an std.Logger for all
// levels except debug. Debug logging can be turned on by setting Debug = true
type StandardLogger struct {
	std.Logger
	Debug bool
}

// Debugf logs the message to the std.Logger using Printf only if Debug is true
func (s *StandardLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	if s.Debug {
		s.Printf(format, args...)
	}
}

// Infof redirects to Printf of std.Logger
func (s *StandardLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	s.Printf(format, args...)
}

// Warningf redirects to Printf of std.Logger
func (s *StandardLogger) Warningf(ctx context.Context, format string, args ...interface{}) {
	s.Printf(format, args...)
}

// Errorf redirects to Printf of std.Logger
func (s *StandardLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	s.Printf(format, args...)
}

// Criticalf redirects to Printf of std.Logger
func (s *StandardLogger) Criticalf(ctx context.Context, format string, args ...interface{}) {
	s.Printf(format, args...)
}

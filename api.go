package jshapi

import (
	"fmt"
	"path"
	"strings"

	"net/http"

	"github.com/derekdowling/go-json-spec-handler"
	"github.com/derekdowling/goji2-logger"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
)

// API is used to direct HTTP requests to resources
type API struct {
	*goji.Mux
	prefix    string
	Resources map[string]*Resource
	Debug     bool
	Logger    LeveledLogger
}

// New initializes a new top level API Resource Handler. The most basic implementation
// is:
//
//	// optionally, set your own logger
//	jshapi.Logger = yourLogger
//
//	// create a new API
//	api := jshapi.New("<prefix>", nil)
func New(prefix string, debug bool) *API {

	// create our new logger
	api := NewWithLogger(prefix, debug, &StandardLogger{Logger, debug})

	// register logger middleware
	gojilogger := gojilogger.New(Logger, debug)
	api.UseC(gojilogger.Middleware)

	return api
}

// NewWithLogger initializes a new top level API Resource Handler with a
// context-aware logger. The most basic implementation is:
//
//	// create a new API
//	api := jshapi.NewWithLogger("<prefix>", true, YourContextAwareLogger)
func NewWithLogger(prefix string, debug bool, logger LeveledLogger) *API {

	// ensure that our top level prefix is "/" prefixed
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	// create our new logger
	api := &API{
		Mux:       goji.NewMux(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
		Debug:     debug,
		Logger:    logger,
	}

	return api
}

// Add implements mux support for a given resource which is effectively handled as:
// pat.New("/(prefix/)resource.Plu*)
func (a *API) Add(resource *Resource) {

	// track our associated resources, will enable auto-generation docs later
	a.Resources[resource.Type] = resource
	resource.api = a

	// Because of how prefix matches work:
	// https://godoc.org/github.com/goji/goji/pat#hdr-Prefix_Matches
	// We need two separate routes,
	// /(prefix/)resources
	matcher := path.Join(a.prefix, resource.Type)
	a.Mux.HandleC(pat.New(matcher), resource)

	// And:
	// /(prefix/)resources/*
	idMatcher := path.Join(a.prefix, resource.Type, "*")
	a.Mux.HandleC(pat.New(idMatcher), resource)
}

// RouteTree prints out all accepted routes for the API that use jshapi implemented
// ways of adding routes through resources: NewCRUDResource(), .Get(), .Post, .Delete(),
// .Patch(), .List(), and .NewAction()
func (a *API) RouteTree() string {
	var routes string

	for _, resource := range a.Resources {
		routes = strings.Join([]string{routes, resource.RouteTree()}, "")
	}

	return routes
}

// SendAndLog is a jsh wrapper function that first prepares a jsh.Sendable response,
// and then handles logging 5XX errors that it encounters in the process.
func (a *API) SendAndLog(ctx context.Context, w http.ResponseWriter, r *http.Request, sendable jsh.Sendable) {
	if a == nil {
		SendAndLog(ctx, w, r, sendable)
		return
	}
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
			a.Logger.Errorf(ctx, "Returning ISE: %s\n", intentionalErr.Error())
		} else if status >= 400 {
			a.Logger.Warningf(ctx, "Client error: %s\n", intentionalErr.Error())
		}
	}

	sendErr := jsh.Send(w, r, sendable)
	if sendErr != nil && sendErr.Status >= 500 {
		a.Logger.Criticalf(ctx, "Error sending response: %s\n", sendErr.Error())
	}
}

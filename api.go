package jshapi

import (
	"fmt"
	"path"
	"strings"

	"goji.io"
	"goji.io/pat"

	"github.com/derekdowling/goji2-logger"
)

// API is used to direct HTTP requests to resources
type API struct {
	*goji.Mux
	prefix    string
	Resources map[string]*Resource
	Debug     bool
}

/*
New initializes a new top level API Resource Handler without. */
func New(prefix string) *API {

	// ensure that our top level prefix is "/" prefixed
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	// create our new API
	return &API{
		Mux:       goji.NewMux(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
	}
}

/*
Default builds a new top-level API with a few out of the box additions to get people
started without needing to add a lot of extra functionality.

The most basic implementation is:

	// optionally, set your own logger that can do custom things
	jshapi.Logger = yourLogger

	// create the API. Specify a http://yourapi/<prefix>/ if required
	api := jshapi.Default("<prefix>", false)
	api.Add(yourResource)

*/
func Default(prefix string, debug bool) *API {

	api := New(prefix)

	// register logger middleware
	gojilogger := gojilogger.New(Logger, debug)
	api.UseC(gojilogger.Middleware)

	return api
}

// Add implements mux support for a given resource which is effectively handled as:
// pat.New("/(prefix/)resource.Plu*)
func (a *API) Add(resource *Resource) {

	// track our associated resources, will enable auto-generation docs later
	a.Resources[resource.Type] = resource

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

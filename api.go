package jshapi

import (
	"fmt"
	"log"
	"os"
	"strings"

	"goji.io"
	"goji.io/pat"
)

// API is used to direct HTTP requests to resources
type API struct {
	*goji.Mux
	prefix    string
	Resources map[string]*Resource
	Logger    *log.Logger
}

// New initializes a new top level API Resource Handler.
func New(prefix string, logger *log.Logger) *API {

	// ensure that our top level prefix is "/" prefixed
	if !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

	if logger == nil {
		logger = log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile)
	}

	return &API{
		Mux:       goji.NewMux(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
		Logger:    logger,
	}
}

// AddResource adds a new resource of type "name" to the API's router
func (a *API) AddResource(resource *Resource) {

	// add prefix and logger
	resource.prefix = a.prefix
	resource.Logger = a.Logger

	a.Resources[resource.Type] = resource

	// Add subrouter to main API mux, use Matcher plus catch all
	a.Mux.HandleC(pat.New(resource.Matcher()+"*"), resource)
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

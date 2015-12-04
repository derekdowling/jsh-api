package jshapi

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zenazn/goji/web"
)

// API is used to direct HTTP requests to resources
type API struct {
	resources []*resource
	Prefix    string
	Logger    *log.Logger
	Mux       *web.Mux
}

// New initializes a Handler object
func New(prefix string) *API {
	return &API{
		resources: []*resource{},
		Prefix:    prefix,
		Logger:    log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile),
		Mux:       web.New(),
	}
}

// Handle is an http.Handle compatible function that should be used to passed the
// jshapi API into a router
func (a *API) Handle(w http.ResponseWriter, r *http.Request) {
	a.Mux.ServeHTTP(w, r)
}

// AddResource adds a new resource of type "name" to the API's router
func (a *API) AddResource(name string, storage Storage) {

	resource := &resource{
		Type:    name,
		Storage: storage,
		Logger:  a.Logger,
		Prefix:  a.Prefix,
		mux:     web.New(),
	}

	plural := fmt.Sprintf("/%ss", name)
	pluralSelecter := fmt.Sprintf("%s/:id", name)

	// setup resource sub-router
	resource.mux.Post(plural, resource.Post)
	resource.mux.Delete(plural, resource.Delete)
	resource.mux.Get(plural, resource.List)
	resource.mux.Get(pluralSelecter, resource.Get)
	resource.mux.Patch(pluralSelecter, resource.Patch)

	a.resources = append(a.resources, resource)

	// Add subrouter to main API mux
	a.Mux.Handle(plural+"/*", resource.mux)
}

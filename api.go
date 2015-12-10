package jshapi

import (
	"log"
	"os"

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

// New initializes a Handler object
func New(prefix string) *API {
	return &API{
		Mux:       goji.NewMux(),
		prefix:    prefix,
		Resources: map[string]*Resource{},
		Logger:    log.New(os.Stdout, "jshapi: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// AddResource adds a new resource of type "name" to the API's router
func (a *API) AddResource(resource *Resource) {

	// add prefix and logger
	resource.Prefix = a.prefix
	resource.Logger = a.Logger

	a.Resources[resource.Name] = resource

	// Add subrouter to main API mux, use Matcher plus catch all
	a.Mux.Handle(pat.New(resource.Matcher()+"*"), resource.Mux)
}

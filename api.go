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

// New initializes a Handler object
func New(prefix string) *API {

	if len(prefix) > 0 && !strings.HasPrefix(prefix, "/") {
		prefix = fmt.Sprintf("/%s", prefix)
	}

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
	resource.prefix = a.prefix
	resource.Logger = a.Logger

	a.Resources[resource.Type] = resource

	// Add subrouter to main API mux, use Matcher plus catch all
	a.Mux.HandleC(pat.New(resource.Matcher()+"*"), resource)
}
